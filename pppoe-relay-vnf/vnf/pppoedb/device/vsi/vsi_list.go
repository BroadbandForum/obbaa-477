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
package vsi

import (
	"context"

	"github.com/BroadbandForum/obbaa-477/common/db/interfaces"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VsiList struct {
	interfaces.NetconfAttributes    `bson:"-" json:"@,omitempty"`
	interfaces.DefaultVNFCollection `bson:"-" json:"-"`
	Vsi                             []Vsi `json:"vsi,omitempty"`
}

func (v *VsiList) Collection(db *mongo.Database, collName string) (*mongo.Collection, error) {
	collection := db.Collection(collName)
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "vsi-name", Value: 1},
			{Key: "device-name", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	return collection, err
}

func (v *VsiList) VNFDocuments() []interfaces.VNFDocument {
	docs := make([]interfaces.VNFDocument, len(v.Vsi))
	for i := 0; i < len(v.Vsi); i++ {
		docs[i] = &v.Vsi[i]
	}
	return docs
}

func (v *VsiList) Get(db *mongo.Database, collName string) (interfaces.VNFCollection, error) {
	collection, err := v.Collection(db, collName)
	if err != nil {
		return nil, err
	}
	vsis := []Vsi{}
	for _, vsiFilter := range v.VNFDocuments() {
		vsisResult, err := vsiFilter.Get(collection)
		if err != nil {
			return nil, err
		}
		for _, vsi := range vsisResult {
			vsis = append(vsis, vsi.(Vsi))
		}
	}
	return &VsiList{Vsi: vsis}, nil
}
