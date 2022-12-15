package storage

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/jordyf15/thullo-api/custom_errors"
	"github.com/jordyf15/thullo-api/models"
)

const (
	baseURL = "https://api.imgur.com"
)

type imgurStorage struct {
	client             *http.Client
	accessToken        string
	accessTokenExpired int64
}

func NewImgurStorage(client *http.Client) Storage {
	return &imgurStorage{client: client}
}

func (api *imgurStorage) UploadFile(respond chan<- error, wg *sync.WaitGroup, currentImage *models.Image, file *os.File, metadata map[string]string) {
	if wg != nil {
		defer wg.Done()
	}

	var albumID string
	if os.Getenv("ENV") == "production" {
		albumID = os.Getenv("IMGUR_ALBUM_ID_PRODUCTION")
	} else {
		albumID = os.Getenv("IMGUR_ALBUM_ID_DEVELOPMENT")
	}

	bits, _ := os.ReadFile(file.Name())
	encoded := base64.StdEncoding.EncodeToString(bits)

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	_ = writer.WriteField("image", encoded)

	album, _ := writer.CreateFormField("album")
	album.Write([]byte(albumID))

	typeField, _ := writer.CreateFormField("type")
	typeField.Write([]byte("base64"))

	name, _ := writer.CreateFormField("name")
	name.Write([]byte(metadata["name"]))

	title, _ := writer.CreateFormField("title")
	title.Write([]byte(metadata["title"]))

	description, _ := writer.CreateFormField("description")
	description.Write([]byte(metadata["description"]))

	writer.Close()

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/3/image", baseURL), buf)
	if err != nil {
		respond <- err
		return
	}
	if api.accessToken == "" || api.accessTokenExpired < time.Now().Unix() {
		api.accessToken, api.accessTokenExpired, err = api.getAccessToken()
		if err != nil {
			respond <- err
			return
		}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.accessToken))
	req.Header.Add("Content-Type", writer.FormDataContentType())

	res, err := api.client.Do(req)
	if err != nil {
		respond <- err
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		respond <- custom_errors.ErrUnknownErrorOccured
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		respond <- err
		return
	}

	uploadFileRes := &uploadFileResponse{}

	err = json.Unmarshal(body, uploadFileRes)
	if err != nil {
		respond <- err
		return
	}

	currentImage.ID = uploadFileRes.Data.ID
	currentImage.URL = uploadFileRes.Data.Link

	respond <- nil
}

func (api *imgurStorage) getAccessToken() (string, int64, error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	refreshToken, _ := writer.CreateFormField("refresh_token")
	refreshToken.Write([]byte(os.Getenv("IMGUR_REFRESH_TOKEN")))

	clientID, _ := writer.CreateFormField("client_id")
	clientID.Write([]byte(os.Getenv("IMGUR_CLIENT_ID")))

	clientSecret, _ := writer.CreateFormField("client_secret")
	clientSecret.Write([]byte(os.Getenv("IMGUR_CLIENT_SECRET")))

	grantType, _ := writer.CreateFormField("grant_type")
	grantType.Write([]byte("refresh_token"))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/oauth2/token", baseURL), buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		return "", 0, err
	}

	res, err := api.client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", 0, err
	}

	authorizationRes := &imgurAuthorizationResponse{}

	err = json.Unmarshal(body, &authorizationRes)
	if err != nil {
		return "", 0, err
	}

	return authorizationRes.AccessToken, authorizationRes.AccessTokenExpiresIn, nil
}

type uploadFileResponse struct {
	Status  int            `json:"status"`
	Success bool           `json:"success"`
	Data    uploadFileData `json:"data"`
}

type uploadFileData struct {
	ID          string `json:"id"`
	DeleteHash  string `json:"deletehash"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Size        int64  `json:"size"`
	Link        string `json:"link"`
}

type imgurAuthorizationResponse struct {
	AccessToken          string `json:"access_token"`
	AccessTokenExpiresIn int64  `json:"expires_in"`
	RefreshToken         string `json:"refresh_token"`
}
