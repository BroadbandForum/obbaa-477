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
package device

import (
	"context"
	"errors"
	"fmt"

	"github.com/BroadbandForum/obbaa-477/common/db/dao"
	log "github.com/BroadbandForum/obbaa-477/common/utils/log"

	"github.com/BroadbandForum/obbaa-477/common/db/interfaces"

	"github.com/BroadbandForum/obbaa-477/pppoe-relay-vnf/vnf/pppoedb/device/subscriber_profile"
	"github.com/BroadbandForum/obbaa-477/pppoe-relay-vnf/vnf/pppoedb/device/vsi"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const VSI_COLL = "vsi-list"
const SUBSCRIBER_PROFILE_COLL = "subscriber-profiles"

type Device struct {
	interfaces.NetconfAttributes `bson:"-" json:"@,omitempty"`
	Name                         string                                 `bson:"name" json:"name"`
	AccessNodeId                 string                                 `bson:"access-node-id,omitempty" json:"access-node-id,omitempty"`
	VsiList                      *vsi.VsiList                           `bson:"-" json:"vsi-list,omitempty"`
	SubscriberProfiles           *subscriber_profile.SubscriberProfiles `bson:"-" json:"subscriber-profiles,omitempty"`
}

func (d *Device) Get(collection *mongo.Collection) ([]interface{}, error) {
	bsonFilter, err := bson.Marshal(d)
	if err != nil {
		return nil, err
	}
	cursor, err := collection.Find(context.TODO(), bsonFilter)
	if err != nil {
		return nil, err
	}
	var devices []interface{}
	for cursor.Next(context.TODO()) {
		var device Device
		if err = cursor.Decode(&device); err != nil {
			return nil, err
		}
		if d.VsiList != nil {
			vsiListFilter := d.VsiList
			for i := 0; i < len(vsiListFilter.Vsi); i++ {
				vsiListFilter.Vsi[i].DeviceName = d.Name
			}
			vsi_list, err := vsiListFilter.Get(collection.Database(), VSI_COLL)
			if err != nil {
				return nil, err
			}
			device.VsiList = vsi_list.(*vsi.VsiList)
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func (d *Device) Create(collection *mongo.Collection) error {
	_, err := collection.InsertOne(context.TODO(), d)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("Created device:\n%+v", d))
	d.setDeviceNameOnVsi()
	err = dao.ProcessVNFCollection(d.VsiList, collection.Database(), VSI_COLL)
	if err != nil {
		return nil
	}
	d.setDeviceNameOnSubscriberProfiles()
	err = dao.ProcessVNFCollection(d.SubscriberProfiles, collection.Database(), SUBSCRIBER_PROFILE_COLL)
	return err
}

func (d *Device) Delete(collection *mongo.Collection) error {
	err := removeVsiListOfDevice(d.Name, collection.Database())
	if err != nil {
		return err
	}
	err = removeSubscriberProfilesOfDevice(d.Name, collection.Database())
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "name", Value: d.Name}}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if result.DeletedCount == 0 {
		return errors.New("Attempted to delete doc but it does not exist")
	} else {
		log.Info(fmt.Sprintf("Deleted device:\n%+v", d))
	}
	return err
}

func (d *Device) Remove(collection *mongo.Collection) error {
	err := removeVsiListOfDevice(d.Name, collection.Database())
	if err != nil {
		return err
	}
	err = removeSubscriberProfilesOfDevice(d.Name, collection.Database())
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "name", Value: d.Name}}
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err == nil {
		log.Info(fmt.Sprintf("Removed device:\n%+v", d))
	}
	return err
}

func (d *Device) Merge(collection *mongo.Collection) error {
	filter := bson.M{"name": d.Name}
	update := bson.M{"$set": d}
	options := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.TODO(), filter, update, options)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("Merged device:\n%+v", d))
	d.setDeviceNameOnVsi()
	err = dao.ProcessVNFCollection(d.VsiList, collection.Database(), VSI_COLL)
	if err != nil {
		return err
	}
	d.setDeviceNameOnSubscriberProfiles()
	err = dao.ProcessVNFCollection(d.SubscriberProfiles, collection.Database(), SUBSCRIBER_PROFILE_COLL)
	return err
}

func (d *Device) Replace(collection *mongo.Collection) error {
	filter := bson.M{"name": d.Name}
	options := options.Replace().SetUpsert(true)
	_, err := collection.ReplaceOne(context.TODO(), filter, d, options)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("Replaced device:\n%+v", d))
	d.setDeviceNameOnVsi()
	models := []mongo.WriteModel{
		mongo.NewDeleteManyModel().SetFilter(bson.M{"device-name": d.Name}),
	}
	for _, vsi := range d.VsiList.Vsi {
		models = append(models, mongo.NewInsertOneModel().SetDocument(vsi))
	}
	_, err = collection.Database().Collection(VSI_COLL).BulkWrite(context.TODO(), models)
	if err != nil {
		return nil
	}
	d.setDeviceNameOnSubscriberProfiles()
	models = []mongo.WriteModel{
		mongo.NewDeleteManyModel().SetFilter(bson.M{"device-name": d.Name}),
	}
	for _, profile := range d.SubscriberProfiles.Profile {
		models = append(models, mongo.NewInsertOneModel().SetDocument(profile))
	}
	_, err = collection.Database().Collection(SUBSCRIBER_PROFILE_COLL).BulkWrite(context.TODO(), models)
	return err
}

func removeVsiListOfDevice(deviceName string, db *mongo.Database) error {
	collection := db.Collection(VSI_COLL)
	filter := bson.M{"device-name": deviceName}
	_, err := collection.DeleteMany(context.TODO(), filter)
	if err == nil {
		log.Info("Removed VSIs of device: ", deviceName)
	}
	return err
}

func (d *Device) setDeviceNameOnVsi() {
	if d.VsiList == nil || d.VsiList.Vsi == nil {
		return
	}
	for i := 0; i < len(d.VsiList.Vsi); i++ {
		d.VsiList.Vsi[i].DeviceName = d.Name
	}
}

func (d *Device) setDeviceNameOnSubscriberProfiles() {
	if d.SubscriberProfiles == nil || d.SubscriberProfiles.Profile == nil {
		return
	}
	for i := 0; i < len(d.SubscriberProfiles.Profile); i++ {
		d.SubscriberProfiles.Profile[i].DeviceName = d.Name
	}
}

func removeSubscriberProfilesOfDevice(deviceName string, db *mongo.Database) error {
	collection := db.Collection(SUBSCRIBER_PROFILE_COLL)
	filter := bson.M{"device-name": deviceName}
	_, err := collection.DeleteMany(context.TODO(), filter)
	if err == nil {
		log.Info("Removed Subscriber Profiles of device: ", deviceName)
	}
	return err
}
