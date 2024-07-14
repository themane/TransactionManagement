package utils

import (
	constants "TxnManagement/contants"
	"TxnManagement/controllers/exceptions"
	"TxnManagement/repositories/models"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
)

func ValidateIdToken(c *gin.Context) (*models.AdminData, error) {
	idToken := extractToken(c)
	return ParseGoogleIdToken(idToken)
}

func ParseGoogleIdToken(idToken string) (*models.AdminData, error) {
	payload, err := idtoken.Validate(context.Background(), idToken, "")
	if err != nil {
		return nil, &exceptions.NoSuchCombinationError{
			Message: err.Error(),
		}
	}
	userDetails := models.AdminData{
		Id:            fmt.Sprintf("%v", payload.Claims["sub"]),
		Name:          fmt.Sprintf("%v", payload.Claims["name"]),
		Email:         fmt.Sprintf("%v", payload.Claims["email"]),
		Authenticator: constants.GoogleAuthenticator,
	}
	return &userDetails, nil
}
