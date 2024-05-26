package repo_adapter

import (
	repository "annotater/internal/bl/documentService/documentMetaDataRepo"
	"annotater/internal/models"
	models_da "annotater/internal/models/modelsDA"
	cache_utils "annotater/internal/pkg/cacheUtils"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	creatorIDPrefix = "crator_id_"
)

type DocumentMetaDataRepositoryAdapter struct {
	db    *gorm.DB
	cache cache_utils.ICache
}

func NewDocumentRepositoryAdapter(srcDB *gorm.DB, cacheSrc cache_utils.ICache) repository.IDocumentMetaDataRepository {
	return &DocumentMetaDataRepositoryAdapter{
		db:    srcDB,
		cache: cacheSrc,
	}
}

func (repo *DocumentMetaDataRepositoryAdapter) AddDocument(doc *models.DocumentMetaData) error {

	tx := repo.db.Create(models_da.ToDaDocument(*doc))
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Error in updating document")
	}
	return nil
}

func (repo *DocumentMetaDataRepositoryAdapter) GetDocumentByID(id uuid.UUID) (*models.DocumentMetaData, error) {
	var documentDa models_da.Document
	idReq := id.String()
	err := repo.cache.Get(idReq, &documentDa)
	if err != nil && err != models.ErrNotFound {
		return nil, err
	}
	if err := repo.cache.Get(idReq, &documentDa); err != nil {

		documentDa.ID = id
		tx := repo.db.Where("id = ?", id).First(&documentDa)
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		if tx.Error != nil {
			return nil, errors.Wrap(tx.Error, "Error getting document by ID")
		}
	}
	if err != nil && err != models.ErrNotFound {
		return nil, err

	}
	err = repo.cache.Set(idReq, documentDa)
	if err != nil {
		return nil, err
	}

	document := models_da.FromDaDocument(&documentDa)
	return &document, nil
}

func (repo *DocumentMetaDataRepositoryAdapter) DeleteDocumentByID(id uuid.UUID) error {

	err := repo.cache.Del(id.String())
	if err != nil {
		return err
	}

	tx := repo.db.Delete(models.DocumentMetaData{}, id) // specifically for gorm
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Error in deleting document")
	}
	return nil
}
func (repo *DocumentMetaDataRepositoryAdapter) GetDocumentsByCreatorID(id uint64) ([]models.DocumentMetaData, error) {
	var documentsDA []models_da.Document
	tx := repo.db.Where("creator_id = ?", id).Find(&documentsDA)
	//strID := strconv.FormatUint(id, 10)
	//err := repo.cache.Get(strID, &documentsDA)
	//if err != nil {

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Error in getting documents by creator")
	}
	//	}
	/*err = repo.cache.Set(strID, documentsDA)
	if err != nil {
		fmt.Print(err.Error()) //  TODO::add logging
	}*/
	documents := models_da.FromDaDocumentSlice(documentsDA)
	return documents, nil
}

func (repo *DocumentMetaDataRepositoryAdapter) GetDocumentCountByCreator(id uint64) (int64, error) {
	var count int64
	tx := repo.db.Model(models_da.Document{}).Where("creator_id = ?", id).Count(&count)

	if tx.Error != nil {
		return -1, errors.Wrap(tx.Error, "Error in getting count by creator")
	}
	return count, nil
}

func (repo *DocumentMetaDataRepositoryAdapter) GetDocumentCountByDocumentName(docName string) (int64, error) {
	var count int64
	tx := repo.db.Model(models_da.Document{}).Where("document_name = ?", docName).Count(&count)

	if tx.Error != nil {
		return -1, errors.Wrap(tx.Error, "Error in getting count by creator")
	}
	return count, nil
}

func (repo *DocumentMetaDataRepositoryAdapter) GetDocumentNotChecked() (*models.DocumentMetaData, error) {
	var documentsDA models_da.Document
	tx := repo.db.Model(models_da.Document{}).Where("checked_status  = ?", models.NotChecked).First(&documentsDA)

	if tx.Error == gorm.ErrRecordNotFound {
		return nil, models.ErrNotFound
	}
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Error in getting document not checked")
	}

	document := models_da.FromDaDocument(&documentsDA)
	return &document, nil
}

func (repo *DocumentMetaDataRepositoryAdapter) UpdateData(id uuid.UUID, data models.DocumentMetaData) error {
	dataDA := models_da.ToDaDocument(data)

	idReq := id.String()
	err := repo.cache.Del(id.String())
	if err != nil {
		return err
	}

	tx := repo.db.Model(models_da.Document{}).Where("id = ?", id).Updates(dataDA)

	if tx.Error == gorm.ErrRecordNotFound {
		return models.ErrNotFound
	}
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Error in getting count by creator")
	}
	err = repo.cache.Set(idReq, dataDA)
	if err != nil {
		return err
	}
	return nil
}
