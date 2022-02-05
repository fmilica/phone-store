package model

type Comment struct {
	Id        string    `json:"id"`
	DisplayId string    `json:"displayId"`
	ParentId  string    `json:"parentId"`
	Content   string    `json:"content"`
	Comments  []Comment `json:"comments"`
	Ratings   []Rating  `json:"ratings"`
}
