package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"phone-store-backend/model"
	"phone-store-backend/service"
)

type RatingController interface {
	GetAll(response http.ResponseWriter, request *http.Request)
	Save(response http.ResponseWriter, request *http.Request)
}

type ratingController struct{}

var (
	ratingService service.RatingService
)

func NewRatingController(service service.RatingService) RatingController {
	ratingService = service
	return &ratingController{}
}

func (*ratingController) Save(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("Content-Type", "application/json")

	var rating model.Rating

	err := json.NewDecoder(request.Body).Decode(&rating)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(model.ErrorMessage{Message: "Error unmarshaling data"})
		fmt.Println("Error 1 rating")
		return
	}

	err1 := ratingService.Validate(&rating)

	if err1 != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(model.ErrorMessage{Message: err1.Error()})
		fmt.Println(err1.Error())
		return
	}

	result, err2 := ratingService.Create(&rating)

	if err2 != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(model.ErrorMessage{Message: "Error saving the rating."})
		return
	}

	response.WriteHeader((http.StatusOK))
	json.NewEncoder(response).Encode(result)
}

func (*ratingController) GetAll(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("Content-Type", "application/json")

	ratings, err := ratingService.FindAll()

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(model.ErrorMessage{Message: "Error getting the ratings."})
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(ratings)
}
