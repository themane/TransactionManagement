package models

type CustomerData struct {
	Id      string          `json:"_id" bson:"_id"`
	Name    string          `json:"name" bson:"name"`
	Father  string          `json:"father_name" bson:"father_name"`
	Mobiles []string        `json:"mobiles" bson:"mobiles"`
	Address CustomerAddress `json:"address" bson:"address"`
	New     bool            `json:"new,omitempty" bson:"-"`
}

type CustomerAddress struct {
	Village string   `json:"village" bson:"village"`
	Tags    []string `json:"tags" bson:"tags"`
}

type CustomerRepository interface {
	FindByIds(ids []string) ([]CustomerData, error)
	FindByName(name string, page int, pageSize int) ([]CustomerData, error)
	FindByMobile(mobile string) (*CustomerData, error)
	FindByVillage(village string, page int, pageSize int) ([]CustomerData, error)
	FindByKeyword(keyword string, page int, pageSize int) ([]CustomerData, error)

	AddUser(customerData CustomerData) error
}
