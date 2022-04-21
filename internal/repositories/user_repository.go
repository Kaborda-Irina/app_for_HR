package repositories

import (
	"context"
	"fmt"
	"github.com/inkoba/app_for_HR/internal/core/domain"
	"github.com/inkoba/app_for_HR/internal/core/ports"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var Ctx = context.Background()

type UserRepository struct {
	mc     *MongoConfig
	logger *logrus.Logger
}

var _ ports.IUserRepository = (*UserRepository)(nil)

func NewUserRepository(mc *MongoConfig, logger *logrus.Logger) ports.IUserRepository {
	return &UserRepository{
		mc,
		logger,
	}
}

func (ur UserRepository) Get(id string) (*domain.User, error) {
	var user domain.User
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ur.logger.Error("Error in repository give objectId", err)
	}
	ctx := context.Background()
	err = ur.mc.collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&user)

	if err != nil {
		ur.logger.Error("Error in find one documents in database mongodb", err)
	}

	return &user, err
}

func (ur UserRepository) GetAll() ([]*domain.User, error) {
	filter := bson.M{}

	cursor, err := ur.mc.collection.Find(Ctx, filter)
	if err != nil {
		return nil, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			ur.logger.Fatal("Get all users to db ", err)
		}
	}(cursor, Ctx)

	var results []*domain.User
	for cursor.Next(Ctx) {
		// create a value into which the single document can be decoded
		var elem domain.User
		err := cursor.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, &elem)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (ur UserRepository) Create(user *domain.User) (string, error) {
	res, err := ur.mc.collection.InsertOne(context.Background(), user)
	if err != nil {
		return "0", err
	}
	id := fmt.Sprintf("%v", res.InsertedID)
	return id, err
}

func (ur UserRepository) Delete(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = ur.mc.collection.DeleteOne(Ctx, bson.D{{"_id", objectId}})
	if err != nil {
		return err
	}
	return nil
}

func (ur UserRepository) GetUserByUsername(username string) (*domain.User, error) {
	var user *domain.User
	err := ur.mc.collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		ur.logger.Error("Error in find one documents in database mongodb", err)
		return nil, err
	}

	return user, nil
}
