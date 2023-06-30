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
	"errors"
	"fmt"
	"strings"

	"github.com/obbaa-477/common/db/dao"
	pb "github.com/obbaa-477/common/pb/tr477"
	log "github.com/obbaa-477/common/utils/log"
	"github.com/obbaa-477/pppoe-relay-vnf/vnf/pppoedb"
	"github.com/obbaa-477/pppoe-relay-vnf/vnf/pppoedb/device"
)

var syntaxMap = map[string]func([]string, int, *pb.CpriMsg) error{
	"access-node-identifier": accessNodeIdentifier,
	"slot":                   slot,
	"sub-slot":               subSlot,
	"chassis":                chassis,
	"rack":                   rack,
	"onuid":                  onuid,
	"port":                   port,
}

const interpretSyntaxError = "Error interpreting syntax token \"%s\""
const separators = "#.,;/ "

func interpretSyntax(syntax string, msg *pb.CpriMsg) (string, error) {
	tokens := strings.FieldsFunc(syntax, func(r rune) bool {
		return strings.ContainsRune(separators, r)
	})
	for i := 0; i < len((tokens)); i++ {
		token := strings.ToLower(tokens[i])
		if f, keyExists := syntaxMap[token]; keyExists {
			err := f(tokens, i, msg)
			if err != nil {
				return "", err
			}
		}
	}
	return strings.Join(tokens, " "), nil
}

func accessNodeIdentifier(tokens []string, pos int, msg *pb.CpriMsg) error {
	log.Info("accessNode: " + msg.String())
	if metadata := msg.GetMetaData(); metadata != nil {
		if genericMetadata := metadata.GetGeneric(); genericMetadata != nil {
			deviceFilter := device.Device{Name: genericMetadata.DeviceName}
			device, err := dao.GetVNFDocument(&deviceFilter, mongoClient.Database(pppoedb.DatabaseName), "devices-using-d-olt-pppoeia")
			if err != nil {
				return err
			}
			tokens[pos] = device.AccessNodeId
			return nil
		}
	}
	return errors.New(fmt.Sprintf(interpretSyntaxError, tokens[pos]))
}

func slot(tokens []string, pos int, msg *pb.CpriMsg) error {
	log.Info("slot: " + msg.String())
	if metadata := msg.GetMetaData(); metadata != nil {
		if pppoeMetadata := metadata.GetPppoe(); pppoeMetadata != nil {
			if hwIdElements := pppoeMetadata.GetHwIdElements(); hwIdElements != nil {
				tokens[pos] = hwIdElements.SlotNumber
				return nil
			}
		}
	}
	return errors.New(fmt.Sprintf(interpretSyntaxError, tokens[pos]))
}

func port(tokens []string, pos int, msg *pb.CpriMsg) error {
	if metadata := msg.GetMetaData(); metadata != nil {
		if pppoeMetadata := metadata.GetPppoe(); pppoeMetadata != nil {
			if hwIdElements := pppoeMetadata.GetHwIdElements(); hwIdElements != nil {
				tokens[pos] = hwIdElements.PortNumber
				return nil
			}
		}
	}
	return errors.New(fmt.Sprintf(interpretSyntaxError, tokens[pos]))
}

func onuid(tokens []string, pos int, msg *pb.CpriMsg) error {
	if metadata := msg.GetMetaData(); metadata != nil {
		if onuMetadata := metadata.GetOnu(); onuMetadata != nil {
			tokens[pos] = onuMetadata.OnuId
			return nil
		}
	}
	return errors.New(fmt.Sprintf(interpretSyntaxError, tokens[pos]))
}

func chassis(tokens []string, pos int, msg *pb.CpriMsg) error {
	if metadata := msg.GetMetaData(); metadata != nil {
		if pppoeMetadata := metadata.GetPppoe(); pppoeMetadata != nil {
			if hwIdElements := pppoeMetadata.GetHwIdElements(); hwIdElements != nil {
				tokens[pos] = hwIdElements.ChassisNumber
				return nil
			}
		}
	}
	return errors.New(fmt.Sprintf(interpretSyntaxError, tokens[pos]))
}

func rack(tokens []string, pos int, msg *pb.CpriMsg) error {
	if metadata := msg.GetMetaData(); metadata != nil {
		if pppoeMetadata := metadata.GetPppoe(); pppoeMetadata != nil {
			if hwIdElements := pppoeMetadata.GetHwIdElements(); hwIdElements != nil {
				tokens[pos] = hwIdElements.RackNumber
				return nil
			}
		}
	}
	return errors.New(fmt.Sprintf(interpretSyntaxError, tokens[pos]))
}

func frame(tokens []string, pos int, msg *pb.CpriMsg) error {
	if metadata := msg.GetMetaData(); metadata != nil {
		if pppoeMetadata := metadata.GetPppoe(); pppoeMetadata != nil {
			if hwIdElements := pppoeMetadata.GetHwIdElements(); hwIdElements != nil {
				tokens[pos] = hwIdElements.FrameNumber
				return nil
			}
		}
	}
	return errors.New(fmt.Sprintf(interpretSyntaxError, tokens[pos]))
}

func subSlot(tokens []string, pos int, msg *pb.CpriMsg) error {
	if metadata := msg.GetMetaData(); metadata != nil {
		if pppoeMetadata := metadata.GetPppoe(); pppoeMetadata != nil {
			if hwIdElements := pppoeMetadata.GetHwIdElements(); hwIdElements != nil {
				tokens[pos] = hwIdElements.SubSlotNumber
				return nil
			}
		}
	}
	return errors.New(fmt.Sprintf(interpretSyntaxError, tokens[pos]))
}
