package service

import (
	"errors"
	"phone-store-backend/model"
	"phone-store-backend/repository"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type PhoneService interface {
	Validate(phone *model.Phone) error
	Create(phone *model.Phone) (*model.Phone, error)
	FindAll() ([]model.Phone, error)
	DeleteAll()
}

type phoneService struct{}

var (
	phoneRepo repository.PhoneRepository
)

func NewPhoneService(repo repository.PhoneRepository) PhoneService {
	phoneRepo = repo
	return &phoneService{}
}

func (*phoneService) Validate(phone *model.Phone) error {

	if phone == nil {
		err := errors.New("The phone is empty.")
		return err
	}

	// if phone.Make == "" {
	// 	err := errors.New("The make is empty.")
	// 	return err
	// }

	return nil
}

func (*phoneService) Create(phone *model.Phone) (*model.Phone, error) {

	phone.Id = uuid.New().String()

	return phoneRepo.Save(phone)
}

func (*phoneService) FindAll() ([]model.Phone, error) {

	return phoneRepo.FindAll()
}

func (*phoneService) DeleteAll() {

	phoneRepo.DeleteAll()
}
