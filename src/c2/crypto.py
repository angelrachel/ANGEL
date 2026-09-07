import os

from cryptography.hazmat.primitives.ciphers.aead import AESGCM

from .shared_key import SHARED_KEY


def encrypt_message(plaintext):
    aesgcm = AESGCM(SHARED_KEY)
    nonce = os.urandom(12)
    ciphertext = aesgcm.encrypt(nonce, plaintext.encode(), None)
    return nonce + ciphertext


def decrypt_message(ciphertext):
    aesgcm = AESGCM(SHARED_KEY)
    nonce = ciphertext[:12]
    ct = ciphertext[12:]
    return aesgcm.decrypt(nonce, ct, None).decode()
