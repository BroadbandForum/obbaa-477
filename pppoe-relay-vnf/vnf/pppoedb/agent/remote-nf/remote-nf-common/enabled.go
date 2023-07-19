package remotenf_common

import (
	"context"
	"errors"
	"fmt"

	"github.com/BroadbandForum/obbaa-477/common/db/interfaces"
	log "github.com/BroadbandForum/obbaa-477/common/utils/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Enabled struct {
	interfaces.NetconfAttributes  `bson:"-" json:"@,omitempty"`
	interfaces.DefaultVNFDocument `bson:"-" json:"-"`
	Enabled                       bool `bson:"enabled"`
}

func (e *Enabled) Merge(collection *mongo.Collection) error {
	update := bson.M{"$set": e}
	options := options.Update().SetUpsert(true)
	filter := bson.M{
		"enabled": bson.M{
			"$exists": true,
		},
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update, options)
	if err == nil {
		log.Info(fmt.Sprintf("Merged agent doc\n%+v", e))
	}
	return err
}

func (e *Enabled) Create(collection *mongo.Collection) error {
	return e.Merge(collection)
}

func (e *Enabled) Replace(collection *mongo.Collection) error {
	return e.Merge(collection)
}

func (e *Enabled) Delete(collection *mongo.Collection) error {
	filter := bson.M{
		"enabled": bson.M{
			"$exists": true,
		},
	}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if result.DeletedCount == 0 {
		return errors.New("Attempted to delete doc but it does not exist")
	} else {
		log.Info(fmt.Sprintf("Deleted subscriber profile:\n%+v", e))
	}
	return err
}

func (e *Enabled) Remove(collection *mongo.Collection) error {
	return e.Delete(collection)
}
