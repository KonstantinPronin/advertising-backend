package model

type Ad struct {
	Id          string   `json:"id,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Price       float64  `json:"price"`
	Images      []string `json:"images"`
	Created     string   `json:"created,omitempty"`
}

//easyjson:json
type AdList []Ad
