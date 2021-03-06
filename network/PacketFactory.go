package network

func CreatePacket(packetType string, data interface{}) Packet {
	return Packet{Ptype: packetType, Pdata: data}
}
