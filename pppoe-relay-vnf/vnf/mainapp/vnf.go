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
* Created by João Correia(Altice Labs) on 02/02/2023
 */
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	kafkaHandler "github.com/BroadbandForum/obbaa-477/common/kafkaHandler"
	kafka "github.com/confluentinc/confluent-kafka-go/kafka"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	log "github.com/BroadbandForum/obbaa-477/common/utils/log"

	pppoePacket "github.com/BroadbandForum/obbaa-477/pppoe-relay-vnf/vnf/pppoepacket"
)

var socketGrpc = os.Getenv("SOCKET_GRPC") // "0.0.0.0:50051" SDN_MC_SERVER_PORT=50053
//TODO: make this a map to support multiple devices
// var globalStream pb.CpriMessage_ListenForCpriRxServer

var mongoClient *mongo.Client

var (
	mongoHost = "mongo"
	mongoPort = 27017
)
var discard_on_error bool = true
var vnfMode = os.Getenv("VNF_MODE")

func main() {
	var err error
	if discard_flag, flag_is_set := os.LookupEnv("DISCARD_ON_ERROR"); flag_is_set {
		discard_on_error, err = strconv.ParseBool(discard_flag)
		if err != nil {
			log.Fatal("VNF - Failed to parse DISCARD_ON_ERROR flag: ", err)
		}
	}
	if os.Getenv("MONGO_HOST") != "" {
		mongoHost = os.Getenv("MONGO_HOST")
	}
	if os.Getenv("MONGO_PORT") != "" {
		mongoPort, err = strconv.Atoi(os.Getenv("MONGO_PORT"))
		if err != nil {
			log.Error("Invalid MONGO_PORT value ", os.Getenv("MONGO_PORT"))
			os.Exit(2)
		}
	}
	var mongoUri = fmt.Sprintf("mongodb://%s:%d", mongoHost, mongoPort)

	mongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatal("VNF - Failed to connect to db: ", err)
	}

	pppoePacket.InitPPPoEDecoder()

	getKafkaHostAndPort()
	kafkaConfig := kafkaHandler.KafkaConnectionConfig{
		CTopics: []string{kafkaConsumeTopic},
		CConfigMap: &kafka.ConfigMap{
			"bootstrap.servers": kafkaHost + ":" + strconv.Itoa(kafkaPort),
			"group.id":          "pppoe-vnf-consumers",
		},
		PConfigMap: &kafka.ConfigMap{
			"bootstrap.servers": kafkaHost + ":" + strconv.Itoa(kafkaPort),
		},
	}
	kafkaHandler.StartKafkaConnection(kafkaConfig, &pppoeVnfKafkaHandler{})
	run()
}

func run() {
	for {
		switch vnfMode {
		case "server":
			log.Info("Running as server")
			runAsServer()
		case "client":
			log.Info("Running as client")
			runAsClient()
		default:
			log.Fatal("Unknown VNF_MODE: ", vnfMode)
		}
	}
}
