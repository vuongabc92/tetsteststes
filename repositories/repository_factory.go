package repositories

import (
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionPrefix string = "octocv_"

type RepositoryFactory struct {
	db *mongo.Database
}

func NewRepositoryFactory(db *mongo.Database) *RepositoryFactory {
	return &RepositoryFactory{
		db: db,
	}
}

func (r *RepositoryFactory) User() *UserRepository {
	userRepo := new(UserRepository)
	collection := r.db.Collection(userRepo.CollectionName())
	userRepo.SetCollection(collection)
	return userRepo
}

func (r *RepositoryFactory) ResetPassword() *ResetPasswordRepository {
	resetPassRepo := new(ResetPasswordRepository)
	collection := r.db.Collection(resetPassRepo.CollectionName())
	resetPassRepo.SetCollection(collection)
	return resetPassRepo
}
