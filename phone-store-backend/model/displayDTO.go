package model

import "time"

type DisplayDTO struct {
	Brand     string    `json:"brand"`
	Model     string    `json:"model"`
	Date      time.Time `json:"date"`
	Processor string    `json:"processor"`
	Battery   string    `json:"battery"`
	RAM       int       `json:"ram"`

	Price int `json:"price"`
}
