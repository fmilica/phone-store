package model

import "time"

type DisplaySearchDTO struct {
	Brand     string `json:"brand"`
	Processor string `json:"processor"`
	Battery   string `json:"battery"`
	RAM       int    `json:"ram"`

	PriceFrom int       `json:"priceFrom"`
	PriceTo   int       `json:"priceTo"`
	DateFrom  time.Time `json:"dateFrom"`
	DateTo    time.Time `json:"dateTo"`

	/*
		'average rate'
		'price up'
		'price down'
		'oldest'
		'latest'
	*/
	Sort string `json:"sort"`
}
