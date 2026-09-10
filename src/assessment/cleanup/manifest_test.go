package cleanup

import "testing"

func TestNewSimulationManifest(t *testing.T) {
	manifest, err := NewSimulationManifest("cleanup-1", []string{"fixture://artifact/one", "fixture://artifact/two"})
	if err != nil {
		t.Fatalf("expected manifest to validate: %v", err)
	}
	if manifest.Mode != "simulate" || len(manifest.Entries) != 2 {
		t.Fatalf("unexpected manifest: %#v", manifest)
	}
}

func TestManifestRejectsUnsafeEntries(t *testing.T) {
	cases := []Manifest{
		{ID: "cleanup-1", Mode: "simulate", Entries: []Entry{{PathRef: "/tmp/artifact", Status: "planned"}}},
		{ID: "cleanup-1", Mode: "simulate", Entries: []Entry{{PathRef: "fixture://artifact/one", Status: "deleted"}}},
		{ID: "cleanup-1", Mode: "active", Entries: []Entry{{PathRef: "fixture://artifact/one", Status: "planned"}}},
	}
	for _, manifest := range cases {
		if err := manifest.Validate(); err == nil {
			t.Fatalf("expected manifest to be rejected: %#v", manifest)
		}
	}
}
