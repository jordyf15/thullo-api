package usecase

import (
	"fmt"
	"math"
	"sync"

	"github.com/jordyf15/thullo-api/board"
	"github.com/jordyf15/thullo-api/custom_errors"
	"github.com/jordyf15/thullo-api/models"
	"github.com/jordyf15/thullo-api/storage"
	"github.com/jordyf15/thullo-api/unsplash"
	"github.com/jordyf15/thullo-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type boardUsecase struct {
	boardRepo    board.Repository
	unsplashRepo unsplash.Repository
	storage      storage.Storage
}

func NewBoardUsecase(boardRepo board.Repository, unsplashRepo unsplash.Repository, storage storage.Storage) board.Usecase {
	return &boardUsecase{boardRepo: boardRepo, unsplashRepo: unsplashRepo, storage: storage}
}

func (usecase *boardUsecase) Create(userID primitive.ObjectID, title string, visibility string, cover map[string]interface{}) error {
	errors := []error{}

	source := cover["source"].(string)
	focalPointY := cover["fp_y"].(float64)
	photoID := cover["photo_id"].(string)

	focalPointY = math.Round(focalPointY*1000) / 1000

	if _, exist := board.BoardCoverSources[source]; !exist {
		errors = append(errors, custom_errors.ErrInvalidCoverSource)
	}

	if focalPointY > 1 {
		errors = append(errors, custom_errors.ErrUnsplashFocalPointYTooHigh)
	}

	if focalPointY < 0 {
		errors = append(errors, custom_errors.ErrUnsplashFocalPointYTooLow)
	}

	if title == "" {
		errors = append(errors, custom_errors.ErrTitleEmpty)
	}

	if visibility != models.BoardVisibilityPrivate && visibility != models.BoardVisibilityPublic {
		errors = append(errors, custom_errors.ErrInvalidVisibility)
	}

	if len(errors) > 0 {
		return &custom_errors.MultipleErrors{Errors: errors}
	}

	imageFiles, err := usecase.unsplashRepo.GetImagesForID(photoID, focalPointY)
	if err != nil {
		return err
	}

	boardCover := &models.BoardCover{
		PhotoID:     photoID,
		Source:      source,
		FocalPointY: focalPointY,
		Images:      make(models.Images, len(board.BoardCoverSizes)),
	}

	for i, width := range board.BoardCoverSizes {
		image := &models.Image{
			Width: width,
		}

		boardCover.Images[i] = image
	}

	var wg sync.WaitGroup
	uploadChannels := make(chan error, len(boardCover.Images))
	wg.Add(len(boardCover.Images))

	for idx, image := range boardCover.Images {
		imageFile := imageFiles[idx]
		name := utils.RandString(8)
		fileName := fmt.Sprintf("%s.%s", name, utils.GetFileExtension(imageFile.Name()))

		metaData := map[string]string{
			"name":        fileName,
			"title":       name,
			"description": fmt.Sprintf("cover picture of board %s with width of %v", title, image.Width),
		}

		go usecase.storage.UploadFile(uploadChannels, &wg, image, imageFile, metaData)
	}

	wg.Wait()
	close(uploadChannels)

	for err = range uploadChannels {
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return &custom_errors.MultipleErrors{Errors: errors}
	}

	<-uploadChannels

	_board := &models.Board{
		Title:   title,
		OwnerID: userID,
		Cover:   boardCover,
	}
	_board.SetVisibility(visibility)
	_board.EmptyImageURLs()

	err = usecase.boardRepo.Create(_board)
	if err != nil {
		return err
	}

	return nil
}
