package ai_orchestrator

import (
"fmt"
"time"
)

type AttackRoute struct {
Target      string
Vector      string
Description string
Priority    int
Timestamp   time.Time
}

type AIRouter struct {
Routes     []AttackRoute
Target     string
LastAction string
}

func NewAIRouter() *AIRouter {
return &AIRouter{
Routes:     []AttackRoute{},
LastAction: "",
}
}

func (a *AIRouter) SetTarget(target string) {
a.Target = target
}

func (a *AIRouter) AssessRoute(vector string, description string, priority int) {
route := AttackRoute{
Target:      a.Target,
Vector:      vector,
Description: description,
Priority:    priority,
Timestamp:   time.Now(),
}
a.Routes = append(a.Routes, route)
}

func (a *AIRouter) SelectBestRoute() AttackRoute {
if len(a.Routes) == 0 {
return AttackRoute{}
}
bestRoute := a.Routes[0]
for _, route := range a.Routes {
if route.Priority > bestRoute.Priority {
bestRoute = route
}
}
return bestRoute
}

func (a *AIRouter) ExecuteRoute(route AttackRoute) string {
a.LastAction = fmt.Sprintf("Executing: %s on target %s", route.Vector, route.Target)
return a.LastAction
}

func (a *AIRouter) GetRoutes() []AttackRoute {
return a.Routes
}

func (a *AIRouter) ClearRoutes() {
a.Routes = []AttackRoute{}
}

func (a *AIRouter) GetTarget() string {
return a.Target
}

func (a *AIRouter) GetLastAction() string {
return a.LastAction
}
