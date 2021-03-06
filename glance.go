package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"gopkg.in/niedbalski/goose.v3/client"
	"gopkg.in/niedbalski/goose.v3/glance"
)

type GlanceExporter struct {
	BaseOpenStackExporter
	Client *glance.Client
}

var defaultGlanceMetrics = []Metric{
	{Name: "images"},
}

func NewGlanceExporter(client client.AuthenticatingClient, prefix string, config *Cloud) (*GlanceExporter, error) {
	exporter := GlanceExporter{BaseOpenStackExporter{
		Name:   "glance",
		Prefix: prefix,
		Config: config,
	}, glance.New(client)}

	for _, metric := range defaultGlanceMetrics {
		exporter.AddMetric(metric.Name, metric.Labels, nil)
	}

	return &exporter, nil
}

func (exporter *GlanceExporter) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range exporter.Metrics {
		ch <- metric
	}
}
func (exporter *GlanceExporter) Collect(ch chan<- prometheus.Metric) {
	log.Infoln("Fetching images list")
	images, err := exporter.Client.ListImagesV2()
	if err != nil {
		log.Errorf("%s", err)
	}

	ch <- prometheus.MustNewConstMetric(exporter.Metrics["images"],
		prometheus.GaugeValue, float64(len(images)))
}
