# ANGEL Operations Runbook

## Before an engagement

Confirm a signed authorization, an active Rules of Engagement record, the target inventory, the emergency contact, the data-retention period, and the approved assessment mode. Use synthetic data in development and staging. Production secrets must come from a managed secret store.

## Start-up checks

Run the test suite and security checks before deployment. Inject `ANGEL_OPERATOR_KEY`, `ANGEL_SHARED_KEY`, and `ANGEL_AUTH_SECRET` from a managed secret store; the operator API intentionally fails closed when `ANGEL_OPERATOR_KEY` is unset and has no source-code default. Bind the service to a private interface and place it behind the hardened reverse proxy. Confirm the health and readiness endpoints, database permissions, and backup destination.

## Assessment workflow

Create the scope with exact hosts and paths. Create an auditable assessment plan using only passive or simulation mode. Review the checks and requester identity before execution. Store captured metadata through the evidence chain after redaction. Generate findings with conservative severity and link every claim to an evidence reference.

The Go package `src/c2/orchestrator` provides an authorization gate for this workflow. A request is allowed only when the engagement is marked authorized, the Rules of Engagement identifier and validity window are present, the requester is identified, the target and technique are allowlisted, and active actions carry an explicit approval ID. Integrate `EngagementScope.Authorize` before dispatching any assessment or simulator task; do not bypass it in adapters.

The remaining integration work is intentionally repository-owner work: persist scopes and approvals in the selected database, bind approvals to authenticated operator identities, append allow/deny decisions to the evidence ledger, and connect the gate to every dispatch path. Until those integrations are complete, the scope type is a tested policy primitive rather than a claim of end-to-end enforcement.

## Evidence handling

Export the evidence chain and its manifest together. Verify sequence numbers, parent links, record hashes, and timestamps before sharing it; a failed verification is an incident, not a recoverable warning. Store backups with restricted permissions and apply the documented retention period. Treat the SQLite database, reports, and manifests as sensitive engagement artifacts.

## Incident response

If an out-of-scope asset is observed, stop the assessment, preserve the current evidence chain, notify the emergency contact, and record the decision in the audit log. If a secret is exposed, revoke and rotate it immediately. Do not delete or alter forensic records to conceal activity.

## Shutdown and review

Expire the engagement, disable temporary credentials, archive the final report and manifest, remove only approved temporary lab resources, and conduct a human review of findings and false positives. No production action is authorized by this runbook without the system owner's approval.
