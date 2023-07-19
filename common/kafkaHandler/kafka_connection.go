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
	"time"

	pb "github.com/BroadbandForum/obbaa-477/common/pb/tr451"
	log "github.com/BroadbandForum/obbaa-477/common/utils/log"
	kafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/protobuf/proto"
)

type KafkaConnectionConfig struct {
	PConfigMap *kafka.ConfigMap
	CConfigMap *kafka.ConfigMap
	CTopics    []string
}

func StartKafkaConnection(config KafkaConnectionConfig, handler KafkaHandler) {
	go func() {
		for {
			tryKafkaConnection(config, handler)
			log.Info("Retrying kafka connection in 10 seconds")
			time.Sleep(10 * time.Second)
		}
	}()
}

func tryKafkaConnection(config KafkaConnectionConfig, handler KafkaHandler) {
	consumer, err := startKafkaConsumer(config.CConfigMap, config.CTopics)
	defer consumer.Close()
	if err != nil {
		return
	}
	producer, err := startKafkaProducer(config.PConfigMap)
	defer producer.Close()
	if err != nil {
		return
	}
	go handleEvents(producer)
	//this function stays running while the connection is up
	waitForPacketsOnKafka(consumer, producer, handler)

}

func startKafkaConsumer(configMap *kafka.ConfigMap, topics []string) (*kafka.Consumer, error) {

	consumer, err := kafka.NewConsumer(configMap)
	if err != nil {
		log.Error("Failed to create Kafka consumer: ", err)
		return consumer, err
	}

	val, _ := configMap.Get("bootstrap.servers", nil)
	log.Info("Connecting to kafka on ", val.(string))

	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		log.Error("Failed to subscribe to consumer topic: ", err)
		return consumer, err
	}

	return consumer, nil
}

func startKafkaProducer(configMap *kafka.ConfigMap) (*kafka.Producer, error) {

	producer, err := kafka.NewProducer(configMap)
	if err != nil {
		log.Error("Failed to create Kafka producer: ", err)
		producer.Close()
		return nil, err
	}
	return producer, nil
}

func handleEvents(producer *kafka.Producer) {
	log.Info("Starting event handling thread")
	for e := range producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Error("Failed to deliver message: ", ev.TopicPartition, string(ev.Value))
			} else {
				log.Info("Messaged sucessfuly delivered to kafka. offset=", ev.TopicPartition.Offset)
			}
		}
	}
	log.Info("Stopping event handling thread")
}

func waitForPacketsOnKafka(consumer *kafka.Consumer, producer *kafka.Producer, handler KafkaHandler) {
	for {
		ev := consumer.Poll(5000)
		switch e := ev.(type) {
		case *kafka.Message:
			packet := &pb.Msg{}
			if err := proto.Unmarshal(e.Value, packet); err != nil {
				log.Error("Could not deserialize the received packet")
				continue
			}
			log.Info("Received packet:\n", packet.String())
			handleKafkaPacket(packet, producer, handler)

		case kafka.PartitionEOF:
			log.Error("Reached the end of the partition", e)
		case kafka.Error:
			log.Error("Kafka Error: ", e)
			return
		default:
			//no other job to do. Keep waiting
			continue
		}
	}
}
