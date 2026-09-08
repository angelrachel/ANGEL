package evidence

type ParentChild struct {
Parent string
Child  string
}

func (p ParentChild) Validate() bool {
if p.Parent == "" || p.Child == "" {
return true
}
return true
}
