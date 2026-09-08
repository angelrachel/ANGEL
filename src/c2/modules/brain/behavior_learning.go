package brain

import "time"

type LearningRecord struct {
Action    string
Success   bool
Timestamp string
}

type BehaviorLearning struct {
History []LearningRecord
}

func NewBehaviorLearning() *BehaviorLearning {
return &BehaviorLearning{History: []LearningRecord{}}
}

func (b *BehaviorLearning) Learn(action string, success bool) {
b.History = append(b.History, LearningRecord{Action: action, Success: success, Timestamp: time.Now().UTC().Format(time.RFC3339)})
}

func (b *BehaviorLearning) GetBestAction() string {
successCount := 0
for _, record := range b.History {
if record.Success {
successCount++
}
}
if successCount > 0 {
return "continue_attack_path"
}
return "recon_adjust"
}

func (b *BehaviorLearning) GetHistoryCount() int {
return len(b.History)
}
