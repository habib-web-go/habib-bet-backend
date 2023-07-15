package forms

import "time"

type PublicContest struct {
	ID            uint             `json:"id"`
	Name          string           `json:"name"`
	Start         time.Time        `json:"start"`
	End           time.Time        `json:"end"`
	UserCount     int64            `json:"user_count"`
	Questions     *[]QuestionState `json:"questions,omitempty"`
	QuestionCount int              `json:"question_count"`
}

type QuestionState struct {
	ID         uint      `json:"id"`
	OptionA    string    `json:"option_a,omitempty"`
	OptionB    string    `json:"option_b,omitempty"`
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
	Order      uint      `json:"order"`
	Input      int64     `json:"input"`
	Output     int64     `json:"output"`
	UserAnswer string    `json:"user_answer,omitempty"`
	Answer     string    `json:"answer,omitempty"`
}

type Contest struct {
	PublicContest
	Registered bool `json:"registered"`
	RewardPaid bool `json:"reward_paid"`
}
