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
package pppoedb

import (
	"os"

	"github.com/obbaa-477/common/db/interfaces"
	"github.com/obbaa-477/pppoe-relay-vnf/vnf/pppoedb/agent"

	"go.mongodb.org/mongo-driver/mongo"
)

var DatabaseName = os.Getenv("DB_NAME")

type PPPoEVnfJson struct {
	Agent *agent.Agent `json:"bbf-d-olt-pppoe-intermediate-agent:d-olt-pppoe-intermediate-agent"`
}

func (p *PPPoEVnfJson) Get(client *mongo.Client) (interfaces.VNFDatabase, error) {
	pppoeJson := PPPoEVnfJson{}
	tmp, err := p.Agent.Get(p.Database(client), "bbf-d-olt-pppoe-intermediate-agent:d-olt-pppoe-intermediate-agent")
	if err != nil {
		return nil, err
	}
	pppoeJson.Agent = tmp.(*agent.Agent)
	return &pppoeJson, nil
}

func (p *PPPoEVnfJson) Database(client *mongo.Client) *mongo.Database {
	return client.Database(DatabaseName)
}

func (p *PPPoEVnfJson) VNFCollections() map[string]interfaces.VNFCollection {
	return map[string]interfaces.VNFCollection{
		"bbf-d-olt-pppoe-intermediate-agent:d-olt-pppoe-intermediate-agent": p.Agent,
	}
}
