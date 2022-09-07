package model

type User struct {
	BaseModel
	Email    string `json:"email"`
	Password string `json:"password"`
	Profile  string `json:"profile"`
}
