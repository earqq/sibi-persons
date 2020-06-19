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

func (r *mutationResolver) Login(ctx context.Context, privateKey string) (*model.User, error) {
	var user model.User
	var userBD = db.GetCollection("users")
	if err := userBD.Find(bson.M{"private_key": privateKey}).Select(bson.M{"_id": 0, "password": 0}).One(&user); err != nil {
		return &model.User{}, errors.New("No existe usuario con esa clave")
	}
	if user.Token == "" {
		var Token = auth.GenerateJWT(privateKey)
		if err := userBD.Update(bson.M{"private_key": privateKey}, bson.M{"$set": bson.M{"token": Token}}); err != nil {
			return &model.User{}, errors.New("No se pudo actualizar token")
		}
		_ = userBD.Find(bson.M{"private_key": privateKey}).Select(bson.M{"_id": 0, "password": 0}).One(&user)
	}
	return &user, nil
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

func (r *queryResolver) Purchases(ctx context.Context, search *string, limit *int) ([]*model.Purchase, error) {
	userContext := auth.GetAuthFromContext(ctx)
	if userContext == nil {
		return nil, errors.New("Acceso denegado")
	}
	var purchases []*model.Purchase
	var fields = bson.M{}
	var purchaseBD = db.GetCollection("purchases")
	var userBD = db.GetCollection("users")
	var user model.User
	if err := userBD.Find(bson.M{"private_key": userContext.PrivateKey}).One(&user); err != nil {
		return nil, errors.New("No se encontró usuario relacionado")
	}
	if search != nil {
		fields["search"] = bson.M{"$regex": *search, "$options": "i"}
	}
	fields["user_id"] = bson.ObjectId(user.ID).Hex()
	if limit != nil {
		purchaseBD.Find(fields).Limit(*limit).Sort("-issue_date").Select(bson.M{"_id": 0}).All(&purchases)
	} else {
		purchaseBD.Find(fields).Sort("-issue_date").Select(bson.M{"_id": 0}).All(&purchases)
	}
	return purchases, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
