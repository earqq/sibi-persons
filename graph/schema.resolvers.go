package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/earqq/gqlgen-easybill/auth"
	"github.com/earqq/gqlgen-easybill/db"
	"github.com/earqq/gqlgen-easybill/graph/generated"
	"github.com/earqq/gqlgen-easybill/graph/model"
	"gopkg.in/mgo.v2/bson"
)

func (r *mutationResolver) Login(ctx context.Context, privateKey *string) (*model.User, error) {
	var user model.User
	var userBD = db.GetCollection("users")
	if err := userBD.Find(bson.M{"private_key": privateKey}).Select(bson.M{"_id": 0, "password": 0}).One(&user); err != nil {
		return &model.User{}, errors.New("No existe usuario con esa clave")
	}
	return &user, nil
}

func (r *mutationResolver) Register(ctx context.Context, input *model.NewUser) (*model.User, error) {
	var user model.User
	var userBD = db.GetCollection("users")
	var PrivateKey = GetRandomNumber() + input.UserID
	var Token = auth.GenerateJWT(PrivateKey)
	if err := userBD.Find(bson.M{"identity": input.Identity}).One(&user); err == nil {
		return &model.User{}, errors.New("Usuario ya registrado")
	}
	userBD.Insert(bson.M{
		"name":        input.Name,
		"identity":    input.Identity,
		"private_key": PrivateKey,
		"token":       Token,
	})
	if err := userBD.Find(bson.M{"private_key": PrivateKey}).Select(bson.M{"_id": 0}).One(&user); err != nil {
		return &model.User{}, err
	}
	return &user, nil
}

func (r *mutationResolver) AddPurchase(ctx context.Context, input *model.NewPurchase) (*model.Purchase, error) {
	userContext := auth.GetAuthFromContext(ctx)
	if userContext == nil {
		return &model.Purchase{}, errors.New("Acceso denegado")
	}
	var userBD = db.GetCollection("users")
	var user model.User
	if err := userBD.Find(bson.M{"private_key": userContext.PrivateKey}).One(&user); err != nil {
		return &model.Purchase{}, errors.New("No se encontró usuario relacionado")
	}
	var purchase model.Purchase
	var purchaseBD = db.GetCollection("purchases")
	if err := purchaseBD.Find(bson.M{"contact_identity": input.ContactIdentity, "serie": input.Serie, "number": input.Number}).One(&purchase); err == nil {
		return &model.Purchase{}, errors.New("La compra ya está registrada")
	}

	id := bson.NewObjectId()
	purchaseBD.Insert(bson.M{
		"_id":              id,
		"serie":            input.Serie,
		"number":           input.Number,
		"contact_identity": input.ContactIdentity,
		"user_id":          user.ID,
		"contact_name":     input.ContactName,
		"total_price":      input.TotalPrice,
		"total_igv":        input.TotalIgv,
		"issue_date":       input.IssueDate,
	})
	if err := purchaseBD.Find(bson.M{"_id": id}).Select(bson.M{"_id": 0}).One(&purchase); err != nil {
		return &model.Purchase{}, errors.New("No se encontró la compra")
	}
	return &purchase, nil
}

func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	userContext := auth.ForContext(ctx)
	if userContext == nil {
		return &model.User{}, errors.New("Acceso denegado")
	}
	var userBD = db.GetCollection("users")
	var user model.User
	if err := userBD.Find(bson.M{"private_key": userContext.PrivateKey}).Select(bson.M{"_id": 0}).One(&user); err != nil {
		return &model.User{}, err
	}
	return &user, nil
}

func (r *queryResolver) Purchases(ctx context.Context) ([]*model.Purchase, error) {
	userContext := auth.GetAuthFromContext(ctx)
	if userContext == nil {
		return nil, errors.New("Acceso denegado")
	}
	var purchases []*model.Purchase
	var purchaseBD = db.GetCollection("purchases")
	var userBD = db.GetCollection("users")
	var user model.User
	if err := userBD.Find(bson.M{"private_key": userContext.PrivateKey}).One(&user); err != nil {
		return nil, errors.New("No se encontró usuario relacionado")
	}
	if err := purchaseBD.Find(bson.M{"user_id": user.ID}).Select(bson.M{"_id": 0}).All(&purchases); err != nil {
		return nil, errors.New("No se encontró la compra")
	}
	return purchases, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
