package network

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"unsafe"
)

//
// Simply represents the connection between the actual client and another client.
//
type Connection struct {
	netConn       net.Conn
	listeningPort int
	readWriter    *bufio.ReadWriter
	ipPortAddress string
	LocalAddr     string
}

type Connectable interface {
	Send(Packet)
	Disconnect()
	Listen()
	GetConn() net.Conn
	GetListeningPort() int
	SetListeningPort(port int)
	SetIpPortAddress(ipport string)
	GetIpPortAddress() string
}

func (c *Connection) GetConn() net.Conn {
	return c.netConn
}

func (c *Connection) GetIpPortAddress() string {
	return c.ipPortAddress
}

func (c *Connection) SetIpPortAddress(ipport string)  {
	c.ipPortAddress = ipport
}

func (c *Connection) GetListeningPort() int {
	return c.listeningPort
}

func (c *Connection) SetListeningPort(port int) {
	c.listeningPort = port

	remoteAddr := c.GetConn().RemoteAddr().String()
	remoteAddr = strings.Split(remoteAddr, ":")[0]
	remoteAddr = remoteAddr + ":" + strconv.Itoa(c.GetListeningPort())
	c.ipPortAddress = remoteAddr
}

func (c *Connection) GetReadWriter() *bufio.ReadWriter {
	return c.readWriter
}

// Source : https://gist.github.com/ecoshub/5be18dc63ac64f3792693bb94f00662f
func IntToByteArray(num int64) []byte {
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	for i := 0; i < size; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[i] = byt
	}
	return arr
}

// Source : https://gist.github.com/ecoshub/5be18dc63ac64f3792693bb94f00662f
func ByteArrayToInt(arr []byte) int64 {
	val := int64(0)
	size := len(arr)
	for i := 0; i < size; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[i]
	}
	return val
}

func (c *Connection) Send(packet Packet) {
	packet.PipSrc = c.GetConn().LocalAddr().String()
	packetBytes := encodeToBytes(packet)
	size := len(packetBytes) // Size used at reception to handle the packet (buffer business)

	// TODO handle err

	// Protocol :
	// size
	// ... n bytes
	// \n
	// packet
	// ... size bytes
	// should stop 'cause we know the size

	_, _ = c.GetReadWriter().Write(IntToByteArray(int64(size))) // size
	_, _ = c.GetReadWriter().Write([]byte("\n"))                // \n
	_, _ = c.GetReadWriter().Write(packetBytes)                 // packet

	_ = c.GetReadWriter().Writer.Flush()

}

func (c *Connection) Disconnect() {
	_ = c.GetConn().Close()
}

// Source: https://gist.github.com/SteveBate/042960baa7a4795c3565
func encodeToBytes(i interface{}) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(i)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

// Source: https://gist.github.com/SteveBate/042960baa7a4795c3565
func decodePacket(s []byte) Packet {
	p := Packet{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil {
		log.Fatal(err)
	}
	return p
}

func (c *Connection) Listen() {
	// Protocol :
	// size
	// ... n bytes
	// \n
	// packet
	// ... size bytes
	// should stop 'cause we know the size

	for {
		//var buffer bytes.Buffer
		basize, err := c.GetReadWriter().Reader.ReadBytes('\n') // size

		if err != nil {
			return // disconnected
		}

		basize = basize[:len(basize)-1]
		size := ByteArrayToInt(basize)

		bapacket := make([]byte, size)
		_, err = io.ReadFull(c.GetReadWriter().Reader, bapacket) // packet

		if err != nil {
			return
		}

		decodedPacket := decodePacket(bapacket)
		decodedPacket.Conn = *c
		Handle(decodedPacket)
	}
}

// Établie un lien TCP et return un Connectable.
func ConnectIP(ip net.IP, port int) Connectable {
	conn, err := net.Dial("tcp", ip.String()+":"+strconv.Itoa(port))

	if err != nil {
		return nil
	}

	return WrapConnection(conn)
}

func WrapConnection(conn net.Conn) Connectable {
	connectable := &Connection{
		conn,
		0,
		bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)),
		conn.RemoteAddr().String(),
		conn.LocalAddr().String(),
	}
	go connectable.Listen()
	return connectable
}
