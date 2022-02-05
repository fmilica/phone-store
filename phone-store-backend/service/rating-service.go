package service

import (
	"errors"
	"phone-store-backend/model"
	"phone-store-backend/repository"

	"github.com/google/uuid"
)

type RatingService interface {
	Validate(rating *model.Rating) error
	Create(rating *model.Rating) (*model.Rating, error)
	FindAll() ([]model.Rating, error)
	// DeleteAll()
}

type ratingService struct{}

var (
	ratingRepo repository.RatingRepository
)

func NewRatingService(repo repository.RatingRepository) RatingService {
	ratingRepo = repo
	return &ratingService{}
}

func (*ratingService) Create(rating *model.Rating) (*model.Rating, error) {

	rating.Id = uuid.New().String()

	return ratingRepo.Save(rating)
}

func (*ratingService) Validate(rating *model.Rating) error {

	if rating == nil {
		err := errors.New("the rating is empty.")
		return err
	}

	if rating.DisplayId == "" {
		err := errors.New("the phone display id is empty.")
		return err
	}

	if rating.Mark < 0 || rating.Mark > 5 {
		err := errors.New("Mark must be in range 1-5.")
		return err
	}

	return nil
}

func (*ratingService) FindAll() ([]model.Rating, error) {

	return ratingRepo.FindAll()
}
