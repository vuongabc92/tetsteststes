package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/vuongabc92/octocv/helpers"
	"github.com/vuongabc92/octocv/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

type repository struct {
	collection *mongo.Collection
}

func (r *repository) SetCollection(c *mongo.Collection) {
	r.collection = c
}

func (r *repository) Insert(ctx context.Context, model interface{}) error {
	// If there is field name CreatedAt or UpdatedAt
	// then auto set field value to current datetime
	var (
		reflectValue = reflect.Indirect(reflect.ValueOf(model))
		createdAt    = reflectValue.FieldByName("CreatedAt")
		updatedAt    = reflectValue.FieldByName("UpdatedAt")
		now          = helpers.Now()
	)

	if createdAt.IsValid() {
		createdAt.Set(reflect.ValueOf(now))
	}

	if updatedAt.IsValid() {
		updatedAt.Set(reflect.ValueOf(now))
	}

	// Insert data and if no errors then set model ID = inserted result ID
	insertResult, err := r.collection.InsertOne(ctx, model)
	if err == nil {
		ID := reflectValue.FieldByName("ID")
		if ID.IsValid() && ID.Type().String() == "primitive.ObjectID" {
			ID.Set(reflect.ValueOf(insertResult.InsertedID))
		}
	}

	return err
}

// Update collection by collection ID
func (r *repository) Update(ctx context.Context, model interface{}) error {
	var (
		reflectValue = reflect.Indirect(reflect.ValueOf(model))
		modelId      = reflectValue.FieldByName("Id")
		updatedAt    = reflectValue.FieldByName("UpdatedAt")
	)

	// Every model must contains a field name `Id` and it is primary key.
	if !modelId.IsValid() {
		return errors.New("model id is invalid")
	}

	// If model contains a field name UpdatedAt then every time model is updated
	// we have to set it to current time.
	if updatedAt.IsValid() {
		updatedAt.Set(reflect.ValueOf(helpers.Now()))
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": modelId.Interface()}, bson.M{"$set": model})
	return err
}

// Fetch data out of cursor.
// The results must be an array of database models
func (r *repository) fetchCursor(ctx context.Context, c *mongo.Cursor) (results []interface{}, err error) {
	for c.Next(ctx) {
		var m interface{}
		if m, err = r.decodeCursor(c); err != nil {
			return
		}

		results = append(results, m)
	}

	return
}

// Decode cursor data into a specified model type
func (r *repository) decodeCursor(c *mongo.Cursor) (model interface{}, err error) {
	var (
		userRepo UserRepository
	)

	switch r.collection.Name() {
	// Decode cursor data into user model
	case userRepo.CollectionName():
		var m models.User
		if err = c.Decode(&m); err != nil {
			return
		}

		model = m
		break
	default:
		err = errors.New(fmt.Sprintf("can not decode model name: %s. model unknown", r.collection.Name()))
	}

	return
}
