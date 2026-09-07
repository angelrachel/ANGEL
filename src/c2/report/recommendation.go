package report

import (
"time"
)

type Recommendation struct {
Title       string
Description string
Priority    string
}

func NewRecommendation() *Recommendation {
return &Recommendation{}
}

func (r *Recommendation) SetTitle(title string) {
r.Title = title
}

func (r *Recommendation) SetDescription(description string) {
r.Description = description
}

func (r *Recommendation) SetPriority(priority string) {
r.Priority = priority
}

func (r *Recommendation) GetTitle() string {
return r.Title
}

func (r *Recommendation) GetDescription() string {
return r.Description
}

func (r *Recommendation) GetPriority() string {
return r.Priority
}

func (r *Recommendation) GetLastSeen() time.Time {
return time.Now()
}
