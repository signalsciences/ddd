package hostmetrics

import (
	"github.com/signalsciences/dogdirect"
)

// HostMetricsWriter collects and writes host metrics to datadog
type HostMetricWriter struct {
	ddog      *dogdirect.Client
	collector *HostMetricCollector
	tags      []string
}

// NewHostMetricWriter will collect hostmetrics and send them to the datadog client on flush.  it will also apply global tags to the system metrics.
func NewHostMetricWriter(ddog *dogdirect.Client, tags []string) (*HostMetricWriter, error) {
	collector, err := NewHostMetricCollector()
	if err != nil {
		return nil, err
	}
	return &HostMetricWriter{
		ddog:      ddog,
		collector: collector,
		tags:      tags,
	}, nil
}

// Flush gets new host metrics, inserts them into the datadog client, and
// then flushes the metrics upstream
func (h *HostMetricWriter) Flush() error {
	hr, err := h.collector.Run()
	if err != nil {
		return err
	}
	h.ddog.Gauge("xsystem.cpu.user", hr.CPUUser, h.tags)
	h.ddog.Gauge("xsystem.cpu.system", hr.CPUSystem, h.tags)
	h.ddog.Gauge("xsystem.cpu.iowait", hr.CPUIowait, h.tags)
	h.ddog.Gauge("xsystem.cpu.idle", hr.CPUIdle, h.tags)
	h.ddog.Gauge("xsystem.cpu.stolen", hr.CPUStolen, h.tags)
	h.ddog.Gauge("xsystem.cpu.guest", hr.CPUGuest, h.tags)
	/*
	h.ddog.Gauge("system.mem.total", hr.MemTotal, h.tags)
	h.ddog.Gauge("system.mem.free", hr.MemFree, h.tags)
	h.ddog.Gauge("system.mem.used", hr.MemUsed, h.tags)
	h.ddog.Gauge("system.mem.usable", hr.MemUsable, h.tags)
	h.ddog.Gauge("system.mem.pct_usable", hr.MemPctUsable, h.tags)
	*/
	return h.ddog.Flush()
}

// Close closes the datadog client
func (h *HostMetricWriter) Close() error {
	return h.ddog.Close()
}
