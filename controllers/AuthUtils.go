package controllers

import (
	constants "TxnManagement/contants"
	"TxnManagement/controllers/exceptions"
	controllerModels "TxnManagement/controllers/models"
	"TxnManagement/controllers/utils"
	"TxnManagement/repositories/models"
	"errors"
	"github.com/gin-gonic/gin"
)

func (a *AuthController) login(c *gin.Context, userDetails *models.AdminData) {
	a.logger.Printf("Logging in user: %s", userDetails.Email)
	var response *models.AdminData
	var err error
	if userDetails.Authenticator == constants.GoogleAuthenticator {
		response, err = a.authService.GoogleLogin(userDetails.Id)
	}
	var noSuchCombinationError *exceptions.NoSuchCombinationError
	if errors.As(err, &noSuchCombinationError) {
		a.logger.Error("admin not registered", err)
		c.JSON(204, err.Error())
		return
	}
	if err != nil {
		a.logger.Error("error in getting admin data", err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting admin data. contact administrators for more info", HttpCode: 500})
		return
	}
	if response == nil {
		msg := "Admin data not found"
		a.logger.Info(msg)
		c.JSON(204, msg)
		return
	}
	a.addTokens(c, response.Email)
}

func (a *AuthController) addTokens(c *gin.Context, email string) {
	token, err := utils.GenerateToken(email, a.apiSecret)
	if err != nil {
		a.logger.Error("error in getting auth token generation", err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	c.Header("X-Api-Token", token)

	refreshToken, err := utils.GenerateRefreshToken(email, a.apiSecret)
	if err != nil {
		a.logger.Error("error in getting auth token generation", err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting user data. contact administrators for more info", HttpCode: 500})
		return
	}
	c.Header("X-Refresh-Token", refreshToken)
}
