package models

type Task struct {
	Date        string `json:"date"`
	Title       string `json:"title"`
	Description string `json:"comment"`
	Repeat      string `json:"repeat"`
}
