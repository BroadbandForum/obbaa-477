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
package dao

import (
	"context"

	"github.com/BroadbandForum/obbaa-477/common/db/interfaces"

	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetVNFDocuments[D interfaces.VNFDocument](filter D, db *mongo.Database, collName string) ([]D, error) {
	bsonFilter, err := bson.Marshal(filter)
	if err != nil {
		return nil, err
	}
	collection := db.Collection(collName)
	cursor, err := collection.Find(context.TODO(), bsonFilter)
	if err != nil {
		return nil, err
	}
	var docs []D
	for cursor.Next(context.TODO()) {
		var doc D
		if err = cursor.Decode(&doc); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

func GetVNFDocument[D interfaces.VNFDocument](filter D, db *mongo.Database, collName string) (D, error) {
	bsonFilter, err := bson.Marshal(filter)
	var doc D
	if err != nil {
		return doc, err
	}
	collection := db.Collection(collName)
	err = collection.FindOne(context.TODO(), bsonFilter).Decode(&doc)
	return doc, err
}

func ProcessVNFDatabase(d interfaces.VNFDatabase, client *mongo.Client) error {
	var err error
	database := d.Database(client)
	for collName, coll := range d.VNFCollections() {
		if !reflect.ValueOf(coll).IsNil() {
			err = ProcessVNFCollection(coll, database, collName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
