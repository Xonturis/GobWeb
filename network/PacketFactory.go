package network

func CreatePacket(ipDest string, packetType string, data interface{}) Packet {
	return Packet{PipDest: ipDest, PipSrc: GetSelfIPAddress(), Ptype: packetType, Pdata: data}
}
