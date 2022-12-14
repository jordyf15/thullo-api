package usecase

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jordyf15/thullo-api/custom_errors"
	"github.com/jordyf15/thullo-api/models"
	"github.com/jordyf15/thullo-api/oauth"
	"github.com/jordyf15/thullo-api/storage"
	"github.com/jordyf15/thullo-api/token"
	"github.com/jordyf15/thullo-api/user"
	"github.com/jordyf15/thullo-api/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	httpHeaderFilenameRegex = regexp.MustCompile("filename=\"([A-Za-z0-9.]+)\"")
	gidPictureSizeRegex     = regexp.MustCompile(`s\d+-c`)
)

type userUsecase struct {
	userRepo  user.Repository
	tokenRepo token.Repository
	oauthRepo oauth.Repository
	storage   storage.Storage
}

type userInstanceUsecase struct {
	user *models.User
	userUsecase
}

func NewUserUsecase(userRepo user.Repository, tokenRepo token.Repository, oauthRepo oauth.Repository, storage storage.Storage) user.Usecase {
	return &userUsecase{userRepo: userRepo, tokenRepo: tokenRepo, oauthRepo: oauthRepo, storage: storage}
}

func (usecase *userUsecase) Create(_user *models.User, imageFile utils.NamedFileReader) (map[string]interface{}, error) {
	var err error
	errors := make([]error, 0)

	validateFieldErrors := _user.VerifyFields()
	if len(validateFieldErrors) > 0 {
		errors = append(errors, validateFieldErrors...)
	}

	isUsernameExist, err := usecase.userRepo.FieldExists("username", _user.Username)
	if err != nil {
		return nil, err
	}
	if isUsernameExist {
		errors = append(errors, custom_errors.ErrUsernameAlreadyExists)
	}

	isEmailExist, err := usecase.userRepo.FieldExists("email", _user.Email)
	if err != nil {
		return nil, err
	}
	if isEmailExist {
		errors = append(errors, custom_errors.ErrEmailAddressAlreadyRegistered)
	}

	err = _user.SetPassword(_user.Password)
	if err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return nil, &custom_errors.MultipleErrors{Errors: errors}
	}

	var userAvatar utils.NamedFileReader
	if imageFile != nil {
		userAvatar = imageFile
	} else {
		avatar, err := utils.GenerateAvatar(_user.Initials(), 800, 400)
		if err != nil {
			return nil, err
		}

		defer os.Remove(avatar.Name())

		userAvatar = utils.NewNamedFileReader(avatar, avatar.Name())
	}

	fileExtension := utils.GetFileExtension(userAvatar.Name())
	switch fileExtension {
	case "jpg", "jpeg", "png":
	default:
		return nil, custom_errors.ErrImageFormatInvalid
	}

	_user.Images = make([]*models.Image, len(user.DisplayPictureSizes))
	for i, width := range user.DisplayPictureSizes {
		image := &models.Image{}
		image.Width = width

		_user.Images[i] = image
	}

	uploadChannels := make(chan error, len(_user.Images))
	var wg sync.WaitGroup

	wg.Add(len(_user.Images))

	for _, img := range _user.Images {
		var resizedImageFile *os.File
		resizedImageFile, err = utils.ResizeImage(userAvatar, int(img.Width))
		if err != nil {
			break
		}

		defer os.Remove(resizedImageFile.Name())

		name := utils.RandString(8)
		fileName := fmt.Sprintf("%s.%s", name, utils.GetFileExtension(resizedImageFile.Name()))

		metaData := map[string]string{
			"name":        fileName,
			"title":       name,
			"description": fmt.Sprintf("profile picture of %s with width and height of %v", _user.Username, img.Width),
		}

		go usecase.storage.UploadFile(uploadChannels, &wg, img, resizedImageFile, metaData)
	}

	wg.Wait()
	close(uploadChannels)

	for err = range uploadChannels {
		if err != nil {
			fmt.Println(err)
		}
	}

	<-uploadChannels

	if err := usecase.userRepo.Create(_user); err != nil {
		return nil, err
	}

	accessToken, refreshToken, _ := usecase.For(_user).GenerateTokens()

	_user.EmptyImageIDs()

	response := utils.DataResponse(_user, map[string]interface{}{
		"access_token":  accessToken.ToJWTString(),
		"refresh_token": refreshToken.ToJWTString(),
		"expires_at":    accessToken.ExpiresAt,
	})

	return response, nil
}

func (usecase *userUsecase) LoginWithGoogle(token string) (map[string]interface{}, error) {
	tokenInfo, err := usecase.oauthRepo.GetGoogleTokenInfo(token)
	if err != nil {
		return nil, err
	}

	if tokenInfo.ExpiresAt < time.Now().Unix() {
		return nil, custom_errors.ErrGoogleOauthTokenExpired
	}

	user, err := usecase.userRepo.GetByEmail(tokenInfo.Email)
	if err == nil {
		accessToken, refreshToken, err := usecase.For(user).GenerateTokens()
		if err != nil {
			return nil, err
		}

		err = usecase.storage.AssignImageURLToUser(user)
		if err != nil {
			return nil, err
		}

		response := utils.DataResponse(user, map[string]interface{}{
			"access_token":  accessToken.ToJWTString(),
			"refresh_token": refreshToken.ToJWTString(),
			"expires_at":    accessToken.ExpiresAt,
		})

		return response, nil
	}

	if err != mongo.ErrNoDocuments {
		return nil, err
	}

	regex := regexp.MustCompile("[^A-Za-z0-9]")
	password := "0.Aa" + utils.RandString(8)

	var username string
	if len(tokenInfo.Name) >= models.MinUsernameLength {
		username = strings.ToLower(regex.ReplaceAllString(tokenInfo.Name, ""))
	} else {
		username = fmt.Sprintf("User %d", rand.Intn(10000))
		username = strings.ToLower(regex.ReplaceAllString(username, ""))
	}

	if exists, err := usecase.userRepo.FieldExists("username", username); err != nil {
		return nil, err
	} else if exists {
		seed := 100
		for exists {
			newUsername := username + strconv.Itoa(rand.Intn(seed))
			exists, err = usecase.userRepo.FieldExists("username", newUsername)
			if err != nil {
				return nil, err
			} else if !exists {
				username = newUsername
			} else {
				seed *= 10
			}
		}
	}

	user = &models.User{Name: tokenInfo.Name, Username: username, Email: tokenInfo.Email, Password: password}

	if len(tokenInfo.Picture) > 0 {
		tmpFile, err := ioutil.TempFile(os.TempDir(), "googleimg-")
		if err != nil {
			return usecase.Create(user, nil)
		}

		defer os.Remove(tmpFile.Name())

		pictureURL := gidPictureSizeRegex.ReplaceAllString(tokenInfo.Picture, "s800-c")
		respHeader, err := utils.DownloadFile(tmpFile.Name(), pictureURL)
		if err != nil {
			return usecase.Create(user, nil)
		}

		contentDisposition := strings.Join(respHeader.Values("Content-Disposition"), "")
		matches := httpHeaderFilenameRegex.FindStringSubmatch(contentDisposition)
		var filename string
		if len(matches) > 1 {
			filename = matches[1]
		} else {
			contentType := strings.Join(respHeader.Values("Content-type"), "")
			if len(contentType) > 0 && strings.Contains(contentType, "/") {
				split := strings.Split(contentType, "/")
				filename = "a." + split[len(split)-1]
			} else {
				filename = filepath.Base(tokenInfo.Picture)
			}
		}

		response, err := usecase.Create(user, utils.NewNamedFileReader(tmpFile, filename))
		return response, err
	}

	return usecase.Create(user, nil)
}

func (usecase *userUsecase) For(user *models.User) user.InstanceUsecase {
	instanceUsecase := &userInstanceUsecase{user: user, userUsecase: *usecase}
	return instanceUsecase
}

func (usecase *userUsecase) Login(email, password string) (map[string]interface{}, error) {
	user, err := usecase.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, custom_errors.ErrCurrentPasswordWrong
		}
		return nil, err
	}

	accessToken, refreshToken, err := usecase.For(user).GenerateTokens()
	if err != nil {
		return nil, err
	}

	err = usecase.storage.AssignImageURLToUser(user)
	if err != nil {
		return nil, err
	}

	user.EmptyImageIDs()

	response := utils.DataResponse(user, map[string]interface{}{
		"access_token":  accessToken.ToJWTString(),
		"refresh_token": refreshToken.ToJWTString(),
		"expires_at":    accessToken.ExpiresAt,
	})

	return response, nil
}

// userInstanceUsecase
func (usecase *userInstanceUsecase) GenerateTokens() (*models.AccessToken, *models.RefreshToken, error) {
	refreshToken := (&models.RefreshToken{UserID: usecase.user.ID})
	refreshToken.Id = utils.RandString(8)

	accessToken := (&models.AccessToken{UserID: usecase.user.ID}).SetExpiration(time.Now().Add(time.Hour * 1))
	accessToken.RefreshTokenID = utils.ToSHA256(refreshToken.Id)

	tokenSet := &models.TokenSet{UserID: usecase.user.ID, RefreshTokenID: accessToken.RefreshTokenID}
	err := usecase.tokenRepo.Create(tokenSet)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}
