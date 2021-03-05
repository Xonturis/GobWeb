package network

//
//A packet is an abstract piece of information that can travel through the cobweb.
//It has a target and Pdata.
//

type Packet struct {
	PipDest string
	PipSrc  string
	Pdata   []byte
}
