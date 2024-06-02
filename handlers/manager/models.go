package manager

import "chatFileBackend/models"

type UserResource struct {
	Files   []models.MetaData `json:"files"`
	FileNum int64             `json:"filenum"`
}
