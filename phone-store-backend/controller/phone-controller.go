package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"phone-store-backend/model"
	"phone-store-backend/service"
)

type PhoneController interface {
	GetAll(response http.ResponseWriter, request *http.Request)
	Save(response http.ResponseWriter, request *http.Request)
	DeleteAll(response http.ResponseWriter, request *http.Request)
}

type phoneController struct{}

var (
	phoneService service.PhoneService
)

func NewPhoneController(service service.PhoneService) PhoneController {
	phoneService = service
	return &phoneController{}
}

func (*phoneController) GetAll(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("Content-Type", "application/json")

	phones, err := phoneService.FindAll()

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(model.ErrorMessage{Message: "Error getting the phones"})
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(phones)
}

func (*phoneController) Save(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("Content-Type", "application/json")

	var phone model.Phone

	err := json.NewDecoder(request.Body).Decode(&phone)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(model.ErrorMessage{Message: "Error unmarshaling data"})
		fmt.Println("Greska1")
		return
	}

	err1 := phoneService.Validate(&phone)

	if err1 != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(model.ErrorMessage{Message: err1.Error()})
		fmt.Println("Greska2")
		return
	}

	result, err2 := phoneService.Create(&phone)

	if err2 != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(model.ErrorMessage{Message: "Error saving the phone."})
		return
	}

	response.WriteHeader((http.StatusOK))
	json.NewEncoder(response).Encode(result)
}

func (*phoneController) DeleteAll(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("Content-Type", "application/json")

	phoneService.DeleteAll()

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response)
}
