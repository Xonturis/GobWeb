package network

func CreatePacket(ipDest string, data []byte) Packet {
	return Packet{PipDest: ipDest, PipSrc: SelfServerAddress, Pdata: data}
}
