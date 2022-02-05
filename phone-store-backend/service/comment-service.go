package service

import (
	"errors"
	"phone-store-backend/model"
	"phone-store-backend/repository"

	"github.com/google/uuid"
)

type CommentService interface {
	Validate(comment *model.Comment) error
	Create(comment *model.Comment) (*model.Comment, error)
	FindAll() ([]model.Comment, error)
	// DeleteAll()
}

type commentService struct{}

var (
	commentRepo repository.CommentRepository
)

func NewCommentService(repo repository.CommentRepository) CommentService {
	commentRepo = repo
	return &commentService{}
}

func (*commentService) Create(comment *model.Comment) (*model.Comment, error) {

	comment.Id = uuid.New().String()

	return commentRepo.Save(comment)
}

func (*commentService) Validate(comment *model.Comment) error {

	if comment == nil {
		err := errors.New("the comment is empty.")
		return err
	}

	if comment.DisplayId == "" {
		err := errors.New("the phone display id is empty.")
		return err
	}

	if comment.Content == "" {
		err := errors.New("the content is empty.")
		return err
	}

	return nil
}

func (*commentService) FindAll() ([]model.Comment, error) {

	return commentRepo.FindAll()
}
