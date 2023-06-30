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
* BBF PPPoE tag encoder/decoder
*
* Created by João Correia(Altice Labs) on 17/02/2023
 */
package pppoepacket

import (
	"encoding/binary"
)

//tag format described in TR-101i2 Appendix A

const (
	BBF_VENDOR_ID                = 0x00000DE9
	BBF_CIRCUIT_ID_SUBTAG_NUMBER = 0x01 //TR-101i2 section 8.3
	BBF_REMOTE_ID_SUBTAG_NUMBER  = 0x02 //TR-101i2 section 8.3
)

type BBFSubTag struct {
	Number uint8
	Length uint8
	Value  []byte
}

type BBFTagValue struct {
	VendorId uint32
	SubTags  []BBFSubTag
}

func (p *BBFTagValue) Serialize() []byte {
	bytes := make([]byte, binary.Size(uint32(0)))
	binary.BigEndian.PutUint32(bytes, p.VendorId)
	for _, subtag := range p.SubTags {
		bytes = append(bytes, subtag.Serialize()...)
	}
	return bytes
}

func (p *BBFSubTag) Serialize() []byte {
	bytes := make([]byte, 2*binary.Size(uint8(0))+len(p.Value))
	bytes[0] = p.Number
	bytes[1] = p.Length
	copy(bytes[2:], p.Value)
	return bytes
}
