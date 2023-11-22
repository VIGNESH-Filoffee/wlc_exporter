package apessid

import (
	"github.com/yankiwi/wlc_exporter/collector"
	"github.com/yankiwi/wlc_exporter/rpc"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const prefix string = "wlc_essid"

var (
	apessidDesc *prometheus.Desc
)

func init() {
	l := []string{"target", "serial", "essid", "aps", "mbssidtxbss", "clients", "vlans", "encryption"}
	apessidDesc = prometheus.NewDesc(prefix+"access_point", "", l, nil)

}

type apessidCollector struct {
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	return &apessidCollector{}
}

// Name returns the name of the collector
func (*apessidCollector) Name() string {
	return "clients"
}

// Describe describes the metrics
func (*apessidCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- apessidDesc
}

// Collect collects metrics from wlc
func (c *apessidCollector) Collect(client *rpc.Client, ch chan<- prometheus.Metric, labelValues []string) error {

	var (
		out   string
		items map[string]Apess
		err   error
	)

	// switch client.OSType {
	// case "wlcInstant":

	// 	out, err = client.RunCommand([]string{"show clients"})
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	switch client.OSType {
	case "ArubaController":

		out, err = client.RunCommand([]string{"show ap essid"})
		if err != nil {
			return err
		}
	}

	items, err = c.Parse(client.OSType, out)
	if err != nil {
		log.Warnf("Parse clients failed for %s: %s\n", labelValues[0], err.Error())
		return nil
	}

	for _, apessData := range items {

		l := append(labelValues, apessData.essid, apessData.aps, apessData.mbssidtxbss, apessData.clients, apessData.vlans, apessData.encryption)
		//  This line should be changed

		ch <- prometheus.MustNewConstMetric(apessidDesc, prometheus.GaugeValue, 1, l...)
	}

	return nil
}
