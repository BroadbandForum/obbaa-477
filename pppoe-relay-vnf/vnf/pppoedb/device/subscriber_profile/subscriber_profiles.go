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
package subscriber_profile

import (
	"context"

	log "github.com/BroadbandForum/obbaa-477/common/utils/log"

	"github.com/BroadbandForum/obbaa-477/common/db/interfaces"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SubscriberProfiles struct {
	interfaces.NetconfAttributes    `bson:"-" json:"@,omitempty"`
	interfaces.DefaultVNFCollection `bson:"-" json:"-"`
	Profile                         []SubscriberProfile `json:"subscriber-profile,omitempty"`
}

func (s *SubscriberProfiles) Get(db *mongo.Database, collName string) (interfaces.VNFCollection, error) {
	collection, err := s.Collection(db, collName)
	if err != nil {
		return nil, err
	}
	profiles := []SubscriberProfile{}
	for _, profileFilter := range s.VNFDocuments() {
		log.Info("filterdoc: ", profileFilter)
		profilesResult, err := profileFilter.Get(collection)
		if err != nil {
			return nil, err
		}
		log.Info("Found doc: ", profilesResult)
		for _, profile := range profilesResult {
			profiles = append(profiles, profile.(SubscriberProfile))
		}
	}
	return &SubscriberProfiles{Profile: profiles}, nil
}

func (s *SubscriberProfiles) VNFDocuments() []interfaces.VNFDocument {
	docs := make([]interfaces.VNFDocument, len(s.Profile))
	for i := 0; i < len(s.Profile); i++ {
		docs[i] = &s.Profile[i]
	}
	return docs
}

func (s *SubscriberProfiles) Collection(db *mongo.Database, collName string) (*mongo.Collection, error) {
	collection := db.Collection(collName)
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
			{Key: "device-name", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	return collection, err
}
