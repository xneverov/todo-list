package models

type Task struct {
	ID          string `json:"id"`
	Date        string `json:"date"`
	Title       string `json:"title"`
	Description string `json:"comment"`
	Repeat      string `json:"repeat"`
}
