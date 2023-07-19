package listen_endpoint

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

type ListenEndpoint struct {
	interfaces.NetconfAttributes `bson:"-" json:"@,omitempty"`
	interfaces.VNFDocument       `bson:"-" json:"-"`
	Name                         string      `bson:"name" json:"name"`
	GrpcServer                   *GrpcServer `bson:"grpc-server" json:"grpc-server"`
}

func (l *ListenEndpoint) Operation() string {
	return l.NetconfAttributes.Op
}

func (l *ListenEndpoint) Create(collection *mongo.Collection) error {
	_, err := collection.InsertOne(context.TODO(), l)
	if err == nil {
		log.Info(fmt.Sprintf("Created listen endpoint:\n%+v", l))
	}
	return err
}

func (l *ListenEndpoint) Delete(collection *mongo.Collection) error {
	filter := bson.D{
		{Key: "name", Value: l.Name},
	}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if result.DeletedCount == 0 {
		return errors.New("Attempted to delete doc but it does not exist")
	} else {
		log.Info(fmt.Sprintf("Deleted subscriber profile:\n%+v", l))
	}
	return err
}

func (l *ListenEndpoint) Remove(collection *mongo.Collection) error {
	filter := bson.D{
		{Key: "name", Value: l.Name},
	}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err == nil {
		log.Info(fmt.Sprintf("Removed subscriber profile:\n%+v", l))
	}
	return err
}

func (l *ListenEndpoint) Merge(collection *mongo.Collection) error {
	filter := bson.D{
		{Key: "name", Value: l.Name},
	}
	update := bson.M{"$set": l}
	options := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.TODO(), filter, update, options)
	if err == nil {
		log.Info(fmt.Sprintf("Merged subscriber profile:\n%+v", l))
	}
	return err
}

func (l *ListenEndpoint) Replace(collection *mongo.Collection) error {
	filter := bson.D{
		{Key: "name", Value: l.Name},
	}
	options := options.Replace().SetUpsert(true)
	_, err := collection.ReplaceOne(context.TODO(), filter, l, options)
	if err == nil {
		log.Info(fmt.Sprintf("Replaced subscriber profile:\n%+v", l))
	}
	return err
}

type GrpcServer struct {
	TCPServerParameters TCPServerParameters `bson:"tcp-server-parameters" json:"tcp-server-parameters"`
	TLSServerParameters TLSServerParameters `bson:"tls-server-parameters" json:"tls-server-parameters"`
}

type TCPServerParameters struct {
	LocalAddress string `bson:"local-address" json:"local-address"`
	LocalPort    uint32 `bson:"local-port" json:"local-port"`
}

type TLSServerParameters struct {
	ServerIdentity ServerIdentity `bson:"server-identity" json:"server-identity"`
}

type ServerIdentity struct {
	Certificate Certificate `bson:"certificate" json:"certificate"`
}

type Certificate struct {
	InlineDefinition InlineDefinition `bson:"inline-definition" json:"inline-definition"`
}

type InlineDefinition struct {
	PublicKeyFormat     string `bson:"public-key-format" json:"public-key-format"`
	PublicKey           string `bson:"public-key" json:"public-key"`
	PrivateKeyFormat    string `bson:"private-key-format" json:"private-key-format"`
	CleartextPrivateKey string `bson:"cleartext-private-key" json:"cleartext-private-key"`
	CertData            string `bson:"cert-data" json:"cert-data"`
}
