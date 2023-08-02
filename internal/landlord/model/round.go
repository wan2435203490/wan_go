package model

// 游戏回合
type Round struct {
	Multiple int `json:"multiple"`
	StepNum  int `json:"stepNum"`
}

func NewRound() *Round {
	return &Round{
		Multiple: 1,
		StepNum:  0,
	}
}
