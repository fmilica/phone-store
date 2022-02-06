package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"phone-store-backend/model"
	"phone-store-backend/service"
)

type CommentController interface {
	GetAll(response http.ResponseWriter, request *http.Request)
	Save(response http.ResponseWriter, request *http.Request)
}

type commentController struct{}

var (
	commentService service.CommentService
)

func NewCommentController(service service.CommentService) CommentController {
	commentService = service
	return &commentController{}
}

func (*commentController) Save(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("Content-Type", "application/json")

	var comment model.Comment

	err := json.NewDecoder(request.Body).Decode(&comment)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(model.ErrorMessage{Message: "Error unmarshaling data"})
		fmt.Println("Error 1 comment")
		return
	}

	err1 := commentService.Validate(&comment)

	if err1 != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(model.ErrorMessage{Message: err1.Error()})
		fmt.Println(err1.Error())
		return
	}

	result, err2 := commentService.Create(&comment)

	if err2 != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(model.ErrorMessage{Message: "Error saving the comment."})
		return
	}

	response.WriteHeader((http.StatusOK))
	json.NewEncoder(response).Encode(result)
}

func (*commentController) GetAll(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("Content-Type", "application/json")

	comments, err := commentService.FindAll()

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(model.ErrorMessage{Message: "Error getting the comments"})
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(comments)
}
