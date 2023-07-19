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
package vnf_kafka

import (
	log "github.com/BroadbandForum/obbaa-477/common/utils/log"

	pb "github.com/BroadbandForum/obbaa-477/common/pb/tr451"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaHandler interface {
	//request
	Hello(*kafka.Producer, *pb.Msg)
	GetData(*kafka.Producer, *pb.Msg)
	ReplaceConfig(*kafka.Producer, *pb.Msg)
	UpdateConfig(*kafka.Producer, *pb.Msg)
	RPC(*kafka.Producer, *pb.Msg)
	Action(*kafka.Producer, *pb.Msg)

	Notification(*kafka.Producer, *pb.Msg)
}

type UnimplementedKafkaHandler struct{}

func (*UnimplementedKafkaHandler) Hello(*kafka.Producer, *pb.Msg)         {}
func (*UnimplementedKafkaHandler) GetData(*kafka.Producer, *pb.Msg)       {}
func (*UnimplementedKafkaHandler) ReplaceConfig(*kafka.Producer, *pb.Msg) {}
func (*UnimplementedKafkaHandler) UpdateConfig(*kafka.Producer, *pb.Msg)  {}
func (*UnimplementedKafkaHandler) RPC(*kafka.Producer, *pb.Msg)           {}
func (*UnimplementedKafkaHandler) Action(*kafka.Producer, *pb.Msg)        {}
func (*UnimplementedKafkaHandler) Notification(*kafka.Producer, *pb.Msg)  {}

func handleKafkaPacket(packet *pb.Msg, producer *kafka.Producer, handler KafkaHandler) {
	body := packet.Body
	if body.GetRequest() != nil {
		handleKafkaRequest(packet, producer, handler)
	} else if body.GetNotification() != nil {
		handler.Notification(producer, packet)
	} else {
		log.Error("Invalid msg body: ", body)
	}
}

func handleKafkaRequest(packet *pb.Msg, producer *kafka.Producer, handler KafkaHandler) {
	req := packet.Body.GetRequest()
	if hello := req.GetHello(); hello != nil {
		log.Info("Received a Hello message.")
		handler.Hello(producer, packet)
	} else if getData := req.GetGetData(); getData != nil {
		log.Info("Received a GetData message.")
		handler.GetData(producer, packet)
	} else if repConfig := req.GetReplaceConfig(); repConfig != nil {
		log.Info("Received a ReplaceConfig message.")
		handler.ReplaceConfig(producer, packet)
	} else if upConfig := req.GetUpdateConfig(); upConfig != nil {
		log.Info("Received an UpdateConfig message.")
		handler.UpdateConfig(producer, packet)
	} else if rpc := req.GetRpc(); rpc != nil {
		log.Info("Received a RPC message.")
		handler.RPC(producer, packet)
	} else if action := req.GetAction(); action != nil {
		log.Info("Received an Action message.")
		handler.Action(producer, packet)
	} else {
		log.Error("Invalid request: ", req)
	}
}
