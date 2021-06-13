// Copyright 2018 Paychex Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"runtime"

	"github.com/jamiealquiza/envy"
	"github.com/paychex/prometheus-isilon-exporter/pkg/collector"
	isiconfig "github.com/paychex/prometheus-isilon-exporter/pkg/config"
	"github.com/paychex/prometheus-isilon-exporter/pkg/isiclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

const (
	namespace = "emcisi" // used for prometheus metrics
)

var (
	log        = logrus.New()
	config     *isiconfig.Config
	debugLevel = flag.Bool("debug", false, "enable debug messages")

	// date is a time label of the moment when the binary was built
	date = "unset"
	// commit is a last commit hash at the moment when the binary was built
	commit = "unset"
	// version is a semantic version of current build
	version = "unset"

	// Metrics about the EMC Isilon exporter itself.
	isiCollectionBuildInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "emcisi_collector_build_info",
			Help: "A metric with a constant '1' value labeled by version, commitid and goversion exporter was built",
		},
		[]string{"version", "commitid", "goversion"},
	)
	isiExporterUp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "emcisi_exporter_up",
			Help: "Indicates if scrape was successful or not.",
		},
		[]string{"clustername"},
	)
)

func init() {
	log.Formatter = new(logrus.TextFormatter)
	envy.Parse("ISIENV") // looks for ISIENV_USERNAME, ISIENV_PASSWORD, ISIENV_BINDPORT etc
	flag.Parse()

	if *debugLevel {
		log.Level = logrus.DebugLevel
		log.Debug("Setting logging to debug level.")
	} else {
		log.Info("Logging set to standard level.")
		log.Level = logrus.InfoLevel
	}

	//
	isiCollectionBuildInfo.WithLabelValues(version, commit, runtime.Version()).Set(1)
	prometheus.MustRegister(isiCollectionBuildInfo)

	// gather our configuration
	config = isiconfig.GetConfig()
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	if target == "" {
		http.Error(w, "'target' parameter must be specified", 400)
		return
	}

	log.Debugf("Scraping target '%s'", target)

	registry := prometheus.NewRegistry()

	log.Info("Connecting to Isilon Cluster: " + target)
	c, err := isiclient.NewIsiClient(config.ISI.UserName, config.ISI.Password, target)
	if err != nil {
		log.Infof("Can't create Isilon Client connection : %s", err)
		isiExporterUp.WithLabelValues(target).Set(0)
		registry.MustRegister(isiExporterUp)
	} else {
		log.Debug("Isilon Cluster version is: " + c.ISIVersion)
		log.Debugf("Isilon Cluster node count: %v", c.NumNodes)

		// cluster summary info
		clusterSummaryExporter, err := collector.NewIsiClusterCollector(c, namespace)
		if err != nil {
			log.Infof("Can't create exporter : %s", err)
			isiExporterUp.WithLabelValues(target).Set(0)
			registry.MustRegister(isiExporterUp)
		} else {
			log.Debugln("Register Cluster Summary exporter")
			registry.MustRegister(clusterSummaryExporter)
		}
	}
	// Delegate http serving to Prometheus client library, which will call collector.Collect.
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func main() {
	log.Info("Starting the Isilon Exporter service...")
	log.Infof("commit: %s, build time: %s, release: %s",
		commit, date, version,
	)
	// This can go one of two ways
	// either just monitor one device or go into a query mode based on flag/env variable "multiquery"
	// to allow for multiple systems querying
	if config.Exporter.MultiQuery {
		log.Info("Running in multiquery mode...")
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<html>
            <head>
            <title>Isilon Cluster Exporter</title>
            <style>
            label{
            display:inline-block;
            width:75px;
            }
            form label {
            margin: 10px;
            }
            form input {
            margin: 10px;
            }
            </style>
            </head>
            <body>
            <h1>Isilon Cluster Exporter</h1>
            <form action="/query">
            <label>Target:</label> <input type="text" name="target" placeholder="X.X.X.X" value="1.2.3.4"><br>
            <input type="submit" value="Submit">
            </form>
            </html>`))
		})

		http.HandleFunc("/query", queryHandler)     // Endpoint to do specific cluster scrapes.
		http.Handle("/metrics", promhttp.Handler()) // endpoint for exporter stats
	} else {
		log.Info("Running in single query mode...")
		// we are only going to be watching one endpoint, so just watch that
		u, err := url.Parse(config.ISI.IsiURL)
		if err != nil {
			log.Fatalf("Issue with Isilon URL: %s\n", err)
		}
		if u.Hostname() == "" {
			log.Fatal("Hostname not defined.")
		}

		log.Info("Connecting to Isilon Cluster: " + u.Hostname())
		c, err := isiclient.NewIsiClient(config.ISI.UserName, config.ISI.Password, u.Hostname())
		if err != nil {
			log.Fatal("Unable to connect to Isilon: ", err)
		}

		log.Debug("Isilon Cluster version is: " + c.ISIVersion)
		log.Debugf("Isilon Cluster node count: %v", c.NumNodes)

		// cluster summary info
		clusterSummaryExporter, err := collector.NewIsiClusterCollector(c, namespace)
		if err != nil {
			log.Infof("Can't create exporter : %s", err)
		} else {
			log.Debugln("Register Cluster Summary exporter")
			prometheus.MustRegister(clusterSummaryExporter)
		}

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<html>
			<head><title>Dell EMC Isilon Exporter</title></head>
			<body>
			<h1>Dell EMC Isilon Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
		})

		http.Handle("/metrics", promhttp.Handler())
	}

	listenPort := fmt.Sprintf(":%v", config.Exporter.BindPort)
	log.Info("Listening on port: ", listenPort)
	log.Fatal(http.ListenAndServe(listenPort, nil))
}
