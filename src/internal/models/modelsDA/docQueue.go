package models_da //stands for data_acess

import "github.com/google/uuid"

type QueueStatus int

const (
	NotStarted QueueStatus = iota // Role check depends on the order
	Finished
	Error
)

type DocumentQueue struct {
	ID     int64       `gorm:"primaryKey;column:id"`
	DocID  uuid.UUID   `gorm:"column:doc_id"`
	Status QueueStatus `gorm:"column:status"`
}
