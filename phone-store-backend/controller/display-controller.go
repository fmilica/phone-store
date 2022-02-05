package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"phone-store-backend/model"
	"phone-store-backend/service"
)

type DisplayController interface {
	// GetAll(response http.ResponseWriter, request *http.Request)
	GetAll2(response http.ResponseWriter, request *http.Request)
	Save(response http.ResponseWriter, request *http.Request)
	// DeleteAll(response http.ResponseWriter, request *http.Request)
}

type displayController struct{}

var (
	displayService service.DisplayService
)

func NewDisplayController(service service.DisplayService) DisplayController {
	displayService = service
	return &displayController{}
}

func (*displayController) Save(response http.ResponseWriter, request *http.Request) {

	fmt.Println("*** Add new display***")

	response.Header().Set("Content-Type", "application/json")

	var displayDTO model.DisplayDTO

	err := json.NewDecoder(request.Body).Decode(&displayDTO)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(model.ErrorMessage{Message: "Error unmarshaling data"})
		fmt.Println("Greska1")
		return
	}

	err1 := displayService.Validate(&displayDTO)

	if err1 != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(model.ErrorMessage{Message: err1.Error()})
		fmt.Println("Greska2")
		return
	}

	result, err2 := displayService.Create(&displayDTO)

	if err2 != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(model.ErrorMessage{Message: "Error saving the phone display."})
		return
	}

	response.WriteHeader((http.StatusOK))
	json.NewEncoder(response).Encode(result)
}

// func (*displayController) GetAll(response http.ResponseWriter, request *http.Request) {

// 	fmt.Println("*** Call GetAll Search Method ***")

// 	response.Header().Set("Content-Type", "application/json")

// 	var search model.Search
// 	fmt.Println(request.Body)

// 	err := json.NewDecoder(request.Body).Decode(&search)

// 	if err != nil {
// 		response.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(response).Encode(model.ErrorMessage{Message: "Error unmarshaling data"})
// 		fmt.Println("Greska1 phone display cont")
// 		return
// 	}

// 	displays, err := displayService.FindAll(&search)

// 	if err != nil {
// 		response.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(response).Encode(model.ErrorMessage{Message: "Error getting the phone displays."})
// 	}

// 	response.WriteHeader(http.StatusOK)
// 	json.NewEncoder(response).Encode(displays)
// }

// /*
// 	Ova metoda je samo dok se razvija klijent.
// 	Samo dobavljamo sve oglase.
// */
func (*displayController) GetAll2(response http.ResponseWriter, request *http.Request) {

	fmt.Println("*** Call GetAll Method ***")

	response.Header().Set("Content-Type", "application/json")

	displays, err := displayService.FindAll2()

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(model.ErrorMessage{Message: "Error getting the phone displays."})
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(displays)
}
