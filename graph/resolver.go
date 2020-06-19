package graph

import (
	"crypto/rand"
	"fmt"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{}

func GetRandomNumber() string {
	var array = make([]byte, 4)
	if _, err := rand.Read(array); err != nil {
		return ""
	}
	Random := fmt.Sprintf("%X", array)
	return Random
}

// func (r *mutationResolver) Register(ctx context.Context, input *model.NewUser) (*model.User, error) {
// 	var user model.User
// 	var userBD = db.GetCollection("users")
// 	var PrivateKey = GetRandomNumber() + input.UserID
// 	var Token = auth.GenerateJWT(PrivateKey)
// 	if err := userBD.Find(bson.M{"identity": input.Identity}).One(&user); err == nil {
// 		return &model.User{}, errors.New("Usuario ya registrado")
// 	}
// 	userBD.Insert(bson.M{
// 		"name":        input.Name,
// 		"identity":    input.Identity,
// 		"private_key": PrivateKey,
// 		"token":       Token,
// 	})
// 	if err := userBD.Find(bson.M{"private_key": PrivateKey}).Select(bson.M{"_id": 0}).One(&user); err != nil {
// 		return &model.User{}, err
// 	}
// 	return &user, nil
// }
// func (r *mutationResolver) AddPurchase(ctx context.Context, input *model.NewPurchase) (*model.Purchase, error) {
// 	userContext := auth.GetAuthFromContext(ctx)
// 	if userContext == nil {
// 		return &model.Purchase{}, errors.New("Acceso denegado")
// 	}
// 	var userBD = db.GetCollection("users")
// 	var user model.User
// 	if err := userBD.Find(bson.M{"private_key": userContext.PrivateKey}).One(&user); err != nil {
// 		return &model.Purchase{}, errors.New("No se encontró usuario relacionado")
// 	}
// 	var purchase model.Purchase
// 	var purchaseBD = db.GetCollection("purchases")
// 	if err := purchaseBD.Find(bson.M{"contact_identity": input.ContactIdentity, "serie": input.Serie, "number": input.Number}).One(&purchase); err == nil {
// 		return &model.Purchase{}, errors.New("La compra ya está registrada")
// 	}

// 	id := bson.NewObjectId()
// 	purchaseBD.Insert(bson.M{
// 		"_id":              id,
// 		"serie":            input.Serie,
// 		"number":           input.Number,
// 		"contact_identity": input.ContactIdentity,
// 		"user_id":          user.ID,
// 		"contact_name":     input.ContactName,
// 		"total_price":      input.TotalPrice,
// 		"total_igv":        input.TotalIgv,
// 		"issue_date":       input.IssueDate,
// 	})
// 	if err := purchaseBD.Find(bson.M{"_id": id}).Select(bson.M{"_id": 0}).One(&purchase); err != nil {
// 		return &model.Purchase{}, errors.New("No se encontró la compra")
// 	}
// 	return &purchase, nil
// }
