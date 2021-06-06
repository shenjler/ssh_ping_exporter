package icmp

type Icmp struct {
	Name   string
	Target string
	Source string

	PingStatus string

	PacketLoss float64
	RttMin     float64
	RttMax     float64
	RttAvg     float64
}
