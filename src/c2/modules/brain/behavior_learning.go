package brain

type BehaviorLearning struct {
History     []string
SuccessRate map[string]float64
}

func NewBehaviorLearning() *BehaviorLearning {
return &BehaviorLearning{
History:     make([]string, 0),
SuccessRate: make(map[string]float64),
}
}

func (b *BehaviorLearning) Learn(action string, success bool) {
b.History = append(b.History, action)
if _, ok := b.SuccessRate[action]; !ok {
b.SuccessRate[action] = 0.0
}
if success {
b.SuccessRate[action] = b.SuccessRate[action] + 0.1
} else {
b.SuccessRate[action] = b.SuccessRate[action] - 0.05
}
}

func (b *BehaviorLearning) GetBestAction(actions []string) string {
var bestAction string
var bestRate float64
for _, action := range actions {
if rate, ok := b.SuccessRate[action]; ok && rate > bestRate {
bestRate = rate
bestAction = action
}
}
return bestAction
}

func (b *BehaviorLearning) Adapt() map[string]string {
adaptations := make(map[string]string)
for action, rate := range b.SuccessRate {
if rate < 0.3 {
adaptations[action] = "avoid"
} else if rate > 0.8 {
adaptations[action] = "prioritize"
}
}
return adaptations
}
