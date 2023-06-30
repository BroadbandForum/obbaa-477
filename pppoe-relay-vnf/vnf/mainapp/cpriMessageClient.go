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
* Created by João Correia(Altice Labs) on 03/04/2023
 */
package main

import (
	"context"
	"os"

	"github.com/obbaa-477/common/pb/tr477"
	log "github.com/obbaa-477/common/utils/log"
	"google.golang.org/grpc"
)

func runAsClient() {
	cc, err := grpc.Dial(socketGrpc, grpc.WithInsecure())
	if err != nil {
		return
	}
	defer cc.Close()

	packetClient := tr477.NewCpriMessageClient(cc)
	helloClient := tr477.NewCpriHelloClient(cc)
	sendHello(helloClient)
	waitForPacketsOnStream(packetClient)

}

func sendHello(client tr477.CpriHelloClient) {

	req := &tr477.HelloCpriRequest{
		LocalEndpointHello: &tr477.Hello{
			EntityName:   os.Getenv("VNF_NAME"),
			EndpointName: os.Getenv("VNF_NAME"),
		},
	}
	_, err := client.HelloCpri(context.Background(), req)
	if err != nil {
		log.Error("- error while calling Hello RPC: ", err)
	}
}

func waitForPacketsOnStream(client tr477.CpriMessageClient) {
	log.Info("Waiting for Packets")
	stream, err := client.TransferCpri(context.Background())
	if err != nil {
		log.Error("TransferCpri failed: ", err)
		return
	}
	for {
		packet, err := stream.Recv()
		if err != nil {
			log.Error("Error receiving packet: ", err)
			break
		}
		if packet != nil {
			log.Info("Received packet: ", packet)
			oltPacket, err := processPacket(packet)
			if err != nil {
				if discard_on_error {
					log.Warning("Discarding packet")
					continue
				}
			}
			err = stream.Send(oltPacket)
			if err != nil {
				log.Error("Error sending packet: ", err)
			}
		}
	}
}
