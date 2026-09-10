package platform

import "testing"

func TestCurrentManifestIsComplete(t *testing.T) {
	manifest := Current()
	if err := manifest.Validate(); err != nil {
		t.Fatal(err)
	}
	if manifest.Name != "ANGEL" || manifest.LayerCount != 25 {
		t.Fatalf("unexpected manifest: %#v", manifest)
	}
}

func TestLayerBounds(t *testing.T) {
	if name, ok := Layer(1); !ok || name == "" {
		t.Fatal("expected first layer")
	}
	if _, ok := Layer(0); ok {
		t.Fatal("layer zero must be rejected")
	}
	if _, ok := Layer(26); ok {
		t.Fatal("layer 26 must be rejected")
	}
}

func TestDeniedCapabilitiesFailClosed(t *testing.T) {
	for _, capability := range []string{"arbitrary-command", "data-exfiltration", "PROCESS-INJECTION"} {
		if !IsDeniedCapability(capability) {
			t.Fatalf("expected denied capability: %s", capability)
		}
	}
	if IsDeniedCapability("surface-map") {
		t.Fatal("safe assessment capability must not be denied")
	}
}
