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

var ipConnectableMap = make(map[string]Connectable)
var packetsChan chan Packet = make(chan Packet)
var SelfServerAddress string
var SelfServerPort int
var listener net.Listener

func GetConnectable(ip string) Connectable {
	return ipConnectableMap[ip] // Todo handle map err if any ?
}

func GetConnectables() *[]Connectable {
	connectables := make([]Connectable, len(ipConnectableMap))
	for ip, conn := range ipConnectableMap {
		_ = append(connectables, conn)
		fmt.Println(conn == nil)
		fmt.Println(ip)
		fmt.Println(conn.GetConn() == nil)
	}
	return &connectables
}

func setConnectable(ip string, connectable Connectable) {
	ipConnectableMap[ip] = connectable
}

func accept(port int) {
	// Accepts new connections
	var err error
	listener, err = net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Listening to : " + listener.Addr().String())

	SelfServerAddress = "127.0.0.1"
	SelfServerPort = port

	for {
		conn, err := listener.Accept()
		fmt.Println("NEW CONNECTION")
		if err != nil {
			fmt.Println(err)
			os.Exit(6)
		}
		ConnectTo(conn)
	}
}

func SendToAll(packet Packet) {
	packet.PipDest = "255.255.255.255"

	//for _, connection := range *GetConnectables() {
	//	fmt.Println(connection == nil)
	//	connection.Send(packet)
	//	//fmt.Println("SENT " + strconv.Itoa(id))
	//}

	ipConnectableMap["127.0.0.1"].Send(packet)
}

func GracefulShutdown() {
	// Todo, not implemented yet
}

func Handle(packet Packet) {
	fmt.Println(packet.PipDest)
	if packet.PipDest == SelfServerAddress || packet.PipDest == "255.255.255.255" {
		packetsChan <- packet
	}
}

func handlePackets() {
	for {
		packet := <-packetsChan
		onReceive(packet)
	}
}

func onReceive(packet Packet) {
	// Todo implement this
	fmt.Println(packet)
}

func ConnectTo(conn net.Conn) {
	ip := strings.Split(conn.RemoteAddr().String(), ":")[0]
	ipConnectableMap[ip] = Connect(conn)
}

func ConnectCobweb(port int, ip net.IP, localport int) {
	go handlePackets()
	ipConnectableMap[ip.String()] = ConnectIP(ip, port)
	accept(localport)
}

func StartCobweb(port int) {
	go handlePackets()
	accept(port)
}
