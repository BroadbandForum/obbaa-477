package main

import (
	"context"
	"flag"
	"io"

	tr477 "github.com/obbaa-477/common/pb/tr477"

	log "github.com/obbaa-477/common/utils/log"
	osadapters "github.com/obbaa-477/common/utils/osadapters"

	"net"
	"os"
	"time"

	"github.com/google/gopacket/layers"

	"github.com/mdlayher/raw"
	"google.golang.org/grpc"
)

var clientName = os.Getenv("HOSTNAME")
var StationID string
var clientSideIfaceName string
var serverSideIfaceName string
var quit = make(chan int)

/*two connections are required to receive packets, as pppoe discovery and session have different ethertypes,
but any one of those connections may be used to send packets
*/
type pppoeInterface struct {
	name                       string
	discoveryConn, sessionConn *raw.Conn
}

func main() {

	var serverAddr string

	flag.StringVar(&clientSideIfaceName, "client_if", "", "available interfaces for supplicants(e.g. ens5,ens6,ens7,ens8)")
	flag.StringVar(&serverSideIfaceName, "server_if", "", "available interfaces for supplicants(e.g. ens5,ens6,ens7,ens8)")
	flag.StringVar(&serverAddr, "vnf-addr:port", "", "IP and port of VNF (e.g. X.X.X.X:50051)")
	flag.Parse()

	if serverAddr == "" || clientSideIfaceName == "" || serverSideIfaceName == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parse()

	// Allow only single instance of goeap_proxy
	// We could potentially tie the lock file to the wan and rtr interfaces
	// But lets keep things simple for now
	l, err := net.Listen("unix", "@/run/pppoe_proxy.lock")
	if err != nil {
		log.Warning(os.Stderr, "pppoe_proxy is already running!")
		os.Exit(1)
	}
	defer l.Close()

	// Start gRPC Connection
	cc, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("could not connect: ", err)
	}
	defer cc.Close()
	log.Info("########### Starting OLT_App - gRPC Client is initializing ###########")

	addHelloServiceClient := tr477.NewCpriHelloClient(cc)
	addPacketServiceClient := tr477.NewCpriMessageClient(cc)

	helloService(addHelloServiceClient)

	stream, err := addPacketServiceClient.TransferCpri(context.Background())
	if err != nil {
		log.Error("Error TransferCpri: ", err)
		quit <- 1
	}
	// Go rotine to receive data from VNF via gRPC
	go waitForPacketsOnStream(stream)

	setupInterface(serverSideIfaceName, stream)
	setupInterface(clientSideIfaceName, stream)

	<-quit
}

func setupInterface(ifaceName string, stream tr477.CpriMessage_TransferCpriClient) {
	iface := newInterface(ifaceName)
	ifaceMAC, ifaceMACFormat := osadapters.GetMAC(ifaceName)
	log.Info("OLT_app - Start listening on port: " + ifaceName + " | MAC Addr: " + ifaceMACFormat)

	time.Sleep((500 * time.Millisecond))
	go listenForPPPoEDiscoveryPackets(iface, ifaceMAC, stream)
	go listenForPPPoESessionPackets(iface, ifaceMAC, stream)

}

func helloService(addHelloService tr477.CpriHelloClient) {

	req := &tr477.HelloCpriRequest{
		LocalEndpointHello: &tr477.Hello{
			EntityName:   clientName,
			EndpointName: clientName,
		},
	}
	_, err := addHelloService.HelloCpri(context.Background(), req)
	if err != nil {
		log.Error("- error while calling Hello RPC: ", err)
		quit <- 1
	}
}

func newInterface(name string) *pppoeInterface {

	intf, err := net.InterfaceByName(name)
	if err != nil {
		log.Error(os.Stderr, "- InterfaceByName failed: ", name, err)
		quit <- 1
	}

	discoveryConn, err := raw.ListenPacket(intf, uint16(layers.EthernetTypePPPoEDiscovery), nil)
	if err != nil {
		log.Error(os.Stderr, "- ListenPacket for discovery failed: ", name, err)
		quit <- 1
	}
	sessionConn, err := raw.ListenPacket(intf, uint16(layers.EthernetTypePPPoESession), nil)
	if err != nil {
		log.Error(os.Stderr, "- ListenPacket for session failed: ", name, err)
		quit <- 1
	}

	pppoeIntf := pppoeInterface{name: name,
		discoveryConn: discoveryConn,
		sessionConn:   sessionConn}

	return &pppoeIntf
}

func listenForPPPoESessionPackets(src *pppoeInterface, oltIfaceMAC []byte, stream tr477.CpriMessage_TransferCpriClient) {
	// This might break for jumbo frames
	recvBuf := make([]byte, 1500)
	for {

		size, _, err := src.sessionConn.ReadFrom(recvBuf)
		if err != nil {
			log.Error(os.Stderr, "- : unexpected read error: ", src.name, err)
			// maybe not necessary, give the system a minute to recover
			time.Sleep(500 * time.Millisecond)
			continue
		}
		log.Info("Received pppoe session packet")

		packetData := recvBuf[:size]

		// Send a gRPC message to VNF
		err = stream.Send(&tr477.CpriMsg{
			MetaData: &tr477.CpriMetaData{
				Generic: &tr477.GenericMetadata{
					DeviceName:      clientName,
					DeviceInterface: src.name,
					Direction:       getPacketDirection(src.name),
				},
			},
			Packet: packetData})
		if err != nil {
			log.Error("Error: ", err)
		}
	}
}

// Function to receive ethernet frame
func listenForPPPoEDiscoveryPackets(src *pppoeInterface, oltIfaceMAC []byte, stream tr477.CpriMessage_TransferCpriClient) {
	// This might break for jumbo frames
	recvBuf := make([]byte, 1500)
	for {
		size, _, err := src.discoveryConn.ReadFrom(recvBuf)
		if err != nil {
			log.Error(os.Stderr, "- : unexpected read error: ", src.name, err)
			// maybe not necessary, give the system a minute to recover
			time.Sleep(500 * time.Millisecond)
			continue
		}
		log.Info("Received pppoe discovery packet")

		packetData := recvBuf[:size]

		// Send a gRPC message to VNF
		err = stream.Send(&tr477.CpriMsg{
			MetaData: &tr477.CpriMetaData{
				Generic: &tr477.GenericMetadata{
					DeviceName:      clientName,
					DeviceInterface: src.name,
					Direction:       getPacketDirection(src.name),
				},
			},
			Packet: packetData})
		if err != nil {
			log.Error("Error: ", err)
		}
	}
}

// Func to receive data from VNF, uncap and send supplicant
func waitForPacketsOnStream(stream tr477.CpriMessage_TransferCpriClient) {

	go on()

	for {
		packet, err := stream.Recv()
		if err == io.EOF {
			log.Fatal("Erro: ", err)
		}
		if packet != nil {
			log.Info("OLT_app - Received - gRPC message from Control_relay")

			var rtr *pppoeInterface
			if packet.MetaData.Generic.Direction == tr477.GenericMetadata_UNI_TO_NNI {
				rtr = newInterface(serverSideIfaceName)
			} else if packet.MetaData.Generic.Direction == tr477.GenericMetadata_NNI_TO_UNI {
				rtr = newInterface(clientSideIfaceName)
			} else {
				log.Error("Unrecognized packet direction: ", packet.MetaData.Generic.Direction)
				continue
			}

			//any one of the connections can be used to send packets
			_, err = rtr.discoveryConn.WriteTo(packet.Packet, nil)
			if err != nil {
				log.Error(os.Stderr, ": unexpected write error: ", packet.MetaData.Generic.DeviceInterface, err)
			}

		}
	}
}

// Send signal to OLT backplane that port is authenticated and authorized
func UnlockPort(olt_port string) {
	log.Warning("Interface is authenticated and authorized: " + olt_port)

}

// Send signal to OLT backplane that port is not authenticated and authorized
func LockPort(olt_port string) {
	log.Warning("Interface is not authenticated and authorized: " + olt_port)

}
func on() {

	for {
		log.Info("OLT_app - Waiting for packets on stream")
		time.Sleep(10 * time.Second)
	}
}

func getPacketDirection(srcIntf string) tr477.GenericMetadata_PacketDirection {
	if srcIntf == clientSideIfaceName {
		return tr477.GenericMetadata_UNI_TO_NNI
	} else if srcIntf == serverSideIfaceName {
		return tr477.GenericMetadata_NNI_TO_UNI
	} else {
		log.Error("Unknown interface: ", srcIntf)
	}
	return -1
}
