# Prometheus exporter for EMC Isilon
[![Build Status](https://api.travis-ci.com/paychex/prometheus-isilon-exporter.svg?branch=master)](https://travis-ci.com/paychex/prometheus-isilon-exporter/builds)
[![Go Report Card](https://goreportcard.com/badge/github.com/paychex/prometheus-isilon-exporter)](https://goreportcard.com/report/github.com/paychex/prometheus-isilon-exporter)

This exporter collects performance and usage stats from Dell/EMC Isilon cluster running version 8.x and above OneFS code and makes it available for Prometheus to scrape.  It is not recommended that you run this tool on the Isilon Cluster node(s), instead it should be run on a separate machine.  The application can be configured to monitor just one cluster, or can be configured to query multiple Isilon clusters.  See configuration options below for how to use this tool.

## Usage

| Flag      | Description                                                                                                                                           | Default Value | Env Name         |
|-----------|-------------------------------------------------------------------------------------------------------------------------------------------------------|---------------|------------------|
| url       | Base URL of the Isilon management interface.  Normally something like https://myisilon.internal.com:8080.  This is ignored when using the multi flag. | none          | ISIENV_URL       |
| username  | Username with which to connect to the Isilon API                                                                                                      | none          | ISIENV_USERNAME  |
| password  | Password with which to connect to the Isilon API                                                                                                      | none          | ISIENV_PASSWORD  |
| bind_port | Port to bind the exporter endpoint to                                                                                                                 | 9437          | ISIENV_BIND_PORT |
| multi     | Enable multi query endpoint                                                                                                                           | false         | ISIENV_MULTI     |

### Running in multi-query mode

While normally one runs one exporter per device, there are times where running one exporter for multiple Isilon devices may make sense.  This setup works similar to the [SNMP exporter](https://github.com/prometheus/snmp_exporter).  Note that you will need to configure each Isilon device to use the same username and password for this to work properly.

When configuring Prometheus to scrape in this manner use the following Prometheus config snippet:

````YAML
scrape_configs:
  - job_name: 'isilon'
    static_configs:
      - targets:
        - 192.168.1.2  # Isilon device.
        - 192.168.1.3  # Isilon device 2
        - 192.168.2.2  # Isilon device 3, etc
    metrics_path: /query
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:9437  # The isilon exporter's real hostname:port running in "multi-query" mode
  - job_name: 'isilon-exporter-stats' # gathers the exporter application process stats if you want this sort of information
    static_configs:
      - targets: 127.0.0.1:9437
````

## Exported Metrics

### Isilon

````
# HELP emcisi_cluster_alerts_critical Number of current critical alerts for the cluster
# TYPE emcisi_cluster_alerts_critical gauge
# HELP emcisi_cluster_alerts_error Number of current error alerts for the cluster
# TYPE emcisi_cluster_alerts_error gauge
# HELP emcisi_cluster_alerts_info Number of current info alerts for the cluster
# TYPE emcisi_cluster_alerts_info gauge
# HELP emcisi_cluster_alerts_warning Number of current warning alerts for the cluster
# TYPE emcisi_cluster_alerts_warning gauge
# HELP emcisi_cluster_cpu_usage The percentage CPU utilization.
# TYPE emcisi_cluster_cpu_usage gauge
# HELP emcisi_cluster_disk_in_throughput Traffic to disk (in bytes/sec).
# TYPE emcisi_cluster_disk_in_throughput gauge
# HELP emcisi_cluster_disk_out_throughput Traffic from disk (in bytes/sec).
# TYPE emcisi_cluster_disk_out_throughput gauge
# HELP emcisi_cluster_ftp_throughput The total throughput (in bytes/sec) for FTP operations.
# TYPE emcisi_cluster_ftp_throughput gauge
# HELP emcisi_cluster_hdfs_throughput The total throughput (in bytes/sec) for HDFS operations.
# TYPE emcisi_cluster_hdfs_throughput gauge
# HELP emcisi_cluster_http_throughput The total throughput (in bytes/sec) for HTTP operations.
# TYPE emcisi_cluster_http_throughput gauge
# HELP emcisi_cluster_iscsi_throughput The total throughput (in bytes/sec) for iSCSI operations.
# TYPE emcisi_cluster_iscsi_throughput gauge
# HELP emcisi_cluster_net_in_throughput Incoming network traffic (in bytes/sec) for all operations.
# TYPE emcisi_cluster_net_in_throughput gauge
# HELP emcisi_cluster_net_out_throughput Outgoing network traffic (in bytes/sec) for all operations.
# TYPE emcisi_cluster_net_out_throughput gauge
# HELP emcisi_cluster_net_total_throughput The total throughput (in bytes/sec) for all protocols listed.
# TYPE emcisi_cluster_net_total_throughput gauge
# HELP emcisi_cluster_nfs_throughput The total throughput (in bytes/sec) for NFS operations.
# TYPE emcisi_cluster_nfs_throughput gauge
# HELP emcisi_cluster_smb_throughput The total throughput (in bytes/sec) for SMB operations.
# TYPE emcisi_cluster_smb_throughput gauge
# HELP emcisi_cluster_version A metric with a constant '1' value labeled by version, and nodecount
# TYPE emcisi_cluster_version gauge
````

## Building

This exporter can run on any go supported platform.  As of version 1.2 we have moved to using Go 1.11 and higher. Testing is done with Go 1.12 but go 1.11 should work for anyone using it.

To build run:
`go build`

You can also run:
`go get github.com/paychex/prometheus-isilon-exporter`

## Refrences

- https://www.emc.com/collateral/TechnicalDocument/docu66301.pdf
- https://thesanguy.com/2017/06/30/custom-reporting-with-isilon-onefs-api-calls/

## Author

This exporter was originally written by [Mark DeNeve](https://github.com/xphyr)