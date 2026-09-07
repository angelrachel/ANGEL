from __future__ import annotations

import pytest

from src.rbac import RBAC, RBACError


def test_rbac_authorizes_and_revokes_principal() -> None:
    rbac = RBAC()
    rbac.define_role("operator", {"task.queue", "audit.read"})
    rbac.upsert_principal("alice", {"operator"})
    rbac.authorize("alice", "task.queue")
    with pytest.raises(RBACError, match="denied"):
        rbac.authorize("alice", "admin.delete")
    rbac.revoke("alice")
    with pytest.raises(RBACError, match="revoked"):
        rbac.authorize("alice", "task.queue")


def test_rbac_rejects_unknown_roles() -> None:
    with pytest.raises(RBACError, match="invalid"):
        RBAC().upsert_principal("alice", {"missing"})
