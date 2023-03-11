package repositories

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/nicodeheza/peersEat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PeerRepositoryI interface{
	Insert(peer models.Peer) (id primitive.ObjectID ,  err error)
	GetById(id  primitive.ObjectID)(models.Peer, error)
	GetAll(excludesUrls []string) ([]models.Peer, error)
	GetSelf() (models.Peer, error)
	Update(peer models.Peer, fields []string ) error
	GetAllUrls(excludes []string) ([]string, error)

}

type PeerRepository struct{
	coll *mongo.Collection
}

func NewPeerRepository(collection *mongo.Collection) *PeerRepository{
	return &PeerRepository{collection}
}



func (p PeerRepository) Insert(peer models.Peer) (id primitive.ObjectID ,  err error) {
	result,err := p.coll.InsertOne(context.Background(), peer)
	if err != nil{ return primitive.NewObjectID(), err}
	return  result.InsertedID.(primitive.ObjectID) , nil
}

func (p PeerRepository) GetById(id  primitive.ObjectID)(models.Peer, error){
	filter := bson.D{{Key:"_id", Value: id}}
	var result models.Peer
	err := p.coll.FindOne(context.Background(), filter).Decode(&result)
	
	return	result, err
}

func (p PeerRepository) GetAll(excludesUrls []string) ([]models.Peer, error) {
	filter := bson.D{}

	for _, excludeUrl := range excludesUrls{
		filter = append(filter, bson.E{Key: "url" , Value: bson.D{{Key: "$ne", Value: excludeUrl}}})
	}

	cursor, err := p.coll.Find(context.Background(), filter)
	if err != nil{return nil, err}
	var results []models.Peer
	if err= cursor.All(context.Background(), &results); err != nil{return nil, err}
	return results, nil
}

func (p PeerRepository) GetSelf() (models.Peer, error) {

	var result models.Peer

	filter:= bson.D{{Key: "url", Value: os.Getenv("HOST")}}

	err := p.coll.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		return models.Peer{}, err
	}

	return result, nil
} 

func (p PeerRepository) Update(peer models.Peer, fields []string ) error {
	filter := bson.D{{Key: "_id",  Value: peer.Id}}

	updateDatan := bson.D{}
	val:= reflect.ValueOf(peer)

	for _, field := range fields{
		f:= reflect.Indirect(val).FieldByName(field)
		if  f.IsZero(){
			message := fmt.Sprintf("field %v not exist in peer struct", field)
			return errors.New(message)
		}
		updateDatan = append(updateDatan, bson.E{Key: field,  Value: f})
	}

	update:= bson.D{{Key: "$set", Value: updateDatan}}

	_, err := p.coll.UpdateOne(context.Background(), filter, update)

	if err != nil{return err}
	return nil
}

func (p PeerRepository) GetAllUrls(excludes []string) ([]string, error){
	opts := options.Find().SetProjection(bson.D{{Key: "url", Value: 1}})
	filter := bson.D{}

	for _, excludeUrl := range excludes{
		filter = append(filter, bson.E{Key: "url" , Value: bson.D{{Key: "$ne", Value: excludeUrl}}})
	}

	cursor, err := p.coll.Find(context.Background(), filter, opts)
	if err != nil{
		return nil, err
	}

	var result []string
	if err = cursor.All(context.Background(), &result); err !=nil{
		return nil, err
	}

	return result, nil
}
