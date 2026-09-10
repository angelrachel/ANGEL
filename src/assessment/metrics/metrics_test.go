package metrics

import (
	"testing"
	"time"
)

func TestCollectorTracksRequestsAndErrors(t *testing.T) {
	collector := New()
	started := time.Now().Add(-time.Millisecond)
	collector.Start()
	collector.Finish(200, started)
	collector.Start()
	collector.Finish(500, started)
	snapshot := collector.Snapshot()
	if snapshot.Requests != 2 || snapshot.Errors != 1 || snapshot.Active != 0 || snapshot.TotalLatencyMS == 0 {
		t.Fatalf("snapshot=%#v", snapshot)
	}
}
