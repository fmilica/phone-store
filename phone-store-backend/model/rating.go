package model

type Rating struct {
	Id        string `json:"id"`
	DisplayId string `json:"displayId"`
	ParentId  string `json:"parentId"`
	Mark      int    `json:"mark"`
}
