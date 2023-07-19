package remotenf

import (
	"github.com/BroadbandForum/obbaa-477/common/db/interfaces"
	nfclient "github.com/BroadbandForum/obbaa-477/pppoe-relay-vnf/vnf/pppoedb/agent/remote-nf/nf-client"
	nfserver "github.com/BroadbandForum/obbaa-477/pppoe-relay-vnf/vnf/pppoedb/agent/remote-nf/nf-server"
	"go.mongodb.org/mongo-driver/mongo"
)

type RemoteNf struct {
	interfaces.NetconfAttributes    `bson:"-" json:"@,omitempty"`
	interfaces.DefaultVNFCollection `bson:"-" json:"-"`
	NfServer                        *nfserver.NfServer `bson:"nf-server,omitempty" json:"nf-server,omitempty"`
	NfClient                        *nfclient.NfClient `bson:"nf-client,omitempty" json:"nf-client,omitempty"`
}

func (r *RemoteNf) Collection(*mongo.Database, string) (*mongo.Collection, error) {
	return nil, nil
}

func (r *RemoteNf) VNFSubCollections() map[string]interfaces.VNFCollection {
	return map[string]interfaces.VNFCollection{
		"nf-server": r.NfServer,
		"nf-client": r.NfClient,
	}
}
