package model

type Statistics struct {
	Submitted  int `json:"submitted"`
	Completed  int `json:"completed"`
	InProgress int `json:"in_progress"`
}
