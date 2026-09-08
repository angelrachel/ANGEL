package evidence

type ParentChild struct {
ParentID string
ChildID  string
}

func NewParentChild(parentID, childID string) *ParentChild {
return &ParentChild{
ParentID: parentID,
ChildID:  childID,
}
}

func (p *ParentChild) GetParentID() string {
return p.ParentID
}

func (p *ParentChild) GetChildID() string {
return p.ChildID
}
