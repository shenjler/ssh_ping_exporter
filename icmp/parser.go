package icmp

import (
	"regexp"
	"strings"

	"github.com/shenjler/ssh_ping_exporter/util"
)

// Parse parses cli output and tries to find interfaces with related stats
func (c *icmpCollector) Parse(ostype string, output string) (Icmp, error) {
	// if ostype != rpc.IOSXE && ostype != rpc.NXOS && ostype != rpc.IOS {
	// 	return nil, errors.New("'show interface' is not implemented for " + ostype)
	// }

	targetRegexp := regexp.MustCompile(`^\s*--- (.*) ping statistics ---.*$`)                                                                      // target
	packetLossRegexp := regexp.MustCompile(`^\s*(?:(?:\d+) packets transmitted, (?:\d+) received, )?((?:[1-9][\d]*|0)(?:\.\d+)?)% packet loss.*$`) // packet loss rate
	rttRegexp := regexp.MustCompile(`^\s*(?:rtt|round-trip)? min/avg/max(?:/mdev)? = ((?:[1-9][\d]*|0)(?:\.[\d]+)?)/((?:[1-9][\d]*|0)(?:\.[\d]+)?)/((?:[1-9][\d]*|0)(?:\.[\d]+)?)(?:/.*)? ms.*$`)

	current := Icmp{}
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if matches := targetRegexp.FindStringSubmatch(line); matches != nil {
			current = Icmp{
				Target: matches[1],
			}
		}
		if current == (Icmp{}) {
			continue
		}
		if matches := packetLossRegexp.FindStringSubmatch(line); matches != nil {
			current.PacketLoss = util.Str2float64(matches[1])
			if current.PacketLoss == 100 {
				current.PingStatus = "down"
			} else {
				current.PingStatus = "up"
			}
		}
		if matches := rttRegexp.FindStringSubmatch(line); matches != nil {
			current.RttMin = util.Str2float64(matches[1])
			current.RttAvg = util.Str2float64(matches[2])
			current.RttMax = util.Str2float64(matches[3])
		}

	}
	return current, nil
}
