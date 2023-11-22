package main

import (
	// "clients"

	// "github.com/yankiwi/wlc_exporter/aps"
	// "github.com/yankiwi/wlc_exporter/clients"
	"github.com/yankiwi/wlc_exporter/collector"
	"github.com/yankiwi/wlc_exporter/config"

	// "github.com/yankiwi/aruba_exporter/collector"
	"github.com/yankiwi/wlc_exporter/apessid"

	// "github.com/yankiwi/wlc_exporter/config"
	"github.com/yankiwi/wlc_exporter/connector"
	// "github.com/yankiwi/wlc_exporter/inp"
	// "github.com/yankiwi/wlc_exporter/interfaces"sudo
	"github.com/yankiwi/wlc_exporter/system"
	// "github.com/yankiwi/wlc_exporter/wireless"
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
	c.addCollectorIfEnabledForDevice(device, "system", f.System, system.NewCollector)
	// c.addCollectorIfEnabledForDevice(device, "interfaces", f.Interfaces, interfaces.NewCollector)
	// c.addCollectorIfEnabledForDevice(device, "clients", f.Clients, clients.NewCollector)
	// c.addCollectorIfEnabledForDevice(device, "aps", f.Aps, aps.NewCollector)
	// c.addCollectorIfEnabledForDevice(device, "wireless", f.Wireless, wireless.NewCollector)
	// c.addCollectorIfEnabledForDevice(device, "inp", f.Inp, inp.NewCollector)
	c.addCollectorIfEnabledForDevice(device, "apessid", f.Apessidd, apessid.NewCollector)

}

func (c *collectors) addCollectorIfEnabledForDevice(device *connector.Device, key string, enabled *bool, newCollector func() collector.RPCCollector) {
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
