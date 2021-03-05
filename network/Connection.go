package network

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strconv"
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

func (c *Connection) Send(packet Packet) {

	packetBytes := encodeToBytes(packet)
	_, err := c.GetReadWriter().Write(packetBytes)
	if err != nil {
		fmt.Println(err)
	}
	err = c.GetReadWriter().Writer.Flush()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Written " + strconv.Itoa(len(packetBytes)))
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
	fmt.Println("Start listen")
	for {
		readable := c.GetReadWriter().Reader.Buffered()
		if readable > 0 {
			read := make([]byte, readable)
			n, err := c.GetReadWriter().Read(read)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("RECEIVED ", n, " bytes over ", readable)
			decodedPacket := decodePacket(read)
			fmt.Println("DECODED")
			Handle(decodedPacket)
			fmt.Println("HANDLED")
		}

	}
}

func ConnectIP(ip net.IP, port int) Connectable {
	conn, err := net.Dial("tcp", ip.String()+":"+strconv.Itoa(port))

	if err != nil {
		panic(err)
	}

	return Connect(conn)
}

func Connect(conn net.Conn) Connectable {
	fmt.Println("New Connection")
	connectable := &Connection{
		conn,
		bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)),
	}
	go connectable.Listen()
	return connectable
}
