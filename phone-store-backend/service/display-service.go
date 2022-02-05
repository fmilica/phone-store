package service

import (
	"errors"
	"phone-store-backend/model"
	"phone-store-backend/repository"
	"time"
)

type DisplayService interface {
	Validate(displayDTO *model.DisplayDTO) error
	Create(display *model.DisplayDTO) (*model.Display, error)
	// FindAll(displaySearchDTO *model.DisplaySearchDTO) ([]model.Display, error)
	FindAll2() ([]model.Display, error)
	// DeleteAll()
}

type displayService struct{}

var (
	displayRepo repository.DisplayRepository
)

func NewDisplayService(repo repository.DisplayRepository) DisplayService {
	displayRepo = repo
	return &displayService{}
}

func (*displayService) Create(displayDTO *model.DisplayDTO) (*model.Display, error) {

	return displayRepo.Save(displayDTO)
}

func (*displayService) Validate(displayDTO *model.DisplayDTO) error {

	if displayDTO == nil {
		err := errors.New("the phone display dto is empty.")
		return err
	}

	if displayDTO.Price == 0 {
		err := errors.New("the phone price is empty.")
		return err
	}

	if displayDTO.Brand == "" {
		err := errors.New("the phone brand is empty.")
		return err
	}

	if displayDTO.Processor == "" {
		err := errors.New("the phone processor is empty.")
		return err
	}

	if displayDTO.Battery == "" {
		err := errors.New("the phone battery is empty.")
		return err
	}

	if displayDTO.RAM == 0 {
		err := errors.New("the phone ram memory is empty.")
		return err
	}

	var zeroTime time.Time

	if displayDTO.Date == zeroTime {
		err := errors.New("the phone production date is empty.")
		return err
	}

	return nil
}

// func (*displayService) FindAll(search *model.Search) ([]model.Display, error) {

// 	return displayRepo.FindAll(search)
// }

func (*displayService) FindAll2() ([]model.Display, error) {

	return displayRepo.FindAll2()
}
