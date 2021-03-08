package network

// Un packet est un élément d'information abstrait qui peut voyager a travers le réseau.
//
//  Connection   conn    La connexion de l'émetteur du packet
//  string       Ptype   Le type du packet
//  interface{}  Pdate   Les données du packet
//
type Packet struct {
	conn   *Connection
	Ptype  string
	Pdata  interface{}
}

type PacketConnectable interface {
	GetConnection() *Connection
	SetConnection(connection *Connection)
}

// Fonction permettant de récupérer la connexion.
func (p *Packet) GetConnection() *Connection {
	return p.conn
}

// Fonction permettant de définir la connexion
func (p *Packet) SetConnection(connection *Connection)  {
	p.conn = connection
}

// Fonction permettant de créer un packet.
//
//  string       packetType  Le type du packet
//  interface{}  data        Les données du packet
//
func CreatePacket(packetType string, data interface{}) Packet {
	return Packet{Ptype: packetType, Pdata: data}
}
