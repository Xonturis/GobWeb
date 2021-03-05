package network

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"unsafe"
)

//
// Simply represents the connection between the actual client and another client.
//
type Connection struct {
	netConn    net.Conn
	readWriter *bufio.ReadWriter
}

type Connectable interface {
	Send(Packet)
	Disconnect()
	Listen()
	GetConn() net.Conn
}

func (c *Connection) GetConn() net.Conn {
	return c.netConn
}

func (c *Connection) GetReadWriter() *bufio.ReadWriter {
	return c.readWriter
}

// Source : https://gist.github.com/ecoshub/5be18dc63ac64f3792693bb94f00662f
func IntToByteArray(num int64) []byte {
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	for i := 0 ; i < size ; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[i] = byt
	}
	return arr
}

// Source : https://gist.github.com/ecoshub/5be18dc63ac64f3792693bb94f00662f
func ByteArrayToInt(arr []byte) int64{
	val := int64(0)
	size := len(arr)
	for i := 0 ; i < size ; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[i]
	}
	return val
}

func (c *Connection) Send(packet Packet) {
	fmt.Println("Send: ", packet)
	packetBytes := encodeToBytes(packet)
	fmt.Println("after encode")
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
	_, _ = c.GetReadWriter().Write([]byte("\n")) // \n
	_, _ = c.GetReadWriter().Write(packetBytes) // packet

	_ = c.GetReadWriter().Writer.Flush()

	fmt.Println("Written " + strconv.Itoa(len(packetBytes)))
}

func (c *Connection) Disconnect() {
	_ = c.GetConn().Close()
}

// Source: https://gist.github.com/SteveBate/042960baa7a4795c3565
func encodeToBytes(i interface{}) []byte {
	fmt.Println("1")
	buf := bytes.Buffer{}
	fmt.Println("2")
	enc := gob.NewEncoder(&buf)
	fmt.Println("3")
	err := enc.Encode(i)
	fmt.Println("4")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("5")
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
	//fmt.Println("Start listen")

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

		//fmt.Println(n)
		decodedPacket := decodePacket(bapacket)
		//fmt.Println("DECODED")
		Handle(decodedPacket)
		//fmt.Println("HANDLED")
	}
}

// Établie un lien TCP et return un Connectable.
func ConnectIP(ip net.IP, port int) Connectable {
	conn, err := net.Dial("tcp", ip.String()+":"+strconv.Itoa(port))

	if err != nil {
		panic(err)
	}

	return WrapConnection(conn)
}

func WrapConnection(conn net.Conn) Connectable {
	//fmt.Println("New Connection")
	connectable := &Connection{
		conn,
		bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)),
	}
	go connectable.Listen()
	return connectable
}
