package evidence

type Sequence struct {
Current int
}

func NewSequence() *Sequence {
return &Sequence{
Current: 0,
}
}

func (s *Sequence) Next() int {
s.Current++
return s.Current
}

func (s *Sequence) GetCurrent() int {
return s.Current
}
