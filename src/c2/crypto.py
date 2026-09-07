"""Authenticated message encryption for the ANGEL control protocol."""

from __future__ import annotations

import base64
import hashlib
import hmac
import json
import os
import threading
import time
from dataclasses import dataclass
from typing import Any, cast

from cryptography.exceptions import InvalidSignature
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives.ciphers.aead import AESGCM
from cryptography.hazmat.primitives.kdf.hkdf import HKDF


class CryptoError(ValueError):
    """Raised when an authenticated message cannot be processed."""


@dataclass(frozen=True)
class KeyPair:
    private_key: ec.EllipticCurvePrivateKey

    @classmethod
    def generate(cls) -> KeyPair:
        return cls(ec.generate_private_key(ec.SECP256R1()))

    @property
    def public_bytes(self) -> bytes:
        return self.private_key.public_key().public_bytes(
            serialization.Encoding.X962,
            serialization.PublicFormat.UncompressedPoint,
        )

    @classmethod
    def from_private_bytes(cls, data: bytes) -> KeyPair:
        key = serialization.load_der_private_key(data, password=None)
        if not isinstance(key, ec.EllipticCurvePrivateKey):
            raise CryptoError("private key is not an EC key")
        return cls(key)

    def private_bytes(self) -> bytes:
        return self.private_key.private_bytes(
            serialization.Encoding.DER,
            serialization.PrivateFormat.PKCS8,
            serialization.NoEncryption(),
        )


def derive_session_key(
    private_key: ec.EllipticCurvePrivateKey, peer_public: bytes, context: bytes = b"ANGEL-v1"
) -> bytes:
    try:
        peer = ec.EllipticCurvePublicKey.from_encoded_point(ec.SECP256R1(), peer_public)
    except ValueError as exc:
        raise CryptoError("invalid peer public key") from exc
    shared = private_key.exchange(ec.ECDH(), peer)
    return HKDF(algorithm=hashes.SHA256(), length=64, salt=None, info=context).derive(shared)


def sign_message(private_key: ec.EllipticCurvePrivateKey, message: bytes) -> bytes:
    return private_key.sign(message, ec.ECDSA(hashes.SHA256()))


def verify_signature(public_key: ec.EllipticCurvePublicKey, message: bytes, signature: bytes) -> bool:
    try:
        public_key.verify(signature, message, ec.ECDSA(hashes.SHA256()))
    except InvalidSignature:
        return False
    return True


def _canonical_bytes(value: dict[str, Any]) -> bytes:
    return json.dumps(value, sort_keys=True, separators=(",", ":")).encode()


@dataclass
class ReplayGuard:
    ttl_seconds: int = 300

    def __post_init__(self) -> None:
        self._seen: set[str] = set()
        self._lock = threading.Lock()

    def accept(self, message_id: str, timestamp: int, now: int | None = None) -> None:
        current = int(time.time()) if now is None else now
        if abs(current - timestamp) > self.ttl_seconds:
            raise CryptoError("message expired")
        with self._lock:
            if message_id in self._seen:
                raise CryptoError("replayed message")
            self._seen.add(message_id)
            if len(self._seen) > 10000:
                self._seen = set(list(self._seen)[-5000:])


class KeyRegistry:
    """Thread-safe key registry supporting explicit revocation and rotation."""

    def __init__(self) -> None:
        self._keys: dict[str, bytes] = {}
        self._revoked: set[str] = set()
        self._lock = threading.Lock()

    def register(self, key_id: str, material: bytes) -> None:
        if not key_id.strip() or len(material) < 64:
            raise CryptoError("invalid key registration")
        with self._lock:
            if key_id in self._revoked:
                raise CryptoError("key id is revoked")
            self._keys[key_id] = bytes(material)

    def revoke(self, key_id: str) -> None:
        with self._lock:
            self._keys.pop(key_id, None)
            self._revoked.add(key_id)

    def get(self, key_id: str) -> bytes:
        with self._lock:
            if key_id in self._revoked:
                raise CryptoError("key id is revoked")
            try:
                return self._keys[key_id]
            except KeyError as exc:
                raise CryptoError("unknown key id") from exc


class SessionCipher:
    """Encrypts JSON messages with AES-GCM and an explicit HMAC integrity tag."""

    def __init__(self, key_material: bytes) -> None:
        if len(key_material) < 64:
            raise CryptoError("session material must contain encryption and MAC keys")
        self._enc_key = key_material[:32]
        self._mac_key = key_material[32:64]

    def encrypt(self, payload: dict[str, Any], *, aad: bytes = b"") -> dict[str, Any]:
        message_id = base64.urlsafe_b64encode(os.urandom(18)).decode().rstrip("=")
        timestamp = int(time.time())
        nonce = os.urandom(12)
        body = {"id": message_id, "ts": timestamp, "payload": payload}
        plaintext = _canonical_bytes(body)
        ciphertext = AESGCM(self._enc_key).encrypt(nonce, plaintext, aad)
        mac_input = aad + nonce + ciphertext
        tag = hmac.new(self._mac_key, mac_input, hashlib.sha256).digest()
        return {
            "v": 1,
            "id": message_id,
            "ts": timestamp,
            "nonce": base64.b64encode(nonce).decode(),
            "data": base64.b64encode(ciphertext).decode(),
            "mac": base64.b64encode(tag).decode(),
        }

    def decrypt(
        self, envelope: dict[str, Any], *, aad: bytes = b"", replay: ReplayGuard | None = None
    ) -> dict[str, Any]:
        try:
            if envelope.get("v") != 1:
                raise CryptoError("unsupported envelope version")
            message_id = str(envelope["id"])
            timestamp = int(envelope["ts"])
            nonce = base64.b64decode(envelope["nonce"], validate=True)
            ciphertext = base64.b64decode(envelope["data"], validate=True)
            supplied_mac = base64.b64decode(envelope["mac"], validate=True)
        except (KeyError, TypeError, ValueError) as exc:
            raise CryptoError("malformed envelope") from exc
        if len(nonce) != 12:
            raise CryptoError("invalid nonce")
        expected = hmac.new(self._mac_key, aad + nonce + ciphertext, hashlib.sha256).digest()
        if not hmac.compare_digest(expected, supplied_mac):
            raise CryptoError("integrity check failed")
        if replay is not None:
            replay.accept(message_id, timestamp)
        try:
            body = json.loads(AESGCM(self._enc_key).decrypt(nonce, ciphertext, aad))
        except Exception as exc:
            raise CryptoError("decryption failed") from exc
        if body.get("id") != message_id or int(body.get("ts", -1)) != timestamp:
            raise CryptoError("envelope metadata mismatch")
        return cast(dict[str, Any], body["payload"])


# Backward-compatible helpers for callers that use the original module API.
def encrypt_message(plaintext: str) -> bytes:
    key = hashlib.sha512(os.environ.get("ANGEL_SHARED_KEY", "development-only-key").encode()).digest()
    cipher = SessionCipher(key)
    envelope = cipher.encrypt({"text": plaintext})
    return json.dumps(envelope, sort_keys=True).encode()


def decrypt_message(ciphertext: bytes) -> str:
    key = hashlib.sha512(os.environ.get("ANGEL_SHARED_KEY", "development-only-key").encode()).digest()
    cipher = SessionCipher(key)
    payload = cipher.decrypt(json.loads(ciphertext.decode()))
    return str(payload["text"])


__all__ = [
    "CryptoError",
    "KeyPair",
    "KeyRegistry",
    "ReplayGuard",
    "SessionCipher",
    "derive_session_key",
    "sign_message",
    "verify_signature",
    "encrypt_message",
    "decrypt_message",
]
