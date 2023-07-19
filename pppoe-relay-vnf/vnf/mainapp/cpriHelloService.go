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
	"context"
	"os"

	pb "github.com/BroadbandForum/obbaa-477/common/pb/tr477"
	log "github.com/BroadbandForum/obbaa-477/common/utils/log"
)

type cpriHelloService struct {
	pb.UnimplementedCpriHelloServer
}

func (c *cpriHelloService) HelloCpri(ctx context.Context, in *pb.HelloCpriRequest) (*pb.HelloCpriResponse, error) {

	if in == nil || ctx == nil {
		log.Error("VNF - Device without information.")
		return nil, nil
	}

	deviceName := in.LocalEndpointHello.EntityName
	if deviceName == "" {
		log.Error("VNF - Device without name.")
		return nil, nil
	}
	log.Info("Received hello from ", deviceName)
	res := &pb.HelloCpriResponse{
		RemoteEndpointHello: &pb.Hello{
			EntityName:   os.Getenv("VNF_NAME"),
			EndpointName: os.Getenv("VNF_NAME"),
		},
	}
	return res, nil
}
