package icmp

import (
	"log"

	"github.com/shenjler/ssh_ping_exporter/rpc"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shenjler/ssh_ping_exporter/collector"
)

const prefix string = "pccw_icmp_"

var (
	packetLossDesc *prometheus.Desc
	pingStatusDesc *prometheus.Desc
	rttAvgDesc     *prometheus.Desc
)

func init() {
	l := []string{"dest", "source"}
	packetLossDesc = prometheus.NewDesc(prefix+"packet_loss", "The ping packet loss rate: 0~100", l, nil)
	rttAvgDesc = prometheus.NewDesc(prefix+"rtt_ms", "The avg rtt of the ping", l, nil)
	pingStatusDesc = prometheus.NewDesc(prefix+"status", "Status of ping, when packet loss is 100, the status is 0. ", l, nil)

}

type icmpCollector struct {
	dest string
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	return &icmpCollector{
		// dest: dest,
	}
}

// Name returns the name of the collector
func (*icmpCollector) Name() string {
	return "Icmp"
}

// Describe describes the metrics
func (*icmpCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- packetLossDesc
	ch <- rttAvgDesc
	ch <- pingStatusDesc
}

// Collect collects metrics from Cisco
func (c *icmpCollector) CollectByDest(client *rpc.Client, ch chan<- prometheus.Metric, labelValues []string) error {
	out, err := client.RunCommand("ping -c 5 " + c.dest)
	if err != nil {
		return err
	}
	item, err := c.Parse(client.OSType, out)
	if err != nil {
		if client.Debug {
			log.Printf("Parse ping for %s: %s\n", labelValues[0], err.Error())
		}
		return nil
	}

	l := append(labelValues, item.Target)
	ch <- prometheus.MustNewConstMetric(packetLossDesc, prometheus.GaugeValue, float64(item.PacketLoss), l...)

	if item.PingStatus == "up" {
		ch <- prometheus.MustNewConstMetric(rttAvgDesc, prometheus.GaugeValue, float64(item.RttAvg), l...)
		ch <- prometheus.MustNewConstMetric(pingStatusDesc, prometheus.GaugeValue, 1, l...)

	} else {
		ch <- prometheus.MustNewConstMetric(pingStatusDesc, prometheus.GaugeValue, 0, l...)
	}
	return nil
}
