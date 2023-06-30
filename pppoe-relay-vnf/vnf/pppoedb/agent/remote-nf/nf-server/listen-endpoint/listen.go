package listen_endpoint

import (
	"github.com/obbaa-477/common/db/interfaces"
	"go.mongodb.org/mongo-driver/mongo"
)

type Listen struct {
	interfaces.NetconfAttributes    `bson:"-" json:"@,omitempty"`
	interfaces.DefaultVNFCollection `bson:"-" json:"-"`
	ListenEndpoint                  []ListenEndpoint `json:"listen-endpoint"`
}

func (*Listen) Collection(db *mongo.Database, collName string) (*mongo.Collection, error) {
	return db.Collection(collName), nil
}

func (l *Listen) VNFDocuments() []interfaces.VNFDocument {
	docs := make([]interfaces.VNFDocument, len(l.ListenEndpoint))
	for i := 0; i < len(l.ListenEndpoint); i++ {
		docs[i] = &l.ListenEndpoint[i]
	}
	return docs
}
