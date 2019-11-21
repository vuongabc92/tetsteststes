package repositories

import (
	"context"
	"github.com/vuongabc92/octocv/models"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository struct {
	repository
}

func (u *UserRepository) FindByEmail(ctx context.Context, email string) (user models.User, err error) {
	err = u.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return
}

func (u *UserRepository) CollectionName() string {
	return collectionPrefix + "users"
}
