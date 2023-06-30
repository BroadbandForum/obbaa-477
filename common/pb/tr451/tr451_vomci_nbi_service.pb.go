// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: tr451_vomci_nbi_service.proto

//**************************************************************************
// TR-451 vOMCI NBI Service Protocol Buffer Schema
//
//  Copyright (c) 2021, Broadband Forum
//
//  Redistribution and use in source and binary forms, with or
//   without modification, are permitted provided that the following
//   conditions are met:
//
//   1. Redistributions of source code must retain the above copyright
//      notice, this list of conditions and the following disclaimer.
//
//   2. Redistributions in binary form must reproduce the above
//      copyright notice, this list of conditions and the following
//      disclaimer in the documentation and/or other materials
//      provided with the distribution.
//
//   3. Neither the name of the copyright holder nor the names of its
//      contributors may be used to endorse or promote products
//      derived from this software without specific prior written
//      permission.
//
//   THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND
//   CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
//   INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
//   MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
//   DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR
//   CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
//   SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
//   NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
//   LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
//   CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
//   STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
//   ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
//   ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
//   The above license is used as a license under copyright only.
//   Please reference the Forum IPR Policy for patent licensing terms
//   <https://www.broadband-forum.org/ipr-policy>.
//
//   Any moral rights which are necessary to exercise under the above
//   license grant are also deemed granted under this license.
//
// | Version           | Name                   | Date       |
// | TR-451 1.0.0      | vOMCI Specification    | TBD, 2021  |
//
// BBF software release registry: http://www.broadband-forum.org/software
//**************************************************************************

package tr451

import (
	empty "github.com/golang/protobuf/ptypes/empty"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_tr451_vomci_nbi_service_proto protoreflect.FileDescriptor

var file_tr451_vomci_nbi_service_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x74, 0x72, 0x34, 0x35, 0x31, 0x5f, 0x76, 0x6f, 0x6d, 0x63, 0x69, 0x5f, 0x6e, 0x62,
	0x69, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x1a, 0x74, 0x72, 0x34, 0x35, 0x31, 0x5f, 0x76, 0x6f, 0x6d, 0x63, 0x69, 0x5f, 0x6e, 0x62, 0x69,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1d, 0x74, 0x72, 0x34,
	0x35, 0x31, 0x5f, 0x76, 0x6f, 0x6d, 0x63, 0x69, 0x5f, 0x6e, 0x62, 0x69, 0x5f, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x98, 0x05, 0x0a, 0x0f, 0x56, 0x6f, 0x6d, 0x63,
	0x69, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4e, 0x62, 0x69, 0x12, 0x49, 0x0a, 0x05, 0x48,
	0x65, 0x6c, 0x6c, 0x6f, 0x12, 0x1f, 0x2e, 0x74, 0x72, 0x34, 0x35, 0x31, 0x5f, 0x76, 0x6f, 0x6d,
	0x63, 0x69, 0x5f, 0x6e, 0x62, 0x69, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x4d, 0x73, 0x67, 0x1a, 0x1f, 0x2e, 0x74, 0x72, 0x34, 0x35, 0x31, 0x5f, 0x76, 0x6f,
	0x6d, 0x63, 0x69, 0x5f, 0x6e, 0x62, 0x69, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x4d, 0x73, 0x67, 0x12, 0x4b, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x44, 0x61, 0x74,
	0x61, 0x12, 0x1f, 0x2e, 0x74, 0x72, 0x34, 0x35, 0x31, 0x5f, 0x76, 0x6f, 0x6d, 0x63, 0x69, 0x5f,
	0x6e, 0x62, 0x69, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d,
	0x73, 0x67, 0x1a, 0x1f, 0x2e, 0x74, 0x72, 0x34, 0x35, 0x31, 0x5f, 0x76, 0x6f, 0x6d, 0x63, 0x69,
	0x5f, 0x6e, 0x62, 0x69, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x4d, 0x73, 0x67, 0x12, 0x51, 0x0a, 0x0d, 0x52, 0x65, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x12, 0x1f, 0x2e, 0x74, 0x72, 0x34, 0x35, 0x31, 0x5f, 0x76, 0x6f, 0x6d,
	0x63, 0x69, 0x5f, 0x6e, 0x62, 0x69, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x4d, 0x73, 0x67, 0x1a, 0x1f, 0x2e, 0x74, 0x72, 0x34, 0x35, 0x31, 0x5f, 0x76, 0x6f,
	0x6d, 0x63, 0x69, 0x5f, 0x6e, 0x62, 0x69, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x4d, 0x73, 0x67, 0x12, 0x57, 0x0a, 0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x12, 0x1f, 0x2e,
	0x74, 0x72, 0x34, 0x35, 0x31, 0x5f, 0x76, 0x6f, 0x6d, 0x63, 0x69, 0x5f, 0x6e, 0x62, 0x69, 0x5f,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x73, 0x67, 0x1a, 0x1f,
	0x2e, 0x74, 0x72, 0x34, 0x35, 0x31, 0x5f, 0x76, 0x6f, 0x6d, 0x63, 0x69, 0x5f, 0x6e, 0x62, 0x69,
	0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x73, 0x67, 0x12,
	0x58, 0x0a, 0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x49,
	0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x1f, 0x2e, 0x74, 0x72, 0x34, 0x35, 0x31, 0x5f,
	0x76, 0x6f, 0x6d, 0x63, 0x69, 0x5f, 0x6e, 0x62, 0x69, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x73, 0x67, 0x1a, 0x1f, 0x2e, 0x74, 0x72, 0x34, 0x35, 0x31,
	0x5f, 0x76, 0x6f, 0x6d, 0x63, 0x69, 0x5f, 0x6e, 0x62, 0x69, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x73, 0x67, 0x12, 0x47, 0x0a, 0x03, 0x52, 0x50, 0x43,
	0x12, 0x1f, 0x2e, 0x74, 0x72, 0x34, 0x35, 0x31, 0x5f, 0x76, 0x6f, 0x6d, 0x63, 0x69, 0x5f, 0x6e,
	0x62, 0x69, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x73,
	0x67, 0x1a, 0x1f, 0x2e, 0x74, 0x72, 0x34, 0x35, 0x31, 0x5f, 0x76, 0x6f, 0x6d, 0x63, 0x69, 0x5f,
	0x6e, 0x62, 0x69, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d,
	0x73, 0x67, 0x12, 0x4a, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x2e, 0x74,
	0x72, 0x34, 0x35, 0x31, 0x5f, 0x76, 0x6f, 0x6d, 0x63, 0x69, 0x5f, 0x6e, 0x62, 0x69, 0x5f, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x73, 0x67, 0x1a, 0x1f, 0x2e,
	0x74, 0x72, 0x34, 0x35, 0x31, 0x5f, 0x76, 0x6f, 0x6d, 0x63, 0x69, 0x5f, 0x6e, 0x62, 0x69, 0x5f,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x73, 0x67, 0x12, 0x52,
	0x0a, 0x15, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x46, 0x6f, 0x72, 0x4e, 0x6f, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a,
	0x1f, 0x2e, 0x74, 0x72, 0x34, 0x35, 0x31, 0x5f, 0x76, 0x6f, 0x6d, 0x63, 0x69, 0x5f, 0x6e, 0x62,
	0x69, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x73, 0x67,
	0x30, 0x01, 0x42, 0x11, 0x5a, 0x0f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x70, 0x62, 0x2f,
	0x74, 0x72, 0x34, 0x35, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_tr451_vomci_nbi_service_proto_goTypes = []interface{}{
	(*Msg)(nil),         // 0: tr451_vomci_nbi_message.v1.Msg
	(*empty.Empty)(nil), // 1: google.protobuf.Empty
}
var file_tr451_vomci_nbi_service_proto_depIdxs = []int32{
	0, // 0: tr451_vomci_nbi_service.v1.VomciMessageNbi.Hello:input_type -> tr451_vomci_nbi_message.v1.Msg
	0, // 1: tr451_vomci_nbi_service.v1.VomciMessageNbi.GetData:input_type -> tr451_vomci_nbi_message.v1.Msg
	0, // 2: tr451_vomci_nbi_service.v1.VomciMessageNbi.ReplaceConfig:input_type -> tr451_vomci_nbi_message.v1.Msg
	0, // 3: tr451_vomci_nbi_service.v1.VomciMessageNbi.UpdateConfigReplica:input_type -> tr451_vomci_nbi_message.v1.Msg
	0, // 4: tr451_vomci_nbi_service.v1.VomciMessageNbi.UpdateConfigInstance:input_type -> tr451_vomci_nbi_message.v1.Msg
	0, // 5: tr451_vomci_nbi_service.v1.VomciMessageNbi.RPC:input_type -> tr451_vomci_nbi_message.v1.Msg
	0, // 6: tr451_vomci_nbi_service.v1.VomciMessageNbi.Action:input_type -> tr451_vomci_nbi_message.v1.Msg
	1, // 7: tr451_vomci_nbi_service.v1.VomciMessageNbi.ListenForNotification:input_type -> google.protobuf.Empty
	0, // 8: tr451_vomci_nbi_service.v1.VomciMessageNbi.Hello:output_type -> tr451_vomci_nbi_message.v1.Msg
	0, // 9: tr451_vomci_nbi_service.v1.VomciMessageNbi.GetData:output_type -> tr451_vomci_nbi_message.v1.Msg
	0, // 10: tr451_vomci_nbi_service.v1.VomciMessageNbi.ReplaceConfig:output_type -> tr451_vomci_nbi_message.v1.Msg
	0, // 11: tr451_vomci_nbi_service.v1.VomciMessageNbi.UpdateConfigReplica:output_type -> tr451_vomci_nbi_message.v1.Msg
	0, // 12: tr451_vomci_nbi_service.v1.VomciMessageNbi.UpdateConfigInstance:output_type -> tr451_vomci_nbi_message.v1.Msg
	0, // 13: tr451_vomci_nbi_service.v1.VomciMessageNbi.RPC:output_type -> tr451_vomci_nbi_message.v1.Msg
	0, // 14: tr451_vomci_nbi_service.v1.VomciMessageNbi.Action:output_type -> tr451_vomci_nbi_message.v1.Msg
	0, // 15: tr451_vomci_nbi_service.v1.VomciMessageNbi.ListenForNotification:output_type -> tr451_vomci_nbi_message.v1.Msg
	8, // [8:16] is the sub-list for method output_type
	0, // [0:8] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_tr451_vomci_nbi_service_proto_init() }
func file_tr451_vomci_nbi_service_proto_init() {
	if File_tr451_vomci_nbi_service_proto != nil {
		return
	}
	file_tr451_vomci_nbi_message_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_tr451_vomci_nbi_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tr451_vomci_nbi_service_proto_goTypes,
		DependencyIndexes: file_tr451_vomci_nbi_service_proto_depIdxs,
	}.Build()
	File_tr451_vomci_nbi_service_proto = out.File
	file_tr451_vomci_nbi_service_proto_rawDesc = nil
	file_tr451_vomci_nbi_service_proto_goTypes = nil
	file_tr451_vomci_nbi_service_proto_depIdxs = nil
}
