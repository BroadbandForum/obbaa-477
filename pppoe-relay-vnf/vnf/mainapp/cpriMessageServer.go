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
	"io"
	"net"

	tr477 "github.com/BroadbandForum/obbaa-477/common/pb/tr477"
	log "github.com/BroadbandForum/obbaa-477/common/utils/log"
	"google.golang.org/grpc"
)

const VENDOR_SPECIFIC_TAG_TYPE = 0x0105

const (
	CID_SUBTAG = "circuit-id"
	RID_SUBTAG = "remote-id"
)

type cpriMessageServer struct {
	tr477.UnimplementedCpriMessageServer
	Quit chan bool
}

func runAsServer() {
	lis, err := net.Listen("tcp", socketGrpc)
	if err != nil {
		log.Fatal("VNF - Failed to listen: ", err)
	}
	s := grpc.NewServer()

	addHelloServiceServer := cpriHelloService{}
	addPacketServiceServer := cpriMessageServer{Quit: make(chan bool)}
	tr477.RegisterCpriHelloServer(s, &addHelloServiceServer)
	tr477.RegisterCpriMessageServer(s, &addPacketServiceServer)

	if err := s.Serve(lis); err != nil {
		log.Fatal("VNF - Failed to serve: %v", err)
	}
}

func (c *cpriMessageServer) TransferCpri(stream tr477.CpriMessage_TransferCpriServer) error {
	for {
		select {
		case <-c.Quit:
			return nil
		default:
			in, err := stream.Recv()
			if err != nil {
				log.Error("Error receiving packet: ", err)
				if err == io.EOF {
					return nil
				}
				continue
			}

			oltPacket, err := processPacket(in)
			if err != nil {
				if discard_on_error {
					log.Warning("Discarding packet")
					continue
				}
			}
			if err := stream.Send(oltPacket); err != nil {
				log.Error("Failed to send packet: ", err)
			}
		}
	}
}
