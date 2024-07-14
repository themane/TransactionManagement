package controllers

import (
	constants "TxnManagement/contants"
	controllerModels "TxnManagement/controllers/models"
	"TxnManagement/controllers/utils"
	repoModels "TxnManagement/repositories/models"
	"TxnManagement/services"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *services.AuthService
	apiSecret   string
	logger      *constants.LoggingUtils
}

func NewAuthController(adminRepository repoModels.AdminRepository,
	apiSecret string,
	logLevel string,
) *AuthController {
	return &AuthController{
		authService: services.NewAuthService(adminRepository, logLevel),
		apiSecret:   apiSecret,
		logger:      constants.NewLoggingUtils("AUTH_CONTROLLER", logLevel),
	}
}

// Register godoc
// @Summary Register API for admin controls
// @Description Registration payload verification and initial assignment of complete admin data
// @Tags registration
// @Accept json
// @Produce json
// @Success 200 {object} models.LoginResponse
// @Router /register [post]
func (a *AuthController) Register(c *gin.Context) {
	adminData, err := utils.ValidateIdToken(c)
	if err != nil {
		a.logger.Error("Error in user authentication", err)
		c.JSON(401, err.Error())
		return
	}
	a.logger.Printf("Registering admin: %s", adminData.Email)

	err = a.authService.AddUser(*adminData)
	if err != nil {
		a.logger.Error("error in adding admin to Mongo", err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting admin registered. contact administrators for more info", HttpCode: 500})
		return
	}
	a.login(c, adminData)
	c.Status(200)
}

// Login godoc
// @Summary Login API
// @Description Login verification and first load of complete user data
// @Tags data retrieval
// @Accept json
// @Produce json
// @Success 200 {object} models.LoginResponse
// @Router /login [post]
func (a *AuthController) Login(c *gin.Context) {
	userDetails, err := utils.ValidateIdToken(c)
	if err != nil {
		a.logger.Error("Error in admin authentication", err)
		c.JSON(401, err.Error())
		return
	}
	a.login(c, userDetails)
	c.Status(200)
}

// RefreshToken godoc
// @Summary Refresh Token API
// @Description Refresh Token
// @Tags data retrieval
// @Accept json
// @Produce json
// @Router /token/refresh [post]
func (a *AuthController) RefreshToken(c *gin.Context) {
	email, err := utils.RefreshTokenValid(c, a.apiSecret)
	if err != nil {
		a.logger.Error("Error in admin authentication", err)
		c.JSON(401, err.Error())
		return
	}
	a.addTokens(c, email)
	c.Status(200)
}
