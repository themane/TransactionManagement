package services

import (
	constants "TxnManagement/contants"
	"TxnManagement/controllers/exceptions"
	"TxnManagement/repositories/models"
	"errors"
)

type AuthService struct {
	adminRepository models.AdminRepository
	logger          *constants.LoggingUtils
}

func NewAuthService(adminRepository models.AdminRepository,
	logLevel string,
) *AuthService {
	return &AuthService{
		adminRepository: adminRepository,
		logger:          constants.NewLoggingUtils("AUTH_SERVICE", logLevel),
	}
}

func (a *AuthService) AddUser(adminData models.AdminData) error {
	if a.userExists(adminData.Email) {
		return errors.New("username already taken")
	}
	return a.adminRepository.AddUser(adminData)
}

func (a *AuthService) GoogleLogin(id string) (*models.AdminData, error) {
	userData, err := a.adminRepository.FindById(id)
	if err != nil {
		return nil, &exceptions.NoSuchCombinationError{
			Message: err.Error(),
		}
	}
	return userData, nil
}

func (a *AuthService) userExists(email string) bool {
	user, err := a.adminRepository.FindByEmail(email)
	return err == nil && user != nil && user.Email == email
}
