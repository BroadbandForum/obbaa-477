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
	"errors"
	"fmt"

	"github.com/BroadbandForum/obbaa-477/common/db/interfaces"
	log "github.com/BroadbandForum/obbaa-477/common/utils/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Vsi struct {
	interfaces.NetconfAttributes `bson:"-" json:"@,omitempty"`
	DeviceName                   string `bson:"device-name,omitempty" json:"-"`
	VsiName                      string `bson:"vsi-name,omitempty" json:"vsi-name"`
	SubscriberProfile            string `bson:"subscriber-profile,omitempty" json:"subscriber-profile,omitempty"`
	PppoeProfile                 string `bson:"pppoe-profile,omitempty" json:"pppoe-profile,omitempty"`
	Vlans                        *Vlans `bson:"vlans,omitempty" json:"vlans,omitempty"`
}

type Vlans struct {
	CVidOnU uint16 `bson:"c-vid-on-u,omitempty" json:"c-vid-on-u,omitempty"`
	CVidOnV uint16 `bson:"c-vid-on-v,omitempty" json:"c-vid-on-v,omitempty"`
	SVidOnV uint16 `bson:"s-vid-on-v,omitempty" json:"s-vid-on-v,omitempty"`
}

func (v *Vsi) Operation() string {
	return v.NetconfAttributes.Op
}

func (v *Vsi) Create(collection *mongo.Collection) error {
	_, err := collection.InsertOne(context.TODO(), v)
	if err == nil {
		log.Info(fmt.Sprintf("Created vsi\n%+v", v))
	}
	return err
}

func (v *Vsi) Delete(collection *mongo.Collection) error {
	filter := bson.D{
		{Key: "vsi-name", Value: v.VsiName},
		{Key: "device-name", Value: v.DeviceName},
	}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if result.DeletedCount == 0 {
		return errors.New("Attempted to delete doc but it does not exist")
	}
	if err == nil {
		log.Info(fmt.Sprintf("Deleted vsi\n%+v", v))
	}
	return err
}

func (v *Vsi) Remove(collection *mongo.Collection) error {
	filter := bson.D{
		{Key: "vsi-name", Value: v.VsiName},
		{Key: "device-name", Value: v.DeviceName},
	}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err == nil {
		log.Info(fmt.Sprintf("Removed vsi\n%+v", v))
	}
	return err
}

func (v *Vsi) Merge(collection *mongo.Collection) error {
	filter := bson.D{
		{Key: "vsi-name", Value: v.VsiName},
		{Key: "device-name", Value: v.DeviceName},
	}
	update := bson.M{"$set": v}
	options := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.TODO(), filter, update, options)
	if err == nil {
		log.Info(fmt.Sprintf("Merged vsi\n%+v", v))
	}
	return err
}

func (v *Vsi) Replace(collection *mongo.Collection) error {
	filter := bson.D{
		{Key: "vsi-name", Value: v.VsiName},
		{Key: "device-name", Value: v.DeviceName},
	}
	options := options.Replace().SetUpsert(true)
	_, err := collection.ReplaceOne(context.TODO(), filter, v, options)
	if err == nil {
		log.Info(fmt.Sprintf("Replaced vsi\n%+v", v))
	}
	return err
}

func (v *Vsi) Get(collection *mongo.Collection) ([]interface{}, error) {
	bsonFilter, err := bson.Marshal(v)
	if err != nil {
		return nil, err
	}
	cursor, err := collection.Find(context.TODO(), bsonFilter)
	if err != nil {
		return nil, err
	}
	var vsis []interface{}
	for cursor.Next(context.TODO()) {
		var vsi Vsi
		if err = cursor.Decode(&vsi); err != nil {
			return nil, err
		}
		vsis = append(vsis, vsi)
	}
	return vsis, nil
}
