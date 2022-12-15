package storage

import (
	"os"
	"sync"

	"github.com/jordyf15/thullo-api/models"
)

type Storage interface {
	UploadFile(respond chan<- error, wg *sync.WaitGroup, currentImage *models.Image, file *os.File, metadata map[string]string)
}
