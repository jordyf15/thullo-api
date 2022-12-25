package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/jordyf15/thullo-api/custom_errors"
	"github.com/jordyf15/thullo-api/unsplash"
	"github.com/jordyf15/thullo-api/utils"
)

const (
	baseURL = "https://api.unsplash.com"
)

type unsplashRepository struct {
	client *http.Client
}

func NewUnsplashRepository(client *http.Client) unsplash.Repository {
	return &unsplashRepository{client: client}
}

func (repo *unsplashRepository) GetImagesForID(photoID string, focalPointY float64) ([]*os.File, error) {
	photoURL := fmt.Sprintf("%s/photos/%s", baseURL, photoID)

	req, err := http.NewRequest("GET", photoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Client-ID "+os.Getenv("UNSPLASH_ACCESS_KEY"))

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, custom_errors.ErrUnknownUnsplashPhotoID
		default:
			return nil, custom_errors.ErrUnknownErrorOccured
		}
	}

	receivedResp := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&receivedResp)
	if err != nil {
		return nil, err
	}

	urls := receivedResp["urls"].(map[string]interface{})
	rawURL, err := url.Parse(urls["raw"].(string))
	if err != nil {
		return nil, err
	}

	imageFiles := []*os.File{}

	q, _ := url.ParseQuery(rawURL.RawQuery)
	q.Del("crop")
	rawURL.RawQuery = q.Encode()

	imageRegular, err := utils.GetImageFromURL(fmt.Sprintf("%s&w=1080&h=360&fit=crop&auto=format&fm=jpg&fp-y=%f", rawURL.String(), focalPointY))
	if err != nil {
		return nil, err
	}

	imageSmall, err := utils.GetImageFromURL(fmt.Sprintf("%s&w=450&h=150&fit=crop&auto=format&fm=jpg&fp-y=%f", rawURL.String(), focalPointY))
	if err != nil {
		return nil, err
	}

	imageThumb, err := utils.GetImageFromURL(fmt.Sprintf("%s&w=150&h=50&fit=crop&auto=format&fm=jpg&fp-y=%f", rawURL.String(), focalPointY))
	if err != nil {
		return nil, err
	}

	imageFiles = append(imageFiles, imageRegular)
	imageFiles = append(imageFiles, imageSmall)
	imageFiles = append(imageFiles, imageThumb)

	return imageFiles, nil
}
