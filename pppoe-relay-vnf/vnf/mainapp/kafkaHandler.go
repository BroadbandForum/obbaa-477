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

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	dbDao "github.com/obbaa-477/common/db/dao"
	"google.golang.org/protobuf/proto"

	kafkaHandler "github.com/obbaa-477/common/kafkaHandler"
	pb "github.com/obbaa-477/common/pb/tr451"
	log "github.com/obbaa-477/common/utils/log"
	pppoedb "github.com/obbaa-477/pppoe-relay-vnf/vnf/pppoedb"
)

var (
	kafkaHost         = "kafka"
	kafkaPort         = 9092
	kafkaConsumeTopic = fmt.Sprintf("%s-request", os.Getenv("VNF_NAME"))
	kafkaProduceTopic = fmt.Sprintf("%s-response", os.Getenv("VNF_NAME"))
)

type pppoeVnfKafkaHandler struct {
	kafkaHandler.UnimplementedKafkaHandler
}

func getKafkaHostAndPort() {
	var err error
	if os.Getenv("KAFKA_HOST") != "" {
		kafkaHost = os.Getenv("KAFKA_HOST")
	}
	if os.Getenv("KAFKA_PORT") != "" {
		kafkaPort, err = strconv.Atoi(os.Getenv("KAFKA_PORT"))
		if err != nil {
			log.Error("Invalid KAFKA_PORT value ", os.Getenv("KAFKA_PORT"))
			os.Exit(2)
		}
	}
	log.Info("kafka bootstrap.server=" + kafkaHost + ":" + strconv.Itoa(kafkaPort))
}

func (*pppoeVnfKafkaHandler) GetData(producer *kafka.Producer, packet *pb.Msg) {
	for _, filter := range packet.Body.GetRequest().GetGetData().Filter {
		pppoeJson := []byte{}
		statusCode := pb.Status_OK
		pppoeJsonFilter := pppoedb.PPPoEVnfJson{}
		err := json.Unmarshal(filter, &pppoeJsonFilter)
		if err != nil {
			log.Error("Error marshaling json: ", err)
			statusCode = pb.Status_ERROR_GENERAL
		} else {
			pppoeJsonResult, err := pppoeJsonFilter.Get(mongoClient)
			if err != nil {
				log.Error("Error getting data from database: ", err)
				statusCode = pb.Status_ERROR_GENERAL
			}
			pppoeJson, err = json.Marshal(pppoeJsonResult)
			if err != nil {
				log.Error("Error marshaling json: ", err)
				statusCode = pb.Status_ERROR_GENERAL
			}
		}

		msg := pb.Msg{
			Header: &pb.Header{
				MsgId:         packet.Header.MsgId,
				SenderName:    packet.Header.RecipientName,
				RecipientName: packet.Header.SenderName,
				ObjectType:    packet.Header.ObjectType,
				ObjectName:    packet.Header.ObjectName,
			},
			Body: &pb.Body{
				MsgBody: &pb.Body_Response{
					Response: &pb.Response{
						RespType: &pb.Response_GetResp{
							GetResp: &pb.GetDataResp{
								StatusResp: &pb.Status{
									StatusCode: statusCode,
								},
								Data: pppoeJson,
							},
						},
					},
				},
			},
		}
		packetBytes, err := proto.Marshal(&msg)
		if err != nil {
			log.Error("Error proto marshaling doc: ", err)
			return
		}
		producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &kafkaProduceTopic,
				Partition: kafka.PartitionAny,
			},
			Value: packetBytes,
		},
			nil)

	}
}

func (*pppoeVnfKafkaHandler) UpdateConfig(producer *kafka.Producer, packet *pb.Msg) {
	statusCode := pb.Status_OK
	var deltaConfig pppoedb.PPPoEVnfJson
	deltaBytes := getDeltaConfig(packet.Body.GetRequest().GetUpdateConfig())

	log.Info("Handling UpdateConfig message")
	json.Unmarshal(deltaBytes, &deltaConfig)
	err := dbDao.ProcessVNFDatabase(&deltaConfig, mongoClient)
	if err != nil {
		log.Error("Error during update config: ", err)
		statusCode = pb.Status_ERROR_GENERAL
	}
	msg := pb.Msg{
		Header: &pb.Header{
			MsgId:         packet.Header.MsgId,
			SenderName:    packet.Header.RecipientName,
			RecipientName: packet.Header.SenderName,
			ObjectName:    packet.Header.ObjectName,
			ObjectType:    packet.Header.ObjectType,
		},
		Body: &pb.Body{
			MsgBody: &pb.Body_Response{
				Response: &pb.Response{
					RespType: &pb.Response_UpdateConfigResp{
						UpdateConfigResp: &pb.UpdateConfigResp{
							StatusResp: &pb.Status{
								StatusCode: statusCode,
							},
						},
					},
				},
			},
		},
	}
	packetBytes, err := proto.Marshal(&msg)
	if err != nil {
		log.Error("Error proto marshaling msg: ", err)
		return
	}
	producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &kafkaProduceTopic,
			Partition: kafka.PartitionAny,
		},
		Value: packetBytes,
	}, nil)
}

func getDeltaConfig(updateConfig *pb.UpdateConfig) []byte {
	if updateConfigInst := updateConfig.GetUpdateConfigInst(); updateConfigInst != nil {
		log.Info("Update Config Instance Request")
		return updateConfigInst.DeltaConfig
	} else if updateConfigReplica := updateConfig.GetUpdateConfigReplica(); updateConfigReplica != nil {
		log.Info("Update Config Replica Request")
		return updateConfigReplica.DeltaConfig
	}
	log.Error("Could not get any config from request.")
	return []byte{}
}

func (*pppoeVnfKafkaHandler) ReplaceConfig(producer *kafka.Producer, packet *pb.Msg) {
	statusCode := pb.Status_OK
	var deltaConfig pppoedb.PPPoEVnfJson
	deltaBytes := packet.Body.GetRequest().GetReplaceConfig().ConfigInst

	log.Info("Handling ReplaceConfig message")
	json.Unmarshal(deltaBytes, &deltaConfig)
	err := dbDao.ProcessVNFDatabase(&deltaConfig, mongoClient)
	if err != nil {
		log.Error("Error during replace config: ", err)
		statusCode = pb.Status_ERROR_GENERAL
	}
	msg := pb.Msg{
		Header: &pb.Header{
			MsgId:         packet.Header.MsgId,
			SenderName:    packet.Header.RecipientName,
			RecipientName: packet.Header.SenderName,
			ObjectName:    packet.Header.ObjectName,
			ObjectType:    packet.Header.ObjectType,
		},
		Body: &pb.Body{
			MsgBody: &pb.Body_Response{
				Response: &pb.Response{
					RespType: &pb.Response_ReplaceConfigResp{
						ReplaceConfigResp: &pb.ReplaceConfigResp{
							StatusResp: &pb.Status{
								StatusCode: statusCode,
							},
						},
					},
				},
			},
		},
	}
	packetBytes, err := proto.Marshal(&msg)
	if err != nil {
		log.Error("Error proto marshaling msg: ", err)
		return
	}
	producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &kafkaProduceTopic,
			Partition: kafka.PartitionAny,
		},
		Value: packetBytes,
	}, nil)
}

func (*pppoeVnfKafkaHandler) Hello(producer *kafka.Producer, packet *pb.Msg) {

	log.Info("Handling Hello message")
	msg := pb.Msg{
		Header: &pb.Header{
			MsgId:         packet.Header.MsgId,
			SenderName:    packet.Header.RecipientName,
			RecipientName: packet.Header.SenderName,
			ObjectName:    packet.Header.ObjectName,
			ObjectType:    packet.Header.ObjectType,
		},
		Body: &pb.Body{
			MsgBody: &pb.Body_Response{
				Response: &pb.Response{
					RespType: &pb.Response_HelloResp{
						HelloResp: &pb.HelloResp{
							ServiceEndpointName: os.Getenv("VNF_NAME"),
							NetworkFunctionInfo: []*pb.NFInformation{
								{
									NfTypes: map[string]string{"vendor-name": "Broadband Forum"},
								},
							},
						},
					},
				},
			},
		},
	}
	packetBytes, err := proto.Marshal(&msg)
	if err != nil {
		log.Error("Error proto marshaling msg: ", err)
		return
	}
	producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &kafkaProduceTopic,
			Partition: kafka.PartitionAny,
		},
		Value: packetBytes,
	}, nil)
}
