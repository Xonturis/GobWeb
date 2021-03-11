package network

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"strconv"
	"unsafe"
)

// Représente une connexion entre un client et un autre client.
//
// net.Conn            netConn        La connexion
// int                 listeningPort  Le port d'écoute
// *bufio.ReadWritter  readWriter     Permet de lire et écrire dans le buffer de la connexion
// string              ipPortAddress  Le port
// string              LocalAddr      L'ip local
//
type Connection struct {
	netConn        net.Conn
	readWriter     *bufio.ReadWriter
	ipPortAddress  string
}

type Connectable interface {
	Send(Packet)
	Listen()
	GetConn() net.Conn
	SetIpPortAddress(ipport string)
	GetIpPortAddress() string
}

// Fonction permettant de récupérer la connexion.
func (c *Connection) GetConn() net.Conn {
	return c.netConn
}

// Fonction permettant de récupérer l'ip et le port.
func (c *Connection) GetIpPortAddress() string {
	return c.ipPortAddress
}

// Fonction permettant de changer l'ip et le port.
func (c *Connection) SetIpPortAddress(ipport string)  {
	c.ipPortAddress = ipport
}

// Fonction permettant de récupérer readWriter.
func (c *Connection) GetReadWriter() *bufio.ReadWriter {
	return c.readWriter
}

// Fonction convertissant un entier en tableau de byte.
//
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

// Fonction convertissant un tableau de byte en entier.
//
// Source : https://gist.github.com/ecoshub/5be18dc63ac64f3792693bb94f00662f
func ByteArrayToInt(arr []byte) int64 {
	val := int64(0)
	size := len(arr)
	for i := 0; i < size; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[i]
	}
	return val
}

// Fonction permettant d'envoyer un packet sur la connexion.
func (c *Connection) Send(packet Packet) {
	packetBytes := encodeToBytes(packet)
	size := len(packetBytes) // Size used at reception to handle the packet (buffer business)

	// Protocol :
	// size
	// ... n bytes
	// \n
	// packet
	// ... size bytes
	// should stop 'cause we know the size

	var err error
	_, err = c.GetReadWriter().Write(IntToByteArray(int64(size))) // size
	if err != nil {
		Warning(err) // Not asked to handle that case
		return
	}
	_, err = c.GetReadWriter().Write([]byte("\n"))                // \n
	if err != nil {
		Warning(err) // Not asked to handle that case
		return
	}
	_, err = c.GetReadWriter().Write(packetBytes)                 // packet
	if err != nil {
		Warning(err) // Not asked to handle that case
		return
	}

	err = c.GetReadWriter().Writer.Flush()
	if err != nil {
		Warning(err) // Not asked to handle that case
		return
	}

	Info("|-->", packet)
}

// Fonction encodant n'importe quoi en un tableau de byte.
//
// Source: https://gist.github.com/SteveBate/042960baa7a4795c3565
func encodeToBytes(i interface{}) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(i)
	if err != nil {
		fmt.Println("Error with gob, try adding `gob.Register(YOUR_INTERFACE_HERE{})` before handling packets")
		Error(err)
	}
	return buf.Bytes()
}

// Fonction décodant un tableau de byte en un packet
//
// Source: https://gist.github.com/SteveBate/042960baa7a4795c3565
func decodePacket(s []byte) Packet {
	p := Packet{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil {
		fmt.Println("Error with gob, try adding `gob.Register(YOUR_INTERFACE_HERE{})` before handling packets")
		Error(err)
	}
	return p
}

// Fonction écoutant sur le réseau pour savoir si un packet est reçu.
func (c *Connection) Listen() {
	// Problème original :
	//  Comment savoir la taille du message tout en évitant d'imposer des limitations comme une taille maximum
	//  ou un caractère interdit, etc etc
	//  donc on s'est inspiré du protocole IP
	// Protocole :
	// taille
	// ... n bytes
	// \n
	// packet <-- contenu de l'utilisateur, pas de délimiteur puisqu'on connait à l'avance la taille
	// ... taille bytes

	for {
		basize, err := c.GetReadWriter().Reader.ReadBytes('\n') // size

		if err != nil {
			return // disconnected
		}

		basize = basize[:len(basize)-1]
		size := ByteArrayToInt(basize)

		bapacket := make([]byte, size)
		_, err = io.ReadFull(c.GetReadWriter().Reader, bapacket) // packet

		if err != nil {
			Warning(err) // Not asked to handle that case
			return
		}

		decodedPacket := decodePacket(bapacket)
		decodedPacket.SetConnection(c)
		Handle(decodedPacket)
	}
}

// Établie un lien TCP et return un Connectable.
func ConnectIP(ip net.IP, port int) *Connection {
	conn, err := net.Dial("tcp", ip.String()+":"+strconv.Itoa(port))

	if err != nil {
		Error(err)// Deux cas : le premier pair où on se connecte est fermé ou l'un des pairs du réseau s'est déconnecté
		return nil
	}

	return WrapConnection(conn)
}

// Fonction wrappant la net.Conn avec des données utiles pour le reste du programme (readWriter, ip et port)
// dans un Connection
//
//  net.Conn conn	la connection a wrapper
func WrapConnection(conn net.Conn) *Connection {
	var connectable *Connection
	connectable = &Connection{
		conn,
		bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)),
		conn.RemoteAddr().String(),
	}
	go connectable.Listen()
	return connectable
}
