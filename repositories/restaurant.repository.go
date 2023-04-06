package repositories

import (
	"context"

	"github.com/nicodeheza/peersEat/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RestaurantRepositoryI interface {
	Insert(restaurant models.Restaurant) (id primitive.ObjectID, err error)
}

type RestaurantRepository struct {
	coll *mongo.Collection
}

func NewRestaurantRepository(collection *mongo.Collection) *RestaurantRepository {
	return &RestaurantRepository{collection}
}

func (r *RestaurantRepository) Insert(restaurant models.Restaurant) (id primitive.ObjectID, err error) {
	result, err := r.coll.InsertOne(context.Background(), restaurant)
	if err != nil {
		return primitive.NewObjectID(), err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}
