# ANGEL Operations Runbook

## Before an engagement

Confirm a signed authorization, an active Rules of Engagement record, the target inventory, the emergency contact, the data-retention period, and the approved assessment mode. Use synthetic data in development and staging. Production secrets must come from a managed secret store.

## Start-up checks

Run the test suite and security checks before deployment. In production or staging, configure non-placeholder operator, shared, and authentication secrets. Bind the service to a private interface and place it behind the hardened reverse proxy. Confirm the health and readiness endpoints, database permissions, and backup destination.

## Assessment workflow

Create the scope with exact hosts and paths. Create an auditable assessment plan using only passive or simulation mode. Review the checks and requester identity before execution. Store captured metadata through the evidence chain after redaction. Generate findings with conservative severity and link every claim to an evidence reference.

## Evidence handling

Export the evidence chain and its manifest together. Verify the chain before sharing it. Store backups with restricted permissions and apply the documented retention period. Treat the SQLite database, reports, and manifests as sensitive engagement artifacts.

## Incident response

If an out-of-scope asset is observed, stop the assessment, preserve the current evidence chain, notify the emergency contact, and record the decision in the audit log. If a secret is exposed, revoke and rotate it immediately. Do not delete or alter forensic records to conceal activity.

## Shutdown and review

Expire the engagement, disable temporary credentials, archive the final report and manifest, remove only approved temporary lab resources, and conduct a human review of findings and false positives. No production action is authorized by this runbook without the system owner's approval.
