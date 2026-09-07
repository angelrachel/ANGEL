package cred_crack

import (
"strings"
)

type RuleEngine struct {
Rules []string
}

func NewRuleEngine() *RuleEngine {
return &RuleEngine{
Rules: []string{
"$1", "$2", "$3", "$!",
"$1$2$3", "$2023", "$2024",
"@", "a",
"0", "o",
"1", "l",
},
}
}

func (r *RuleEngine) AddRule(rule string) {
r.Rules = append(r.Rules, rule)
}

func (r *RuleEngine) ApplyRules(word string) []string {
var results []string
for _, rule := range r.Rules {
results = append(results, word+rule)
results = append(results, rule+word)
results = append(results, strings.Title(word)+rule)
}
return results
}

func (r *RuleEngine) ApplyAllRules(words []string) []string {
var results []string
for _, word := range words {
results = append(results, r.ApplyRules(word)...)
}
return results
}

func (r *RuleEngine) GetRules() []string {
return r.Rules
}

func (r *RuleEngine) GetRuleCount() int {
return len(r.Rules)
}

func (r *RuleEngine) GenerateHashcatRuleFile() string {
var sb strings.Builder
for _, rule := range r.Rules {
sb.WriteString(rule + "\n")
}
return sb.String()
}

func (r *RuleEngine) SetRules(rules []string) {
r.Rules = rules
}

func (r *RuleEngine) ClearRules() {
r.Rules = []string{}
}

func (r *RuleEngine) GetRule(index int) string {
if index < 0 || index >= len(r.Rules) {
return ""
}
return r.Rules[index]
}
