package unsplash

import "os"

type Repository interface {
	GetImagesForID(photoID string, focalPointY float64) ([]*os.File, error)
}
