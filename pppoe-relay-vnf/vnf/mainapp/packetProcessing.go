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
	"encoding/binary"
	"fmt"

	dbDao "github.com/BroadbandForum/obbaa-477/common/db/dao"
	pb "github.com/BroadbandForum/obbaa-477/common/pb/tr477"
	log "github.com/BroadbandForum/obbaa-477/common/utils/log"
	"github.com/BroadbandForum/obbaa-477/pppoe-relay-vnf/vnf/pppoedb"
	"github.com/BroadbandForum/obbaa-477/pppoe-relay-vnf/vnf/pppoedb/device/subscriber_profile"
	"github.com/BroadbandForum/obbaa-477/pppoe-relay-vnf/vnf/pppoedb/device/vsi"
	"github.com/BroadbandForum/obbaa-477/pppoe-relay-vnf/vnf/pppoedb/pppoe_profile"
	pppoePacket "github.com/BroadbandForum/obbaa-477/pppoe-relay-vnf/vnf/pppoepacket"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"golang.org/x/exp/slices"
)

func processPacket(in *pb.CpriMsg) (*pb.CpriMsg, error) {
	packet := gopacket.NewPacket(in.Packet, layers.LayerTypeEthernet, gopacket.NoCopy)
	var oltPacket *pb.CpriMsg
	//if packet is a PPPoE Discovery packet, insert circuitId and remoteId tags
	if packetPPPoE := packet.Layer(pppoePacket.LayerTypePPPoEDiscovery); packetPPPoE != nil {
		pppoePacketDiscovery := packetPPPoE.(*pppoePacket.PPPoEDiscovery)
		//TODO: insert circuitId and remoteId according to configuration
		cidValue, ridValue, err := getIds(in)
		if err != nil {
			log.Error("Error getting circuit-id and remote-id: ", err)
			return in, err
		}
		deviceName := in.MetaData.Generic.DeviceName
		vsiName := in.MetaData.Generic.DeviceInterface
		log.Info(fmt.Sprintf("Found circuit id '%s' and remote id '%s' for device/intf (%s/%s)", cidValue, ridValue, deviceName, vsiName))
		bbfTagValue := pppoePacket.BBFTagValue{
			VendorId: pppoePacket.BBF_VENDOR_ID,
			SubTags: []pppoePacket.BBFSubTag{
				{
					Number: pppoePacket.BBF_CIRCUIT_ID_SUBTAG_NUMBER,
					Length: uint8(binary.Size(cidValue)),
					Value:  cidValue,
				},
				{
					Number: pppoePacket.BBF_REMOTE_ID_SUBTAG_NUMBER,
					Length: uint8(binary.Size(ridValue)),
					Value:  ridValue,
				},
			},
		}
		pppoePacketDiscovery.Tags = append(pppoePacketDiscovery.Tags, pppoePacket.PPPoETag{Type: VENDOR_SPECIFIC_TAG_TYPE,
			Value: bbfTagValue.Serialize()})
		buffer := gopacket.NewSerializeBuffer()
		options := gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		}
		err = gopacket.SerializePacket(buffer, options, packet)
		if err != nil {
			log.Error("Error serializing packet: ", err)
			return in, err
		}
		oltPacket = &pb.CpriMsg{
			Header:   in.Header,
			MetaData: in.MetaData,
			Packet:   buffer.Bytes(),
		}
	} else {
		//if packet is a PPPoE Session packet, simply send it back to olt as is
		oltPacket = in
	}
	return oltPacket, nil
}

func getIds(msg *pb.CpriMsg) ([]byte, []byte, error) {
	database := mongoClient.Database(pppoedb.DatabaseName)
	cidValue := []byte{}
	ridValue := []byte{}
	vsiFilter := vsi.Vsi{
		DeviceName: msg.MetaData.Generic.DeviceName,
		VsiName:    msg.MetaData.Generic.DeviceInterface,
	}
	vsi, err := dbDao.GetVNFDocument(&vsiFilter, database, "vsi-list")
	if err != nil {
		log.Error(fmt.Sprintf("Unable to find vsi %+v: %s", vsiFilter, err))
		return nil, nil, err
	}
	if sProfileName := vsi.SubscriberProfile; sProfileName != "" {
		sProfileFilter := subscriber_profile.SubscriberProfile{
			Name:       sProfileName,
			DeviceName: msg.MetaData.Generic.DeviceName,
		}
		sProfile, err := dbDao.GetVNFDocument(&sProfileFilter, database, "subscriber-profiles")
		if err != nil {
			log.Error(fmt.Sprintf("Unable to find subscriber profile %+v: %s", sProfileFilter, err))
			return nil, nil, err
		}
		cidValue = []byte(sProfile.CircuitID)
		ridValue = []byte(sProfile.RemoteID)
	} else {
		pProfileFilter := pppoe_profile.PPPoEProfile{Name: vsi.PppoeProfile}
		pProfile, err := dbDao.GetVNFDocument(&pProfileFilter, database, "pppoe-profiles")
		if err != nil {
			log.Error(fmt.Sprintf("Unable to find pppoe profile %v: %s", pProfileFilter, err))
			return nil, nil, err
		}
		cidValue, ridValue, err = getIdsFromSyntax(pProfile, msg)
		if err != nil {
			return nil, nil, err
		}
	}
	return cidValue, ridValue, nil
}

func getIdsFromSyntax(pProfile *pppoe_profile.PPPoEProfile, msg *pb.CpriMsg) ([]byte, []byte, error) {
	var cidValue []byte = nil
	var ridValue []byte = nil
	var err error
	if slices.Contains(pProfile.PppoeVendorSpecificTag.Subtag, CID_SUBTAG) {
		cidValueSyntax, err := interpretSyntax(pProfile.PppoeVendorSpecificTag.DefaultCircuitIDSyntax, msg)
		if err != nil {
			return nil, nil, err
		}
		cidValue = []byte(cidValueSyntax)
	}
	if slices.Contains(pProfile.PppoeVendorSpecificTag.Subtag, RID_SUBTAG) {
		ridValueSyntax, err := interpretSyntax(pProfile.PppoeVendorSpecificTag.DefaultRemoteIDSyntax, msg)
		if err != nil {
			return nil, nil, err
		}
		ridValue = []byte(ridValueSyntax)
	}
	return cidValue, ridValue, err
}
