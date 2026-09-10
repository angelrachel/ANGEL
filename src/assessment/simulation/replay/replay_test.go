package replay

import "testing"

func TestReplayAcceptsAuthorizedFixtureResult(t *testing.T) {
	task := []byte(`{"id":"task-1","agent_type":"recon","target_ref":"fixture://lab/web-app-01","technique":"http-header-observation","mode":"simulate","requested_by":"operator","approval_id":null}`)
	result := []byte(`{"task_id":"task-1","status":"completed","simulation":true,"authorized":true,"message":"fixture observation completed","evidence_refs":["fixture://evidence/1"]}`)
	got, err := Replay(task, result)
	if err != nil {
		t.Fatalf("expected replay to succeed: %v", err)
	}
	if got.TaskID != "task-1" || got.Status != "completed" {
		t.Fatalf("unexpected replay result: %#v", got)
	}
}

func TestReplayRejectsNonFixtureTarget(t *testing.T) {
	task := []byte(`{"id":"task-1","agent_type":"recon","target_ref":"https://example.test","technique":"http-header-observation","mode":"simulate","requested_by":"operator"}`)
	result := []byte(`{"task_id":"task-1","status":"completed","simulation":true,"authorized":true,"evidence_refs":[]}`)
	if _, err := Replay(task, result); err == nil {
		t.Fatal("expected non-fixture target to be rejected")
	}
}

func TestReplayRejectsExecutionModeAndUnauthorizedResult(t *testing.T) {
	task := []byte(`{"id":"task-1","agent_type":"recon","target_ref":"fixture://lab/web-app-01","technique":"http-header-observation","mode":"active","requested_by":"operator"}`)
	result := []byte(`{"task_id":"task-1","status":"completed","simulation":true,"authorized":true,"evidence_refs":[]}`)
	if _, err := Replay(task, result); err == nil {
		t.Fatal("expected active mode to be rejected")
	}

	task = []byte(`{"id":"task-1","agent_type":"recon","target_ref":"fixture://lab/web-app-01","technique":"http-header-observation","mode":"simulate","requested_by":"operator"}`)
	result = []byte(`{"task_id":"task-1","status":"completed","simulation":true,"authorized":false,"evidence_refs":[]}`)
	if _, err := Replay(task, result); err == nil {
		t.Fatal("expected unauthorized result to be rejected")
	}
}

func TestReplayRejectsMismatchedTaskAndExternalEvidence(t *testing.T) {
	task := []byte(`{"id":"task-1","agent_type":"recon","target_ref":"fixture://lab/web-app-01","technique":"http-header-observation","mode":"simulate","requested_by":"operator"}`)
	result := []byte(`{"task_id":"task-2","status":"completed","simulation":true,"authorized":true,"evidence_refs":["https://example.test/evidence"]}`)
	if _, err := Replay(task, result); err == nil {
		t.Fatal("expected mismatched task to be rejected")
	}
}

func TestReplayRejectsUnknownFieldsAndTrailingJSON(t *testing.T) {
	task := []byte(`{"id":"task-1","agent_type":"recon","target_ref":"fixture://lab/web-app-01","technique":"http-header-observation","mode":"simulate","requested_by":"operator","extra":true}`)
	result := []byte(`{"task_id":"task-1","status":"completed","simulation":true,"authorized":true,"evidence_refs":[]}`)
	if _, err := Replay(task, result); err == nil {
		t.Fatal("expected unknown task field to be rejected")
	}

	task = []byte(`{"id":"task-1","agent_type":"recon","target_ref":"fixture://lab/web-app-01","technique":"http-header-observation","mode":"simulate","requested_by":"operator"}`)
	result = []byte(`{"task_id":"task-1","status":"completed","simulation":true,"authorized":true,"evidence_refs":[]} {}`)
	if _, err := Replay(task, result); err == nil {
		t.Fatal("expected trailing JSON to be rejected")
	}
}
