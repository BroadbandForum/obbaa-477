/*
 * Copyright 2023 Broadband Forum
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
* PPPOE VNF main file
*
* Created by João Correia(Altice Labs) on 20/03/2023
 */
package agent

import (
	"reflect"

	"github.com/BroadbandForum/obbaa-477/common/db/interfaces"
	remotenf "github.com/BroadbandForum/obbaa-477/pppoe-relay-vnf/vnf/pppoedb/agent/remote-nf"
	"github.com/BroadbandForum/obbaa-477/pppoe-relay-vnf/vnf/pppoedb/device"
	"github.com/BroadbandForum/obbaa-477/pppoe-relay-vnf/vnf/pppoedb/pppoe_profile"

	"go.mongodb.org/mongo-driver/mongo"
)

type Agent struct {
	interfaces.NetconfAttributes    `bson:"-" json:"@,omitempty"`
	interfaces.DefaultVNFCollection `bson:"-" json:"-"`
	PPPoEProfiles                   *pppoe_profile.PPPoEProfiles `json:"pppoe-profiles,omitempty"`
	Devices                         *device.Devices              `json:"devices-using-d-olt-pppoeia,omitempty"`
	RemoteNf                        *remotenf.RemoteNf           `json:"remote-nf,omitempty"`
}

func (a *Agent) Collection(db *mongo.Database, collName string) (*mongo.Collection, error) {
	return nil, nil
}

func (a *Agent) VNFSubCollections() map[string]interfaces.VNFCollection {
	m := map[string]interfaces.VNFCollection{
		"pppoe-profiles":              a.PPPoEProfiles,
		"devices-using-d-olt-pppoeia": a.Devices,
		"remote-nf":                   a.RemoteNf,
	}
	for key, value := range m {
		if reflect.ValueOf(value).IsNil() {
			delete(m, key)
		}
	}
	return m
}

func (a *Agent) Get(db *mongo.Database, collName string) (interfaces.VNFCollection, error) {
	var agent Agent
	if a.PPPoEProfiles != nil {
		tmp, err := a.PPPoEProfiles.Get(db, "pppoe-profiles")
		if err != nil {
			return nil, err
		}
		agent.PPPoEProfiles = tmp.(*pppoe_profile.PPPoEProfiles)

	}
	if a.Devices != nil {
		tmp, err := a.Devices.Get(db, "devices-using-d-olt-pppoeia")
		if err != nil {
			return nil, err
		}
		agent.Devices = tmp.(*device.Devices)
	}

	return &agent, nil
}
