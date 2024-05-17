package repo_adapter

import (
	repository "annotater/internal/bl/annotationService/annotattionRepo"
	"annotater/internal/models"
	models_da "annotater/internal/models/modelsDA"
	cache_utils "annotater/internal/pkg/cacheUtils"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	ErrNothingDelete  = errors.New("nothing was deleted")
	userIDCachePrefix = "user_id"
)

type AnotattionRepositoryAdapter struct {
	db    *gorm.DB
	cache cache_utils.ICache
}

func NewAnotattionRepositoryAdapter(srcDB *gorm.DB, cacheSrc cache_utils.ICache) repository.IAnotattionRepository {
	return &AnotattionRepositoryAdapter{
		db:    srcDB,
		cache: cacheSrc,
	}
}

func (repo *AnotattionRepositoryAdapter) AddAnottation(markUp *models.Markup) error {
	markUpDa, err := models_da.ToDaMarkup(*markUp)
	if err != nil {
		return errors.Wrap(err, "Error in getting anotattion type")
	}
	tx := repo.db.Create(markUpDa)
	if tx.Error == gorm.ErrForeignKeyViolated {
		return models.ErrViolatingKeyAnnot
	}
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Error in adding anotattion")
	}
	return nil
}

func (repo *AnotattionRepositoryAdapter) DeleteAnotattion(id uint64) error { // do we need transactions here?
	strID := strconv.FormatUint(id, 10)
	err := repo.cache.Del(strID)
	if err != nil {
		return errors.Wrap(err, "Error in deleting anotattion cache")
	}
	tx := repo.db.Where("id = ?", id) //using that because if id is equal to 0 then the first found row will be deleted
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Error in deleting anotattion")
	}
	fmt.Print(tx.Error)

	tx = tx.Delete(&models_da.Markup{})
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Error in deleting anotattion")
	}
	/*if tx.RowsAffected == 0 {
		return ErrNothingDelete TODO:: think wether it must be and error
	}*/
	return nil
}

func (repo *AnotattionRepositoryAdapter) GetAnottationByID(id uint64) (*models.Markup, error) {
	var markUpDA models_da.Markup
	//strID := strconv.FormatUint(id, 10)
	//if err := repo.cache.Get(strID, &markUpDA); err != nil {

	tx := repo.db.Where("id = ?", id).First(&markUpDA)

	if tx.Error == gorm.ErrRecordNotFound {
		return nil, models.ErrNotFound
	}

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Error in getting anotattion type")
	}
	//	}
	/*	err := repo.cache.Set(strID, &markUpDA)
		if err != nil {
			return nil, err
		}*/

	markUpType, err := models_da.FromDaMarkup(&markUpDA)
	if err != nil {
		return nil, errors.Wrap(err, "Error in getting anotattion type")
	}
	return &markUpType, nil
}
func (repo *AnotattionRepositoryAdapter) GetAnottationsByUserID(id uint64) ([]models.Markup, error) {
	var markUpsDA []models_da.Markup
	tx := repo.db.Where("creator_id = ?", id).Find(&markUpsDA)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Error in getting anotattion type")
	}
	markUps, err := models_da.FromDaMarkupSlice(markUpsDA)
	if err != nil {
		return nil, errors.Wrap(err, "Error in getting markups by userID")
	}
	return markUps, err
}

func (repo *AnotattionRepositoryAdapter) GetAllAnottations() ([]models.Markup, error) {
	var markUpsDA []models_da.Markup

	tx := repo.db.Find(&markUpsDA)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Error in getting anotattion type")
	}

	markUps, err := models_da.FromDaMarkupSlice(markUpsDA)
	if err != nil {
		return nil, errors.Wrap(err, "Error in getting all markups")
	}
	return markUps, err
}

func (repo *AnotattionRepositoryAdapter) UpdateAnotattion(id uint64, markUp *models.Markup) error {
	markUpDa, err := models_da.ToDaMarkup(*markUp)

	strID := strconv.FormatUint(id, 10)
	if err != nil {
		return errors.Wrap(err, "Error in updating anotattion")
	}
	tx := repo.db.Where("id = ?", id).Updates(markUpDa)

	if tx.Error == gorm.ErrRecordNotFound {
		return models.ErrNotFound
	}
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Error in updating anotattion")
	}
	err = repo.cache.Set(strID, markUpDa)
	if err != nil {
		return err
	}
	return nil
}
