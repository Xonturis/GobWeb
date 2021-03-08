package network

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

//
//Cobweb :
//The center of the inner network of the client where the program is running.
//This piece will be useful to send packets to all connections, accept connections etc...
//

var allCurrentConnectables = make([]Connectable, 0)
var packetsChan = make(chan Packet)
var selfServerAddress = "127.0.0.1"
var SelfServerPort int
var listener net.Listener

type packetTypeHandler = func(packet Packet)

var packetTypeHandlerMap = make(map[string]packetTypeHandler)

var OnReady func()

func GetConnectable(ip string) Connectable {
	fmt.Println("All connectables: ")
	for _, conn := range allCurrentConnectables {
		fmt.Println(conn.GetConn().RemoteAddr().String())
		if strings.Index(conn.GetConn().RemoteAddr().String(), ip) >= 0 {
			return conn
		}
	}
	return nil
}

func GetAllConnectedIPListeningPortString() []string {
	ips := make([]string, 0, len(allCurrentConnectables))

	for _, conn := range allCurrentConnectables {
		ips = append(ips, conn.GetIpPortAddress())
	}

	return ips
}

func accept(port int) {
	// Accepts new connections
	var err error
	listener, err = net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		return
	}
	fmt.Println("Listening to : " + listener.Addr().String())

	selfServerAddress = "127.0.0.1"
	SelfServerPort = port

	for {
		conn, err := listener.Accept()
		fmt.Println("New connection: ", conn.RemoteAddr().String())
		if err != nil {
			os.Exit(6)
		}
		AddToNetwork(conn)
	}
}

func SendToAll(packet Packet) {

	for _, connection := range allCurrentConnectables {
		connection.Send(packet)
	}

}

func GracefulShutdown() {
	// Todo, not implemented yet
}

func Handle(packet Packet) {
	//fmt.Println("Handle: ", packet)
	packetsChan <- packet
}

func handlePackets() {
	for {
		packet := <-packetsChan
		onReceive(packet)
	}
}

func onReceive(packet Packet) {
	if packetTypeHandlerMap[packet.Ptype] != nil {
		packetTypeHandlerMap[packet.Ptype](packet)
		return
	} else {
		fmt.Println("Received unhandled packet") // Pit, unhandled packet type
	}
}

func RegisterHandler(packetType string, handler packetTypeHandler) {
	packetTypeHandlerMap[packetType] = handler
}

// Connecte a une connexion.
func AddToNetwork(conn net.Conn) {
	split := strings.Split(conn.RemoteAddr().String(), ":")
	port, _ := strconv.Atoi(split[1])
	newConnectable := WrapConnection(conn)
	newConnectable.SetListeningPort(port)
	allCurrentConnectables = append(allCurrentConnectables, newConnectable)
}

func ConnectCobweb(port int, ip net.IP, localport int) {
	// handshake
	go handlePackets()
	SelfServerPort = localport
	Handshake(ip, port)
	// start server
	//go handlePackets()
	//ipConnectableMap[ip.String()] = ConnectIP(ip, port)
	//go accept(localport)
	//callOnReady()
}

func StartCobweb(port int) {
	RegisterHandshakeHandler()
	go handlePackets()
	go accept(port)
	callOnReady()
}

func callOnReady() {
	if OnReady != nil {
		OnReady()
	}
}
func GetSelfIPPortAddress() string {
	return selfServerAddress + ":" + strconv.Itoa(SelfServerPort)
}
