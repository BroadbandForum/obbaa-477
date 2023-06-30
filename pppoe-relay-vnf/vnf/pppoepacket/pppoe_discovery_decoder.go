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
* PPPOE VNF decoder
*
* Created by João Correia(Altice Labs) on 17/02/2023
 */
package pppoepacket

import (
	"encoding/binary"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type PPPoETag struct {
	Type  uint16
	Value []byte
}

type PPPoEDiscovery struct {
	layers.BaseLayer
	Tags []PPPoETag
}

var LayerTypePPPoEDiscovery = gopacket.RegisterLayerType(1000, gopacket.LayerTypeMetadata{Name: "PPPoE Discovery", Decoder: gopacket.DecodeFunc(decodePPPoEDiscovery)})

// LayerType returns gopacket.LayerTypePPPoE.
func (p *PPPoEDiscovery) LayerType() gopacket.LayerType {
	return LayerTypePPPoEDiscovery
}

// decodePPPoE decodes the PPPoE header (see http://tools.ietf.org/html/rfc2516).
func decodePPPoEDiscovery(data []byte, p gopacket.PacketBuilder) error {
	pppoeDiscovery := &PPPoEDiscovery{}
	byteCount := 0 //count the bytes decoded to determine where in the packet the next tag is
	for byteCount < len(data) {
		tagType := binary.BigEndian.Uint16(data[byteCount:])
		tagLength := binary.BigEndian.Uint16(data[byteCount+binary.Size(uint16(0)):])
		pppoeTag := PPPoETag{Type: tagType,
			Value: data[byteCount+2*binary.Size(uint16(0)) : byteCount+2*binary.Size(uint16(0))+int(tagLength)]}
		pppoeDiscovery.Tags = append(pppoeDiscovery.Tags, pppoeTag)
		//at this point, the type(2 bytes), length(2 bytes) and value (length bytes) of a tag were decoded
		byteCount += int(tagLength) + 2*binary.Size(uint16(0))
	}
	pppoeDiscovery.BaseLayer = layers.BaseLayer{Contents: data, Payload: []byte{}}
	p.AddLayer(pppoeDiscovery)
	return nil
}

func (p *PPPoEDiscovery) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	for _, tag := range p.Tags {
		//append a number of bytes equal to the length of the value, plus the size of type and length (2 bytes each) of a PPPoE Tag
		//(appending bytes instead of prepending to keep the order of the tags)
		bytes, err := b.AppendBytes(len(tag.Value) + 2*binary.Size(uint16(0)))
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint16(bytes[0:], tag.Type)
		binary.BigEndian.PutUint16(bytes[2:], uint16(len(tag.Value)))
		copy(bytes[4:], tag.Value)
	}
	return nil
}

func InitPPPoEDecoder() {
	layers.PPPoECodeMetadata[layers.PPPoECodePADI] = layers.EnumMetadata{DecodeWith: gopacket.DecodeFunc(decodePPPoEDiscovery), Name: "PADI"}
	layers.PPPoECodeMetadata[layers.PPPoECodePADO] = layers.EnumMetadata{DecodeWith: gopacket.DecodeFunc(decodePPPoEDiscovery), Name: "PADO"}
	layers.PPPoECodeMetadata[layers.PPPoECodePADR] = layers.EnumMetadata{DecodeWith: gopacket.DecodeFunc(decodePPPoEDiscovery), Name: "PADR"}
	layers.PPPoECodeMetadata[layers.PPPoECodePADS] = layers.EnumMetadata{DecodeWith: gopacket.DecodeFunc(decodePPPoEDiscovery), Name: "PADS"}
}
