package model

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	PrivateKey string `json:"private_key" bson:"private_key"`
	Token      string `json:"token"`
	Identity   string `json:"identity"`
}
