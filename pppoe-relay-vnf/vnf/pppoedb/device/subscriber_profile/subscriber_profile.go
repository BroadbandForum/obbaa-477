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
	"errors"
	"fmt"

	"github.com/obbaa-477/common/db/interfaces"
	log "github.com/obbaa-477/common/utils/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SubscriberProfile struct {
	interfaces.NetconfAttributes `bson:"-" json:"@,omitempty"`
	Name                         string `bson:"name" json:"name,omitempty"`
	DeviceName                   string `bson:"device-name" json:"-"`
	CircuitID                    string `bson:"circuit-id,omitempty" json:"circuit-id,omitempty"`
	RemoteID                     string `bson:"remote-id,omitempty" json:"remote-id,omitempty"`
}

func (s *SubscriberProfile) Get(collection *mongo.Collection) ([]interface{}, error) {
	bsonFilter, err := bson.Marshal(s)
	if err != nil {
		return nil, err
	}
	cursor, err := collection.Find(context.TODO(), bsonFilter)
	if err != nil {
		return nil, err
	}
	var profiles []interface{}
	for cursor.Next(context.TODO()) {
		var profile SubscriberProfile
		if err = cursor.Decode(&profile); err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil

}

func (s *SubscriberProfile) Create(collection *mongo.Collection) error {
	_, err := collection.InsertOne(context.TODO(), s)
	if err == nil {
		log.Info(fmt.Sprintf("Created subscriber profile:\n%+v", s))
	}
	return err
}

func (s *SubscriberProfile) Delete(collection *mongo.Collection) error {
	filter := bson.D{
		{Key: "name", Value: s.Name},
		{Key: "device-name", Value: s.DeviceName},
	}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if result.DeletedCount == 0 {
		return errors.New("Attempted to delete doc but it does not exist")
	} else {
		log.Info(fmt.Sprintf("Deleted subscriber profile:\n%+v", s))
	}
	return err
}

func (s *SubscriberProfile) Remove(collection *mongo.Collection) error {
	filter := bson.D{
		{Key: "name", Value: s.Name},
		{Key: "device-name", Value: s.DeviceName},
	}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err == nil {
		log.Info(fmt.Sprintf("Removed subscriber profile:\n%+v", s))
	}
	return err
}

func (s *SubscriberProfile) Merge(collection *mongo.Collection) error {
	filter := bson.D{
		{Key: "name", Value: s.Name},
		{Key: "device-name", Value: s.DeviceName},
	}
	update := bson.M{"$set": s}
	options := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.TODO(), filter, update, options)
	if err == nil {
		log.Info(fmt.Sprintf("Merged subscriber profile:\n%+v", s))
	}
	return err
}

func (s *SubscriberProfile) Replace(collection *mongo.Collection) error {
	filter := bson.D{
		{Key: "name", Value: s.Name},
		{Key: "device-name", Value: s.DeviceName},
	}
	options := options.Replace().SetUpsert(true)
	_, err := collection.ReplaceOne(context.TODO(), filter, s, options)
	if err == nil {
		log.Info(fmt.Sprintf("Replaced subscriber profile:\n%+v", s))
	}
	return err
}
