package remote_server

import (
	"context"
	"errors"
	"fmt"

	"github.com/obbaa-477/common/db/interfaces"
	log "github.com/obbaa-477/common/utils/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RemoteServer struct {
	interfaces.NetconfAttributes `bson:"-" json:"@,omitempty"`
	interfaces.VNFDocument       `bson:"-" json:"-"`
	Name                         string     `bson:"name" json:"name"`
	NfType                       string     `bson:"nf-type" json:"nf-type"`
	OnDemand                     bool       `bson:"on-demand" json:"on-demand"`
	GrpcClient                   GrpcClient `bson:"grpc-client" json:"grpc-client"`
	MfcType                      string     `bson:"mfc-type" json:"mfc-type"`
}

func (r *RemoteServer) Operation() string {
	return r.NetconfAttributes.Op
}

func (r *RemoteServer) Create(collection *mongo.Collection) error {
	_, err := collection.InsertOne(context.TODO(), r)
	if err == nil {
		log.Info(fmt.Sprintf("Created listen endpoint:\n%+v", r))
	}
	return err
}

func (r *RemoteServer) Delete(collection *mongo.Collection) error {
	filter := bson.D{
		{Key: "name", Value: r.Name},
	}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if result.DeletedCount == 0 {
		return errors.New("Attempted to delete doc but it does not exist")
	} else {
		log.Info(fmt.Sprintf("Deleted subscriber profile:\n%+v", r))
	}
	return err
}

func (r *RemoteServer) Remove(collection *mongo.Collection) error {
	filter := bson.D{
		{Key: "name", Value: r.Name},
	}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err == nil {
		log.Info(fmt.Sprintf("Removed subscriber profile:\n%+v", r))
	}
	return err
}

func (r *RemoteServer) Merge(collection *mongo.Collection) error {
	filter := bson.D{
		{Key: "name", Value: r.Name},
	}
	update := bson.M{"$set": r}
	options := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.TODO(), filter, update, options)
	if err == nil {
		log.Info(fmt.Sprintf("Merged subscriber profile:\n%+v", r))
	}
	return err
}

func (r *RemoteServer) Replace(collection *mongo.Collection) error {
	filter := bson.D{
		{Key: "name", Value: r.Name},
	}
	options := options.Replace().SetUpsert(true)
	_, err := collection.ReplaceOne(context.TODO(), filter, r, options)
	if err == nil {
		log.Info(fmt.Sprintf("Replaced subscriber profile:\n%+v", r))
	}
	return err
}

type GrpcClient struct {
	AccessPoint []AccessPoint `bson:"access-point" json:"access-point"`
}

type AccessPoint struct {
	Name                    string                  `bson:"name" json:"name"`
	GrpcTransportParameters GrpcTransportParameters `bson:"grpc-transport-parameters" json:"grpc-transport-parameters"`
}

type GrpcTransportParameters struct {
	TCPClientParameters TCPClientParameters `bson:"tcp-client-parameters" json:"tcp-client-parameters"`
	TLSClientParameters TLSClientParameters `bson:"tls-client-parameters" json:"tls-client-parameters"`
}

type TCPClientParameters struct {
	RemoteAddress string `bson:"remote-address" json:"remote-address"`
	RemotePort    uint32 `bson:"remote-port" json:"remote-port"`
}

type TLSClientParameters struct {
	ServerAuthentication ServerAuthentication `bson:"server-authentication" json:"server-authentication"`
}

type ServerAuthentication struct {
	CaCerts CaCerts `bson:"ca-certs" json:"ca-certs"`
}

type CaCerts struct {
	InlineDefinition InlineDefinition `bson:"inline-definition" json:"inline-definition"`
}

type InlineDefinition struct {
	Certificate []Certificate `bson:"certificate" json:"certificate"`
}

type Certificate struct {
	Name     string `bson:"name" json:"name"`
	CertData string `bson:"cert-data" json:"cert-data"`
}
