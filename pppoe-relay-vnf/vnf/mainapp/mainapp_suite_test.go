package main

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/BroadbandForum/obbaa-477/common/pb/tr477"
	pppoedb "github.com/BroadbandForum/obbaa-477/pppoe-relay-vnf/vnf/pppoedb"
	"github.com/BroadbandForum/obbaa-477/pppoe-relay-vnf/vnf/pppoepacket"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"google.golang.org/grpc"
)

var mongoTest *mtest.T

func TestMainapp(t *testing.T) {
	RegisterFailHandler(Fail)
	mongoTest = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	RunSpecs(t, "Mainapp Suite")
}

var _ = Describe("Vnf", Ordered, func() {
	Describe("as Server", func() {
		BeforeAll(func() {
			os.Setenv("VNF_NAME", "PPPOE_HELLO_ENTITY_NAME")
		})
		Describe("Hello Service", func() {
			var testHelloService cpriHelloService

			BeforeEach(func() {
				testHelloService = cpriHelloService{}
				tr477.RegisterCpriHelloServer(grpc.NewServer(), &testHelloService)
			})

			Context("with nil request", func() {
				It("Should return a nil response", func() {
					resp, _ := testHelloService.HelloCpri(context.Background(), nil)
					Expect(resp).To(BeNil())
				})
			})

			Context("with empty Entity Name", func() {
				It("should return a nil response", func() {
					resp, _ := testHelloService.HelloCpri(context.Background(), &tr477.HelloCpriRequest{
						LocalEndpointHello: &tr477.Hello{
							EntityName:   "",
							EndpointName: "",
						},
					})
					Expect(resp).To(BeNil())
				})
			})

			Context("With valid ctx and valid request", func() {
				It("Should return a valid response", func() {
					resp, err := testHelloService.HelloCpri(context.Background(), &tr477.HelloCpriRequest{
						LocalEndpointHello: &tr477.Hello{
							EntityName:   "DEVICE_HELLO_ENTITY_NAME",
							EndpointName: "DEVICE_HELLO_ENDPOINT_NAME",
						},
					})
					Expect(err).To(BeNil())
					Expect(resp.RemoteEndpointHello.EntityName).To(BeEquivalentTo("PPPOE_HELLO_ENTITY_NAME"))
					Expect(resp.RemoteEndpointHello.EndpointName).To(BeEquivalentTo("PPPOE_HELLO_ENTITY_NAME"))
				})
			})
		})

		Describe("Message Service", func() {
			var (
				testMessageService cpriMessageServer
				testMsg            *tr477.CpriMsg
				// timeout            = 1 * time.Second
				stream *MockTransferCpriServer
			)
			BeforeAll(func() {
				pppoepacket.InitPPPoEDecoder()

			})

			BeforeEach(func() {
				stream = &MockTransferCpriServer{}
				testMessageService = cpriMessageServer{Quit: make(chan bool)}
				tr477.RegisterCpriMessageServer(grpc.NewServer(), &testMessageService)
			})

			Context("with 'Recv' returning an EOF error", func() {
				It("should not call 'Send'", func() {
					stream.On("Recv").Return(testMsg, io.EOF)
					testMessageService.TransferCpri(stream)

					stream.AssertCalled(GinkgoT(), "Recv")
					stream.AssertNotCalled(GinkgoT(), "Send")
				})
			})

			Context("without database configurations", func() {
				Context("with a non pppoe discovery packet", func() {
					BeforeEach(func() {
						testMsg = &tr477.CpriMsg{Packet: []byte{1, 2, 3}}
					})
					It("should call 'Send' with the received message as argument", func() {

						stream.On("Recv").Return(testMsg, nil)
						stream.On("Send", testMsg).Return(nil)
						go func() {
							for len(stream.Calls) < 2 {
							}
							testMessageService.Quit <- true
						}()
						testMessageService.TransferCpri(stream)

						stream.AssertExpectations(GinkgoT())
					})
				})
				Context("with a pppoe discovery packet", func() {
					var (
						padiRequest gopacket.Packet
					)
					Context("discarding on error", func() {
						BeforeAll(func() {
							discard_on_error = true
						})
						BeforeEach(func() {
							padiRequest, _ = getTestPADIPackets()
							testMsg = &tr477.CpriMsg{
								MetaData: &tr477.CpriMetaData{
									Generic: &tr477.GenericMetadata{
										DeviceName:      "TEST_DEVICE_NAME",
										DeviceInterface: "TEST_VSI_NAME",
									},
								},
								Packet: padiRequest.Data(),
							}
						})
						It("should discard the packet", func() {
							mongoTest.Run("discarding", func(mt *mtest.T) {
								mongoClient = mt.Client
								stream.On("Recv").Return(testMsg, nil)
								stream.On("Send", testMsg).Return(nil)

								go func() {
									for len(stream.Calls) < 2 {
									}
									testMessageService.Quit <- true
								}()
								testMessageService.TransferCpri(stream)

								stream.AssertCalled(mt.T, "Recv")
							})

						})
					})
					Context("not discarding on error", func() {
						BeforeAll(func() {
							discard_on_error = false
						})
						It("should call 'Send' with the received message", func() {
							stream.On("Recv").Return(testMsg, nil)
							stream.On("Send", testMsg).Return(nil)
							go func() {
								for len(stream.Calls) < 2 {
								}
								testMessageService.Quit <- true
							}()
							testMessageService.TransferCpri(stream)
							stream.AssertCalled(GinkgoT(), "Recv")
							stream.AssertCalled(GinkgoT(), "Send", testMsg)
						})
					})
				})
			})
			Context("with database configurations", func() {
				var (
					padiRequest           gopacket.Packet
					padiResponse          gopacket.Packet
					testVsi               bson.D
					testSubscriberProfile bson.D
					expectedMsg           *tr477.CpriMsg
				)
				BeforeAll(func() {
					padiRequest, padiResponse = getTestPADIPackets()

					pppoedb.DatabaseName = "TEST_DATABASE"
					testMsg = &tr477.CpriMsg{
						MetaData: &tr477.CpriMetaData{
							Generic: &tr477.GenericMetadata{
								DeviceName:      "TEST_DEVICE_NAME",
								DeviceInterface: "TEST_VSI_NAME",
							},
						},
						Packet: padiRequest.Data(),
					}

					expectedMsg = &tr477.CpriMsg{
						MetaData: testMsg.MetaData,
						Packet:   padiResponse.Data(),
					}

				})
				Context("with vsi requiring subscriber-profile", func() {
					BeforeAll(func() {
						testVsi = bson.D{
							{Key: "_id", Value: 1},
							{Key: "vsi-name", Value: testMsg.MetaData.Generic.DeviceInterface},
							{Key: "device-name", Value: testMsg.MetaData.Generic.DeviceName},
							{Key: "subscriber-profile", Value: "TEST_SUBSCRIBER_PROFILE"},
						}
					})
					Context("without valid subscriber-profile in database", func() {
						It("should call 'Send' with the test packet", func() {
							stream.On("Recv").Return(testMsg, nil)
							stream.On("Send", testMsg).Return(nil)
							mongoTest.Run("", func(mt *mtest.T) {
								mt.AddMockResponses(
									mtest.CreateCursorResponse(1, ".vsi-list", mtest.FirstBatch,
										testVsi),
								)
								mongoClient = mt.Client
								go func() {
									for len(stream.Calls) < 2 {
									}
									testMessageService.Quit <- true
								}()
								testMessageService.TransferCpri(stream)
							})
							stream.AssertExpectations(GinkgoT())
						})
					})
					Context("with valid subscriber-profile in database", func() {
						BeforeAll(func() {
							discard_on_error = true

							testSubscriberProfile = bson.D{
								{Key: "_id", Value: 1},
								{Key: "name", Value: "TEST_SUBSCRIBER_PROFILE"},
								{Key: "circuit-id", Value: "TEST_CIRCUIT_ID"},
								{Key: "remote-id", Value: "TEST_REMOTE_ID"},
							}

						})
						It("should place subscriber profile info on packet and send it", func() {
							Skip("Work in progress. Skipping test.")
							stream.On("Recv").Return(testMsg, nil)
							stream.On("Send", expectedMsg).Return(nil)
							mongoTest.Run("", func(mt *mtest.T) {
								mt.Client.Database(pppoedb.DatabaseName).Collection("vsi-list")
								mongoClient.Database(pppoedb.DatabaseName).Collection("subscriber-profiles")
								mt.AddMockResponses(
									mtest.CreateCursorResponse(1, pppoedb.DatabaseName+".vsi-list", mtest.FirstBatch,
										testVsi),
									mtest.CreateCursorResponse(2, pppoedb.DatabaseName+".subscriber-profiles", mtest.NextBatch,
										testSubscriberProfile),
								)
								go func() {
									for len(stream.Calls) < 2 {
									}
									testMessageService.Quit <- true
								}()
								mongoClient = mt.Client
								testMessageService.TransferCpri(stream)
							})
							stream.AssertExpectations(GinkgoT())

						})
					})

				})
				Context("with vsi requiring pppoe-profile", func() {
					BeforeAll(func() {
						discard_on_error = false
						testVsi = bson.D{
							{Key: "_id", Value: 1},
							{Key: "vsi-name", Value: testMsg.MetaData.Generic.DeviceInterface},
							{Key: "device-name", Value: testMsg.MetaData.Generic.DeviceName},
							{Key: "pppoe-profile", Value: "TEST_PPPOE_PROFILE"},
						}
					})
					Context("without valid pppoe-profile in database", func() {
						It("should call 'Send' with the test packet", func() {
							stream.On("Recv").Return(testMsg, nil)
							stream.On("Send", testMsg).Return(nil)
							mongoTest.Run("", func(mt *mtest.T) {
								mt.AddMockResponses(
									mtest.CreateCursorResponse(1, ".vsi-list", mtest.FirstBatch,
										testVsi),
								)
								mongoClient = mt.Client
								go func() {
									for len(stream.Calls) < 2 {
									}
									testMessageService.Quit <- true
								}()
								testMessageService.TransferCpri(stream)
							})
							stream.AssertExpectations(GinkgoT())
						})
					})
					// Context("with valid pppoe-profile in database", func() {
					// 	testPPPPoEProfile := bson.D{}
					// 	testVsi := bson.D{
					// 		{Key: "vsi-name", Value: testMsg.MetaData.Generic.DeviceInterface},
					// 		{Key: "device-name", Value: testMsg.MetaData.Generic.DeviceName},
					// 		{Key: "pppoe-profile", Value: testPPPPoEProfile.Map()["name"]},
					// 	}
					// 	expectedMsg := &tr477.CpriMsg{
					// 		MetaData: testMsg.MetaData,
					// 		Packet:   []byte{},
					// 	}
					// 	stream.On("Recv").Return(testMsg, nil)
					// 	stream.On("Send", expectedMsg).Return(nil)
					// 	mongoTest.Run("", func(mt *mtest.T) {
					// 		mt.AddMockResponses( //TODO: change this V
					// 			mtest.CreateCursorResponse(1, pppoedb.DatabaseName+".vsi-list", mtest.FirstBatch,
					// 				testVsi),
					// 			mtest.CreateCursorResponse(1, pppoedb.DatabaseName+".pppoe-profiles", mtest.FirstBatch,
					// 				testPPPPoEProfile),
					// 		)
					// 		go func() {
					// 			time.Sleep(timeout)
					// 			testMessageService.Quit <- true
					// 		}()
					// 		testMessageService.TransferCpri(stream)
					// 		stream.AssertCalled(mt.T, "Recv")
					// 		sentMsg := stream.Calls[1].Arguments[0].(*tr477.CpriMsg)
					// 		sentPacket := gopacket.NewPacket(sentMsg.Packet, layers.LayerTypeEthernet, gopacket.NoCopy).
					// 			Layer(pppoePacket.LayerTypePPPoEDiscovery).(*pppoePacket.PPPoEDiscovery)
					// 		Expect(sentPacket).To(BeEquivalentTo(expectedMsg.Packet))
					// 	})
					// })
				})
			})
		})
	})
})

func getTestPADIPackets() (gopacket.Packet, gopacket.Packet) {
	padiPacketsHandler, err := pcap.OpenOffline("./test_padi.pcap")
	if err != nil {
		Fail("Failed to load pcap")
	}
	padiSource := gopacket.NewPacketSource(padiPacketsHandler, padiPacketsHandler.LinkType())
	padiRequest, err := padiSource.NextPacket()
	if err != nil {
		Fail("Failed to get packet")
	}
	padiResponse, err := padiSource.NextPacket()
	if err != nil {
		Fail("Failed to get packet")
	}
	return padiRequest, padiResponse
}

type MockTransferCpriServer struct {
	mock.Mock
	tr477.CpriMessage_TransferCpriServer
}

func (m *MockTransferCpriServer) Recv() (*tr477.CpriMsg, error) {
	args := m.Called()
	return args.Get(0).(*tr477.CpriMsg), args.Error(1)
}

func (m *MockTransferCpriServer) Send(in *tr477.CpriMsg) error {
	args := m.Called(in)
	return args.Error(0)
}
