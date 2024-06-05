package models_dto // stands for data_transfer_objec

import (
	"annotater/internal/models"
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID            uuid.UUID          `json:"id"`
	PageCount     int                `json:"page_count"`
	DocumentName  string             `json:"document_name"` //pdf file -- the whole file
	ChecksCount   uint64             `json:"checks_count"`
	CreatorID     uint64             `json:"creator_id"`
	CreationTime  time.Time          `json:"creation_time"`
	HasPassed     bool               `json:"has_passed"`
	CheckedStatus models.CheckStatus `json:"check_status"`
}

func FromDtoDocument(document *Document) models.DocumentMetaData {

	doc := models.DocumentMetaData{ //TODO::Think about changing only the pointer
		ID:            document.ID,
		PageCount:     document.PageCount,
		DocumentName:  document.DocumentName,
		CreatorID:     document.CreatorID,
		CreationTime:  document.CreationTime,
		HasPassed:     document.HasPassed,
		ChecksCount:   document.ChecksCount,
		CheckedStatus: document.CheckedStatus,
	}
	return doc

}

func ToDtoDocument(document models.DocumentMetaData) Document {
	dtoDoc := Document{
		ID:            document.ID,
		PageCount:     document.PageCount,
		DocumentName:  document.DocumentName,
		CreatorID:     document.CreatorID,
		CreationTime:  document.CreationTime,
		HasPassed:     document.HasPassed,
		ChecksCount:   document.ChecksCount,
		CheckedStatus: document.CheckedStatus,
	}
	return dtoDoc
}
