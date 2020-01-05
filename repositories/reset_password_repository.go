package repositories

import (
	"context"
	"github.com/vuongabc92/octocv/models"
	"go.mongodb.org/mongo-driver/bson"
)

type ResetPasswordRepository struct {
	repository
}

func (r *ResetPasswordRepository) FindByEmail(ctx context.Context, email string) (resetPassword models.ResetPassword, err error) {
	err = r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&resetPassword)
	return
}

func (r *ResetPasswordRepository) FindByToken(ctx context.Context, token string) (resetPassword models.ResetPassword, err error) {
	err = r.collection.FindOne(ctx, bson.M{"token": token}).Decode(&resetPassword)
	return
}

func (r *ResetPasswordRepository) CollectionName() string {
	return collectionPrefix + "reset_password"
}
