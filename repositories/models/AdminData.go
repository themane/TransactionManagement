package models

type AdminData struct {
	Id            string `json:"-" bson:"_id"`
	Name          string `json:"name" bson:"name"`
	Email         string `json:"email" bson:"email"`
	Authenticator string `json:"-" bson:"authenticator"`
}

type AdminRepository interface {
	FindById(userId string) (*AdminData, error)
	FindByEmail(id string) (*AdminData, error)

	AddUser(adminData AdminData) error
}
