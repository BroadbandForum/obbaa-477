package nfclient

import (
	"context"

	"github.com/obbaa-477/common/db/interfaces"
	remote_server "github.com/obbaa-477/pppoe-relay-vnf/vnf/pppoedb/agent/remote-nf/nf-client/remote-server"
	remotenf_common "github.com/obbaa-477/pppoe-relay-vnf/vnf/pppoedb/agent/remote-nf/remote-nf-common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NfClient struct {
	interfaces.NetconfAttributes    `bson:"-" json:"@,omitempty"`
	interfaces.DefaultVNFCollection `bson:"-" json:"-"`
	Enabled                         bool                    `bson:"enabled" json:"enabled"`
	Initiate                        *remote_server.Initiate `bson:"initiate" json:"initiate"`
}

func (nf *NfClient) VNFDocuments() []interfaces.VNFDocument {
	return []interfaces.VNFDocument{&remotenf_common.Enabled{
		NetconfAttributes: nf.NetconfAttributes,
		Enabled:           nf.Enabled,
	}}
}

func (nf *NfClient) VNFSubCollections() map[string]interfaces.VNFCollection {
	return map[string]interfaces.VNFCollection{
		"nf-client": nf.Initiate,
	}
}

func (nf *NfClient) Collection(db *mongo.Database, collName string) (*mongo.Collection, error) {
	// Access the collection on which you want to create the index
	collection := db.Collection(collName)

	// Specify the index model with the unique and sparse options
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetSparse(true),
	}

	// Create the index
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return collection, err
}
