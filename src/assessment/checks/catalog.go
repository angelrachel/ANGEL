package checks

import "sort"

type Descriptor struct {
	ID    string `json:"id"`
	Layer int    `json:"layer"`
}

func Catalog() []Descriptor {
	out := make([]Descriptor, 0, len(Registry()))
	for _, check := range Registry() {
		out = append(out, Descriptor{ID: check.ID(), Layer: check.Layer()})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Layer == out[j].Layer {
			return out[i].ID < out[j].ID
		}
		return out[i].Layer < out[j].Layer
	})
	return out
}
