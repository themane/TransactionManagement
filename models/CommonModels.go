package models


type Product struct {
	Type       string   `json:"type" bson:"type"`
	Weight     float32  `json:"weight" bson:"weight"`
	Tags       []string `json:"tags" bson:"tags"`
	UnitPrice  int      `json:"unit_price" bson:"unit_price"`
	TotalPrice int      `json:"total_price" bson:"total_price"`
}