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
	"errors"
	"fmt"

	"github.com/obbaa-477/common/db/interfaces"
	log "github.com/obbaa-477/common/utils/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PPPoEProfile struct {
	interfaces.NetconfAttributes `bson:"-" json:"@,omitempty"`
	Name                         string                  `bson:"name" json:"name"`
	PppoeVendorSpecificTag       *PPPoEVendorSpecificTag `bson:"pppoe-vendor-specific-tag,omitempty" json:"pppoe-vendor-specific-tag,omitempty"`
}

type PPPoEVendorSpecificTag struct {
	Subtag                 []string `bson:"subtag,omitempty" json:"subtag,omitempty"`
	DefaultCircuitIDSyntax string   `bson:"default-circuit-id-syntax,omitempty" json:"default-circuit-id-syntax,omitempty"`
	DefaultRemoteIDSyntax  string   `bson:"default-remote-id-syntax,omitempty" json:"default-remote-id-syntax,omitempty"`
}

func (p *PPPoEProfile) Get(collection *mongo.Collection) ([]interface{}, error) {
	bsonFilter, err := bson.Marshal(p)
	if err != nil {
		return nil, err
	}
	cursor, err := collection.Find(context.TODO(), bsonFilter)
	if err != nil {
		return nil, err
	}
	var profiles []interface{}
	for cursor.Next(context.TODO()) {
		var device PPPoEProfile
		if err = cursor.Decode(&device); err != nil {
			return nil, err
		}
		profiles = append(profiles, device)
	}
	return profiles, nil

}

func (p *PPPoEProfile) Create(collection *mongo.Collection) error {
	_, err := collection.InsertOne(context.TODO(), p)
	if err == nil {
		log.Info(fmt.Sprintf("Created pppoe profile:\n%+v", p))
	}
	return err
}

func (p *PPPoEProfile) Delete(collection *mongo.Collection) error {
	filter := bson.D{{Key: "name", Value: p.Name}}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if result.DeletedCount == 0 {
		return errors.New("Attempted to delete doc but it does not exist")
	} else {
		log.Info(fmt.Sprintf("Deleted subscriber profile:\n%+v", p))
	}
	return err
}

func (p *PPPoEProfile) Remove(collection *mongo.Collection) error {
	filter := bson.D{{Key: "name", Value: p.Name}}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err == nil {
		log.Info(fmt.Sprintf("Removed pppoe profile:\n%+v", p))
	}
	return err
}

func (p *PPPoEProfile) Merge(collection *mongo.Collection) error {
	filter := bson.M{"name": p.Name}
	update := bson.M{"$set": p}
	options := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.TODO(), filter, update, options)
	if err == nil {
		log.Info(fmt.Sprintf("Merged pppoe profile:\n%+v", p))
	}
	return err
}

func (p *PPPoEProfile) Replace(collection *mongo.Collection) error {
	filter := bson.M{"name": p.Name}
	options := options.Replace().SetUpsert(true)
	_, err := collection.ReplaceOne(context.TODO(), filter, p, options)
	if err == nil {
		log.Info(fmt.Sprintf("Replaced pppoe profile:\n%+v", p))
	}
	return err
}
