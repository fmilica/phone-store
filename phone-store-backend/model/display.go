package model

import "time"

type Display struct {
	Id       string    `json:"id"`
	Phone    Phone     `json:"phone"`
	Price    int       `json:"price"`
	Date     time.Time `json:"date"` // publish date
	Ratings  []Rating  `json:"ratings"`
	Comments []Comment `json:"comments"`
}
