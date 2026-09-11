package checks

import "testing"

func TestCatalogContainsStableLayeredDescriptors(t *testing.T) {
	catalog := Catalog()
	if len(catalog) < 20 {
		t.Fatalf("catalog too small: %d", len(catalog))
	}
	for i, descriptor := range catalog {
		if descriptor.ID == "" || descriptor.Layer < 1 || descriptor.Layer > 25 {
			t.Fatalf("invalid descriptor at %d: %+v", i, descriptor)
		}
		if i > 0 && (catalog[i-1].Layer > descriptor.Layer || (catalog[i-1].Layer == descriptor.Layer && catalog[i-1].ID > descriptor.ID)) {
			t.Fatal("catalog is not stably sorted")
		}
	}
}
