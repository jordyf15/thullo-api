package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jordyf15/thullo-api/custom_errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func respondBasedOnError(c *gin.Context, err error) {
	statusCode := getStatusCodeForError(err)

	switch actualErr := err.(type) {
	case *custom_errors.MultipleErrors:
		c.JSON(statusCode, actualErr)
		return
	case *custom_errors.Error:
		c.JSON(statusCode, custom_errors.MultipleErrors{Errors: []error{actualErr}})
		return
	default:
		break
	}

	switch err {
	case mongo.ErrNoDocuments:
		c.JSON(http.StatusNotFound, custom_errors.MultipleErrors{Errors: []error{custom_errors.ErrRecordNotFound}})
	case nil:
		c.Status(statusCode)
	default:
		c.JSON(statusCode, custom_errors.MultipleErrors{Errors: []error{custom_errors.ErrUnknownErrorOccured}})
	}
}

func getStatusCodeForError(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if _, ok := err.(*custom_errors.MultipleErrors); ok {
		return http.StatusBadRequest
	} else if modelError, ok := err.(*custom_errors.Error); ok {
		switch modelError {
		case custom_errors.ErrMalformedRefreshToken, custom_errors.ErrInvalidRefreshToken:
			return http.StatusForbidden
		default:
			return http.StatusBadRequest
		}
	}

	switch err {
	case mongo.ErrNoDocuments:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
