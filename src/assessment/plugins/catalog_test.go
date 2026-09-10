package plugins

import "testing"

func TestCatalogContains650SafeModulesAcross25Layers(t *testing.T) {
	catalog := Catalog()
	if len(catalog) != 650 {
		t.Fatalf("catalog size=%d", len(catalog))
	}
	seen := map[string]bool{}
	for _, module := range catalog {
		if seen[module.ID] {
			t.Fatalf("duplicate module %s", module.ID)
		}
		seen[module.ID] = true
		if !module.Safe {
			t.Fatalf("unsafe module %s", module.ID)
		}
	}
	if len(Layers()) != 25 {
		t.Fatalf("layer count=%d", len(Layers()))
	}
}
