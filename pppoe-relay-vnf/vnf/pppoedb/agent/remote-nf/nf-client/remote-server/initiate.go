package remote_server

import (
	"github.com/BroadbandForum/obbaa-477/common/db/interfaces"
	"go.mongodb.org/mongo-driver/mongo"
)

type Initiate struct {
	interfaces.NetconfAttributes    `bson:"-" json:"@,omitempty"`
	interfaces.DefaultVNFCollection `bson:"-" json:"-"`
	RemoteServer                    []RemoteServer `bson:"remote-server" json:"remote-server"`
}

func (*Initiate) Collection(db *mongo.Database, collName string) (*mongo.Collection, error) {
	return db.Collection(collName), nil
}

func (in *Initiate) VNFDocuments() []interfaces.VNFDocument {
	docs := make([]interfaces.VNFDocument, len(in.RemoteServer))
	for i := 0; i < len(in.RemoteServer); i++ {
		docs[i] = &in.RemoteServer[i]
	}
	return docs
}
