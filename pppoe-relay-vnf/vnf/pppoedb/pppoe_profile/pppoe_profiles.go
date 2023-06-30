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
package pppoe_profile

import (
	"context"

	"github.com/obbaa-477/common/db/interfaces"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PPPoEProfiles struct {
	interfaces.NetconfAttributes    `bson:"-" json:"@,omitempty"`
	interfaces.DefaultVNFCollection `bson:"-" json:"-"`
	Profile                         []PPPoEProfile `json:"pppoe-profile,omitempty"`
}

func (p *PPPoEProfiles) Get(db *mongo.Database, collName string) (interfaces.VNFCollection, error) {
	collection, err := p.Collection(db, collName)
	if err != nil {
		return nil, err
	}
	profiles := []PPPoEProfile{}
	for _, profileFilter := range p.VNFDocuments() {
		profilesResult, err := profileFilter.Get(collection)
		if err != nil {
			return nil, err
		}
		for _, profile := range profilesResult {
			profiles = append(profiles, profile.(PPPoEProfile))
		}
	}
	return &PPPoEProfiles{Profile: profiles}, nil
}

func (p *PPPoEProfiles) VNFDocuments() []interfaces.VNFDocument {
	docs := make([]interfaces.VNFDocument, len(p.Profile))
	for i := 0; i < len(p.Profile); i++ {
		docs[i] = &p.Profile[i]
	}
	return docs
}

func (p *PPPoEProfiles) Collection(db *mongo.Database, collName string) (*mongo.Collection, error) {
	collection := db.Collection(collName)
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	return collection, err
}
