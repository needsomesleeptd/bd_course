package service

import (
	repository "annotater/internal/bl/annotationService/annotattionRepo"
	"annotater/internal/models"
	"bytes"
	"image"
	"image/png"

	"github.com/pkg/errors"
)

const (
	ADDING_ANNOT_ERR_STR   = "Error in adding anotattion"
	DELETING_ANNOT_ERR_STR = "Error in deleting anotattion"
	GETTING_ANNOT_ERR_STR  = "Error in getting anotattion"
)

var (
	ErrBoundingBoxes   = models.NewUserErr("Invalid markups bounding boxes")
	ErrInvalidFileType = models.NewUserErr("Invalid filetype")
)

type IAnotattionService interface {
	AddAnottation(anotattion *models.Markup) error
	DeleteAnotattion(id uint64) error
	GetAnottationByID(id uint64) (*models.Markup, error)
	GetAnottationByUserID(user_id uint64) ([]models.Markup, error)
	GetAllAnottations() ([]models.Markup, error)
	CheckAnotattion(markup *models.Markup) error
	GetNotCheckedAnnots(count uint64) ([]models.Markup, error)
}

type AnotattionService struct {
	repo repository.IAnotattionRepository
}

func NewAnnotattionService(pRep repository.IAnotattionRepository) IAnotattionService {
	image.RegisterFormat("png", "\x89PNG\r\n\x1a\n", png.Decode, png.DecodeConfig) //for checking file formats
	return &AnotattionService{
		repo: pRep,
	}
}

func AreBBsValid(slice []float32) bool { //TODO:: think if i want to export everything
	if len(slice) != 4 {
		return false
	}
	for _, value := range slice {
		if value < 0.0 || value > 1.0 {
			return false
		}
	}
	return true
}

func CheckPngFile(pngFile []byte) error {
	_, _, err := image.Decode(bytes.NewReader(pngFile))
	if err != nil {
		return err
	}
	return nil

}

func (serv *AnotattionService) AddAnottation(anotattion *models.Markup) error {
	if !AreBBsValid(anotattion.ErrorBB) {
		return ErrBoundingBoxes
	}

	err := CheckPngFile(anotattion.PageData)
	if err != nil {
		return ErrInvalidFileType //maybe user wants to get why his file is broken
	}

	err = serv.repo.AddAnottation(anotattion)
	if err != nil {
		return err
	}
	return err
}

func (serv *AnotattionService) DeleteAnotattion(id uint64) error {
	err := serv.repo.DeleteAnotattion(id)
	if err != nil {
		return errors.Wrap(err, DELETING_ANNOT_ERR_STR)
	}
	return err
}

func (serv *AnotattionService) GetAnottationByID(id uint64) (*models.Markup, error) {
	markup, err := serv.repo.GetAnottationByID(id)
	if err != nil {
		return markup, errors.Wrap(err, GETTING_ANNOT_ERR_STR)
	}
	return markup, err
}

func (serv *AnotattionService) GetAnottationByUserID(user_id uint64) ([]models.Markup, error) {
	markups, err := serv.repo.GetAnottationsByUserID(user_id)
	if err != nil {
		return nil, errors.Wrap(err, GETTING_ANNOT_ERR_STR)
	}
	return markups, nil
}

func (serv *AnotattionService) GetAllAnottations() ([]models.Markup, error) {
	markups, err := serv.repo.GetAllAnottations()
	if err != nil {
		return nil, errors.Wrap(err, GETTING_ANNOT_ERR_STR)
	}
	return markups, nil
}

func (serv *AnotattionService) GetNotCheckedAnnots(count uint64) ([]models.Markup, error) {
	markups, err := serv.repo.GetNotCheckedAnotattions(count)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting not checked markups")
	}
	for _, markup := range markups {
		err = serv.repo.UpdateAnotattion(markup.ID, &models.Markup{CheckedStatus: models.IsBeingChecked})
		if err != nil {
			return nil, errors.Wrap(err, "error in updateing checked markups")
		}
	}
	return markups, nil
}

func (serv *AnotattionService) CheckAnotattion(markup *models.Markup) error {
	if !AreBBsValid(markup.ErrorBB) {
		return ErrBoundingBoxes
	}
	markup.CheckedStatus = models.WasChecked

	err := serv.repo.UpdateAnotattion(markup.ID, markup)
	if err != nil {
		return errors.Wrap(err, "error checking anotattion")
	}
	return nil
}
