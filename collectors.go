package main

import (
	"github.com/shenjler/ssh_ping_exporter/collector"
	"github.com/shenjler/ssh_ping_exporter/config"
	"github.com/shenjler/ssh_ping_exporter/connector"
	"github.com/shenjler/ssh_ping_exporter/icmp"
)

type collectors struct {
	collectors map[string]collector.RPCCollector
	devices    map[string][]collector.RPCCollector
	cfg        *config.Config
}

func collectorsForDevices(devices []*connector.Device, cfg *config.Config) *collectors {
	c := &collectors{
		collectors: make(map[string]collector.RPCCollector),
		devices:    make(map[string][]collector.RPCCollector),
		cfg:        cfg,
	}

	for _, d := range devices {
		c.initCollectorsForDevice(d)
	}

	return c
}

func (c *collectors) initCollectorsForDevice(device *connector.Device) {
	f := c.cfg.FeaturesForDevice(device.Host)

	c.devices[device.Host] = make([]collector.RPCCollector, 0)
	c.addCollectorIfEnabledForDevice(device, "icmp", f.Icmp, icmp.NewCollector)

	// c.addCollectorIfEnabledForDevice(device, "bgp", f.BGP, bgp.NewCollector)
	// c.addCollectorIfEnabledForDevice(device, "environment", f.Environment, environment.NewCollector)
	// c.addCollectorIfEnabledForDevice(device, "facts", f.Facts, facts.NewCollector)
	// c.addCollectorIfEnabledForDevice(device, "interfaces", f.Interfaces, interfaces.NewCollector)
	// c.addCollectorIfEnabledForDevice(device, "optics", f.Optics, optics.NewCollector)

}

func (c *collectors) addCollectorIfEnabledForDevice(device *connector.Device, key string, enabled *bool, newCollector func() collector.RPCCollector, dest string) {
	if !*enabled {
		return
	}

	col, found := c.collectors[key]
	if !found {
		col = newCollector()
		c.collectors[key] = col
	}

	c.devices[device.Host] = append(c.devices[device.Host], col)
}

func (c *collectors) allEnabledCollectors() []collector.RPCCollector {
	collectors := make([]collector.RPCCollector, len(c.collectors))

	i := 0
	for _, collector := range c.collectors {
		collectors[i] = collector
		i++
	}

	return collectors
}

func (c *collectors) collectorsForDevice(device *connector.Device) []collector.RPCCollector {
	cols, found := c.devices[device.Host]
	if !found {
		return []collector.RPCCollector{}
	}

	return cols
}
