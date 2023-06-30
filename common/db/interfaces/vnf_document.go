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
package interfaces

import "go.mongodb.org/mongo-driver/mongo"

type VNFDocument interface {
	Operation() string
	Create(*mongo.Collection) error
	Delete(*mongo.Collection) error
	Remove(*mongo.Collection) error
	Replace(*mongo.Collection) error
	Merge(*mongo.Collection) error
	Get(*mongo.Collection) ([]interface{}, error)
}

type DefaultVNFDocument struct{}

func (*DefaultVNFDocument) Get(*mongo.Collection) ([]interface{}, error) {
	return nil, nil
}
