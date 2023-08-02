package dto

type Chat struct {
	Content   string `json:"content"`
	Type      int    `json:"type"`
	Dimension string `json:"dimension"`
}
