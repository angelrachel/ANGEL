package evidence

type ParentChild struct {
	Parent string
	Child  string
}

func (p ParentChild) Validate() bool {
	return p.Parent != "" && p.Child != "" && p.Parent != p.Child
}
