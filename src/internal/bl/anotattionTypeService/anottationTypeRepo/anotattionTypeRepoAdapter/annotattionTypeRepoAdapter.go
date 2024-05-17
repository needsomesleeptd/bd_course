package repo_adapter

import (
	repository "annotater/internal/bl/anotattionTypeService/anottationTypeRepo"
	"annotater/internal/models"
	models_da "annotater/internal/models/modelsDA"
	cache_utils "annotater/internal/pkg/cacheUtils"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type AnotattionTypeRepositoryAdapter struct {
	db    *gorm.DB
	cache cache_utils.ICache
}

func NewAnotattionTypeRepositoryAdapter(srcDB *gorm.DB, cacheSrc *cache_utils.ICache) repository.IAnotattionTypeRepository {
	return &AnotattionTypeRepositoryAdapter{
		db:    srcDB,
		cache: *cacheSrc,
	}
}

func (repo *AnotattionTypeRepositoryAdapter) AddAnottationType(markUp *models.MarkupType) error {
	tx := repo.db.Create(models_da.ToDaMarkupType(*markUp))
	if tx.Error == gorm.ErrDuplicatedKey {
		return models.ErrDuplicateMarkupType
	}
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Error in adding anotattion type db")
	}
	return nil
}

func (repo *AnotattionTypeRepositoryAdapter) DeleteAnotattionType(id uint64) error { //note that it is cascade deletion, gorm doesn't support cascade deletion((

	err := repo.db.Transaction(func(tx *gorm.DB) error {

		err := tx.Where("class_label = ?", id).Delete(&models_da.Markup{}).Error
		if err != nil {
			return errors.Wrap(tx.Error, "Error in deleting anotattion type")
		}

		err = tx.Where("id = ?", id).Delete(&models_da.MarkupType{}).Error
		if err != nil {
			return errors.Wrap(tx.Error, "Error in deleting anotattion type db")
		}
		return nil
	})
	if err != nil {
		return err
	}
	strID := strconv.FormatUint(id, 10)
	err = repo.cache.Del(strID)
	if err != nil {
		fmt.Print(err)
	}
	return err
}

func (repo *AnotattionTypeRepositoryAdapter) GetAnottationTypeByID(id uint64) (*models.MarkupType, error) {
	var markUpTypeDA models_da.MarkupType
	markUpTypeDA.ID = id

	strId := strconv.FormatUint(id, 10)
	err := repo.cache.Get(strId, &markUpTypeDA)
	if err != nil {

		tx := repo.db.Where("id = ?", id).First(&markUpTypeDA)
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		if tx.Error != nil {
			return nil, errors.Wrap(tx.Error, "Error in getting anotattion type db")
		}
	}
	repo.cache.Set(strId, markUpTypeDA)
	markUpType := models_da.FromDaMarkupType(&markUpTypeDA)
	return &markUpType, nil
}

func (repo *AnotattionTypeRepositoryAdapter) GetAnottationTypesByIDs(ids []uint64) ([]models.MarkupType, error) {
	var markUpTypesDA []models_da.MarkupType

	tx := repo.db.Find(&markUpTypesDA, ids) // works only when the primary key is set and is a valid ID
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Error in getting anotattion type db")
	}
	markUpTypes := models_da.FromDaMarkupTypeSlice(markUpTypesDA)
	return markUpTypes, nil
}

func (repo *AnotattionTypeRepositoryAdapter) GetAnottationTypesByUserID(creator_id uint64) ([]models.MarkupType, error) {
	var markUpsTypeDA []models_da.MarkupType
	tx := repo.db.Where("creator_id = ?", creator_id).Find(&markUpsTypeDA)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Error in getting anotattion type db")
	}
	markUpTypes := models_da.FromDaMarkupTypeSlice(markUpsTypeDA)

	return markUpTypes, nil
}

func (repo *AnotattionTypeRepositoryAdapter) GetAllAnottationTypes() ([]models.MarkupType, error) {
	var markUpsTypeDA []models_da.MarkupType
	tx := repo.db.Find(&markUpsTypeDA)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Error in getting anotattion type db")
	}
	markUpTypes := models_da.FromDaMarkupTypeSlice(markUpsTypeDA)

	return markUpTypes, nil
}
