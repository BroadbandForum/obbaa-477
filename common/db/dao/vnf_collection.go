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
	"reflect"

	"github.com/BroadbandForum/obbaa-477/common/db/interfaces"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateVNFCollection(coll interfaces.VNFCollection, db *mongo.Database, collName string) error {
	collection, err := coll.Collection(db, collName)
	if err != nil {
		return err
	}
	if collection != nil {
		for _, doc := range coll.VNFDocuments() {
			if doc != nil {
				err = doc.Create(collection)
				if err != nil {
					return err
				}
			}
		}
	}
	for subCollName, subCollection := range coll.VNFSubCollections() {
		if !reflect.ValueOf(subCollection).IsNil() {
			err = CreateVNFCollection(subCollection, db, subCollName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func DeleteVNFCollection(coll interfaces.VNFCollection, db *mongo.Database, collName string) error {
	var err error
	for subCollName, subCollection := range coll.VNFSubCollections() {
		if !reflect.ValueOf(subCollection).IsNil() {
			err = DeleteVNFCollection(subCollection, db, subCollName)
			if err != nil {
				return err
			}
		}
	}
	collection, err := coll.Collection(db, collName)
	if err != nil {
		return err
	}
	if collection != nil {
		_, err = collection.DeleteMany(context.TODO(), bson.M{})
	}
	return err
}

func MergeVNFCollection(coll interfaces.VNFCollection, db *mongo.Database, collName string) error {
	collection, err := coll.Collection(db, collName)
	if err != nil {
		return err
	}
	if collection != nil {
		for _, doc := range coll.VNFDocuments() {
			err = ProcessVNFDocument(doc, collection)
			if err != nil {
				return err
			}
		}
	}
	for subCollName, subCollection := range coll.VNFSubCollections() {
		if !reflect.ValueOf(subCollection).IsNil() {
			err = ProcessVNFCollection(subCollection, db, subCollName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ReplaceVNFCollection(coll interfaces.VNFCollection, db *mongo.Database, collName string) error {
	err := DeleteVNFCollection(coll, db, collName)
	if err != nil {
		return err
	}
	err = CreateVNFCollection(coll, db, collName)
	return err
}

func ProcessVNFCollection(coll interfaces.VNFCollection, db *mongo.Database, collName string) error {
	var err error = nil
	if !reflect.ValueOf(coll).IsNil() {
		switch coll.Operation() {
		case interfaces.CREATE:
			err = CreateVNFCollection(coll, db, collName)
		case interfaces.DELETE:
		case interfaces.REMOVE:
			err = DeleteVNFCollection(coll, db, collName)
		case interfaces.REPLACE:
			err = ReplaceVNFCollection(coll, db, collName)
		default:
			err = MergeVNFCollection(coll, db, collName)
		}
	}
	return err
}
