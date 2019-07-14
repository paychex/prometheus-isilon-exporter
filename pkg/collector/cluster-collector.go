package collector

import (
	"fmt"
	"strconv"
	"time"

	"github.com/paychex/prometheus-isilon-exporter/pkg/isiclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/tidwall/gjson"
)

var (
	clusterSummaryCPU = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "cpu_usage"),
		"The percentage CPU utilization.",
		[]string{"clustername"}, nil,
	)
	clusterSummaryFTPthroughput = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "ftp_throughput"),
		"The total throughput (in bytes/sec) for FTP operations.",
		[]string{"clustername"}, nil,
	)
	clusterSummaryHTTPthroughput = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "http_throughput"),
		"The total throughput (in bytes/sec) for HTTP operations.",
		[]string{"clustername"}, nil,
	)
	clusterSummaryHDFSthroughput = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "hdfs_throughput"),
		"The total throughput (in bytes/sec) for HDFS operations.",
		[]string{"clustername"}, nil,
	)
	clusterSummaryiSCSIthroughput = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "iscsi_throughput"),
		"The total throughput (in bytes/sec) for iSCSI operations.",
		[]string{"clustername"}, nil,
	)
	clusterSummarySMBthroughput = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "smb_throughput"),
		"The total throughput (in bytes/sec) for SMB operations.",
		[]string{"clustername"}, nil,
	)
	clusterSummaryNFSthroughput = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "nfs_throughput"),
		"The total throughput (in bytes/sec) for NFS operations.",
		[]string{"clustername"}, nil,
	)
	clusterSummaryNetOutthroughput = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "net_out_throughput"),
		"Outgoing network traffic (in bytes/sec) for all operations.",
		[]string{"clustername"}, nil,
	)
	clusterSummaryNetInthroughput = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "net_in_throughput"),
		"Incoming network traffic (in bytes/sec) for all operations.",
		[]string{"clustername"}, nil,
	)
	clusterSummaryNetTotalthroughput = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "net_total_throughput"),
		"The total throughput (in bytes/sec) for all protocols listed.",
		[]string{"clustername"}, nil,
	)
	clusterSummaryDiskInthroughput = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "disk_in_throughput"),
		"Traffic to disk (in bytes/sec).",
		[]string{"clustername"}, nil,
	)
	clusterSummaryDiskOutthroughput = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "disk_out_throughput"),
		"Traffic from disk (in bytes/sec).",
		[]string{"clustername"}, nil,
	)
	clusterIFSBytesAvail = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "ifs_bytes_avail"),
		"Traffic from disk (in bytes/sec).",
		[]string{"clustername"}, nil,
	)
	clusterIFSBytesFree = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "ifs_bytes_free"),
		"Traffic from disk (in bytes/sec).",
		[]string{"clustername"}, nil,
	)
	clusterIFSBytesTotal = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "ifs_bytes_total"),
		"Traffic from disk (in bytes/sec).",
		[]string{"clustername"}, nil,
	)
	clusterSSDIFSBytesAvail = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "ifs_ssd_bytes_avail"),
		"Traffic from disk (in bytes/sec).",
		[]string{"clustername"}, nil,
	)
	clusterSSDIFSBytesFree = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "ifs_ssd_bytes_free"),
		"Traffic from disk (in bytes/sec).",
		[]string{"clustername"}, nil,
	)
	clusterSSDIFSBytesTotal = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "ifs_ssd_bytes_total"),
		"Traffic from disk (in bytes/sec).",
		[]string{"clustername"}, nil,
	)
	alertsnumcritical = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "alerts_critical"),
		"Number of current critical alerts for the cluster",
		[]string{"clustername"}, nil,
	)
	alertsnumerror = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "alerts_error"),
		"Number of current error alerts for the cluster",
		[]string{"clustername"}, nil,
	)
	alertsnuminfo = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "alerts_info"),
		"Number of current info alerts for the cluster",
		[]string{"clustername"}, nil,
	)
	alertsnumwarning = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "alerts_warning"),
		"Number of current warning alerts for the cluster",
		[]string{"clustername"}, nil,
	)
	nodeDiskBusy = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "node", "disk_busy"),
		"The percentage of time the drive was busy.",
		[]string{"clustername", "drive_id", "type"}, nil,
	)
	nodeDiskAccessLatency = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "node", "disk_access_latency"),
		"The average operation latency.",
		[]string{"clustername", "drive_id", "type"}, nil,
	)
	nodeDiskBytesIn = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "node", "disk_bytes_in"),
		"The rate of bytes written.",
		[]string{"clustername", "drive_id", "type"}, nil,
	)
	nodeDiskBytesOut = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "node", "disk_bytes_out"),
		"The rate of bytes read.",
		[]string{"clustername", "drive_id", "type"}, nil,
	)
	pathHardQuota = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "hard_quota"),
		"HardQuota of a path bytes",
		[]string{"clustername", "path"}, nil,
	)
	pathAdvisoryQuota = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "advisory_quota"),
		"Advisory Quota of a path bytes",
		[]string{"clustername", "path"}, nil,
	)
	pathLogicalUsed = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "logical_used"),
		"Used data w/o overhead of a path bytes",
		[]string{"clustername", "path"}, nil,
	)
	pathPhysicalUsed = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "physical_used"),
		"Used Data w/overhead of a path bytes",
		[]string{"clustername", "path"}, nil,
	)
	exporterUp = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "exporter", "up"),
		"Indicates if scrape was succesful or not.",
		[]string{"clustername"}, nil,
	)
	isiClusterInfo = prometheus.NewDesc(
		prometheus.BuildFQName("emcisi", "cluster", "version"),
		"A metric with a constant '1' value labeled by version, and nodecount",
		[]string{"version", "nodecount", "clustername"}, nil,
	)
	isiCollectionDuration = prometheus.NewDesc(
		"emcisi_collection_duration_seconds",
		"Duration of collections by the EMC Isilon exporter",
		[]string{"clustername"}, nil,
	)
)

// A IsiClusterCollector implements the prometheus.Collector.
type IsiClusterCollector struct {
	isiClient *isiclient.ISIClient
	namespace string
}

// NewIsiClusterCollector returns an initialized Isilon Cluster Collector.
func NewIsiClusterCollector(emcisi *isiclient.ISIClient, namespace string) (*IsiClusterCollector, error) {

	log.Debugln("Init exporter")
	return &IsiClusterCollector{
		isiClient: emcisi,
		namespace: namespace,
	}, nil
}

// Collect fetches the stats from the Isilon cluster and delivers them
// as Prometheus metrics.
// It implements prometheus.Collector.
func (e *IsiClusterCollector) Collect(ch chan<- prometheus.Metric) {
	log.Debugln("Isilon Cluster collect starting")
	start := time.Now()

	if e.isiClient == nil {
		log.Errorf("Isilon client not configured.")
		duration := float64(time.Since(start).Seconds())
		ch <- prometheus.MustNewConstMetric(isiCollectionDuration, prometheus.GaugeValue, duration, e.isiClient.ClusterName)
		ch <- prometheus.MustNewConstMetric(exporterUp, prometheus.GaugeValue, 0, e.isiClient.ClusterName)
		return
	}

	ch <- prometheus.MustNewConstMetric(isiClusterInfo, prometheus.GaugeValue, 1, e.isiClient.ISIVersion, strconv.FormatInt(e.isiClient.NumNodes, 10), e.isiClient.ClusterName)

	// Get base system summary status
	reqStatusURL := "https://" + e.isiClient.ClusterAddress + ":8080/platform/3/statistics/summary/system"
	s, err := e.isiClient.CallIsiAPI(reqStatusURL, 1)

	if err != nil {
		duration := float64(time.Since(start).Seconds())
		ch <- prometheus.MustNewConstMetric(isiCollectionDuration, prometheus.GaugeValue, duration, e.isiClient.ClusterName)
		ch <- prometheus.MustNewConstMetric(exporterUp, prometheus.GaugeValue, 0, e.isiClient.ClusterName)
		return
	}

	ch <- prometheus.MustNewConstMetric(clusterSummaryCPU, prometheus.GaugeValue, gjson.Get(s, "system.0.cpu").Float(), e.isiClient.ClusterName)
	ch <- prometheus.MustNewConstMetric(clusterSummaryFTPthroughput, prometheus.GaugeValue, gjson.Get(s, "system.0.ftp").Float(), e.isiClient.ClusterName)
	ch <- prometheus.MustNewConstMetric(clusterSummaryHTTPthroughput, prometheus.GaugeValue, gjson.Get(s, "system.0.http").Float(), e.isiClient.ClusterName)
	ch <- prometheus.MustNewConstMetric(clusterSummaryHDFSthroughput, prometheus.GaugeValue, gjson.Get(s, "system.0.hdfs").Float(), e.isiClient.ClusterName)
	ch <- prometheus.MustNewConstMetric(clusterSummaryiSCSIthroughput, prometheus.GaugeValue, gjson.Get(s, "system.0.iscsi").Float(), e.isiClient.ClusterName)
	ch <- prometheus.MustNewConstMetric(clusterSummarySMBthroughput, prometheus.GaugeValue, gjson.Get(s, "system.0.smb").Float(), e.isiClient.ClusterName)
	ch <- prometheus.MustNewConstMetric(clusterSummaryNFSthroughput, prometheus.GaugeValue, gjson.Get(s, "system.0.nfs").Float(), e.isiClient.ClusterName)
	ch <- prometheus.MustNewConstMetric(clusterSummaryNetInthroughput, prometheus.GaugeValue, gjson.Get(s, "system.0.net_in").Float(), e.isiClient.ClusterName)
	ch <- prometheus.MustNewConstMetric(clusterSummaryNetOutthroughput, prometheus.GaugeValue, gjson.Get(s, "system.0.net_out").Float(), e.isiClient.ClusterName)
	ch <- prometheus.MustNewConstMetric(clusterSummaryDiskInthroughput, prometheus.GaugeValue, gjson.Get(s, "system.0.disk_in").Float(), e.isiClient.ClusterName)
	ch <- prometheus.MustNewConstMetric(clusterSummaryDiskOutthroughput, prometheus.GaugeValue, gjson.Get(s, "system.0.disk_out").Float(), e.isiClient.ClusterName)
	ch <- prometheus.MustNewConstMetric(clusterSummaryNetTotalthroughput, prometheus.GaugeValue, gjson.Get(s, "system.0.total").Float(), e.isiClient.ClusterName)

	// Get cluster space information
	reqStatusURL = "https://" + e.isiClient.ClusterAddress + ":8080/platform/1/statistics/current?key=ifs.bytes.total&key=ifs.ssd.bytes.total&key=ifs.bytes.free&key=ifs.ssd.bytes.free&key=ifs.bytes.avail&key=ifs.ssd.bytes.avail&devid=all"
	s, err = e.isiClient.CallIsiAPI(reqStatusURL, 1)
	if err != nil {
		duration := float64(time.Since(start).Seconds())
		ch <- prometheus.MustNewConstMetric(isiCollectionDuration, prometheus.GaugeValue, duration, e.isiClient.ClusterName)
		ch <- prometheus.MustNewConstMetric(exporterUp, prometheus.GaugeValue, 0, e.isiClient.ClusterName)
		return
	}
	result := gjson.Get(s, "stats")
	result.ForEach(func(key, value gjson.Result) bool {
		switch gjson.Get(value.String(), "key").String() {
		case "ifs.bytes.avail":
			ch <- prometheus.MustNewConstMetric(clusterIFSBytesAvail, prometheus.GaugeValue, gjson.Get(value.String(), "value").Float(), e.isiClient.ClusterName)
		case "ifs.bytes.free":
			ch <- prometheus.MustNewConstMetric(clusterIFSBytesFree, prometheus.GaugeValue, gjson.Get(value.String(), "value").Float(), e.isiClient.ClusterName)
		case "ifs.bytes.total":
			ch <- prometheus.MustNewConstMetric(clusterIFSBytesTotal, prometheus.GaugeValue, gjson.Get(value.String(), "value").Float(), e.isiClient.ClusterName)
		case "ifs.ssd.bytes.avail":
			ch <- prometheus.MustNewConstMetric(clusterSSDIFSBytesAvail, prometheus.GaugeValue, gjson.Get(value.String(), "value").Float(), e.isiClient.ClusterName)
		case "ifs.ssd.bytes.free":
			ch <- prometheus.MustNewConstMetric(clusterSSDIFSBytesFree, prometheus.GaugeValue, gjson.Get(value.String(), "value").Float(), e.isiClient.ClusterName)
		case "ifs.ssd.bytes.total":
			ch <- prometheus.MustNewConstMetric(clusterSSDIFSBytesTotal, prometheus.GaugeValue, gjson.Get(value.String(), "value").Float(), e.isiClient.ClusterName)
		default:
			fmt.Println("Got something else")
		}
		return true
	})

	// Retrieve individual drive stats
	reqStatusURL = "https://" + e.isiClient.ClusterAddress + ":8080/platform/3/statistics/summary/drive"
	s, err = e.isiClient.CallIsiAPI(reqStatusURL, 1)
	if err != nil {
		duration := float64(time.Since(start).Seconds())
		ch <- prometheus.MustNewConstMetric(isiCollectionDuration, prometheus.GaugeValue, duration, e.isiClient.ClusterName)
		ch <- prometheus.MustNewConstMetric(exporterUp, prometheus.GaugeValue, 0, e.isiClient.ClusterName)
		return
	}
	result = gjson.Get(s, "drive")
	result.ForEach(func(key, value gjson.Result) bool {
		// Cuz I am getting this info multiple times
		did := gjson.Get(value.String(), "drive_id").String()
		dtype := gjson.Get(value.String(), "type").String()
		ch <- prometheus.MustNewConstMetric(nodeDiskBusy, prometheus.GaugeValue, gjson.Get(value.String(), "busy").Float(), e.isiClient.ClusterName, did, dtype)
		ch <- prometheus.MustNewConstMetric(nodeDiskAccessLatency, prometheus.GaugeValue, gjson.Get(value.String(), "access_latency").Float(), e.isiClient.ClusterName, did, dtype)
		ch <- prometheus.MustNewConstMetric(nodeDiskBytesIn, prometheus.GaugeValue, gjson.Get(value.String(), "bytes_in").Float(), e.isiClient.ClusterName, did, dtype)
		ch <- prometheus.MustNewConstMetric(nodeDiskBytesOut, prometheus.GaugeValue, gjson.Get(value.String(), "bytes_out").Float(), e.isiClient.ClusterName, did, dtype)
		return true
	})

	// get count of errors in "information", "warning" and "error" states that are not resolved
	reqStatusURL = "https://" + e.isiClient.ClusterAddress + ":8080/platform/3/event/eventgroup-occurrences?resolved=false&ignore=false"
	s, err = e.isiClient.CallIsiAPI(reqStatusURL, 1)
	if err != nil {
		duration := float64(time.Since(start).Seconds())
		ch <- prometheus.MustNewConstMetric(isiCollectionDuration, prometheus.GaugeValue, duration, e.isiClient.ClusterName)
		ch <- prometheus.MustNewConstMetric(exporterUp, prometheus.GaugeValue, 0, e.isiClient.ClusterName)
		return
	}
	result = gjson.Get(s, `eventgroups.#[severity=="warning"]#`)
	ch <- prometheus.MustNewConstMetric(alertsnumwarning, prometheus.GaugeValue, arrayCount(result), e.isiClient.ClusterName)
	result = gjson.Get(s, `eventgroups.#[severity=="information"]#`)
	ch <- prometheus.MustNewConstMetric(alertsnuminfo, prometheus.GaugeValue, arrayCount(result), e.isiClient.ClusterName)
	result = gjson.Get(s, `eventgroups.#[severity=="error"]#`)
	ch <- prometheus.MustNewConstMetric(alertsnumerror, prometheus.GaugeValue, arrayCount(result), e.isiClient.ClusterName)
	result = gjson.Get(s, `eventgroups.#[severity=="error"]#`)
	ch <- prometheus.MustNewConstMetric(alertsnumcritical, prometheus.GaugeValue, arrayCount(result), e.isiClient.ClusterName)

	//Quota Collection
	reqStatusURL = "https://" + e.isiClient.ClusterAddress + ":8080/platform/1/quota/quotas"
	s, err = e.isiClient.CallIsiAPI(reqStatusURL, 1)
	if err != nil {
		duration := float64(time.Since(start).Seconds())
		ch <- prometheus.MustNewConstMetric(isiCollectionDuration, prometheus.GaugeValue, duration, e.isiClient.ClusterName)
		ch <- prometheus.MustNewConstMetric(exporterUp, prometheus.GaugeValue, 0, e.isiClient.ClusterName)
		return
	}
	result = gjson.Get(s, "quotas")
	result.ForEach(func(key, value gjson.Result) bool {
		path := gjson.Get(value.String(), "path").String()
		ch <- prometheus.MustNewConstMetric(pathHardQuota, prometheus.GaugeValue, gjson.Get(value.String(), "thresholds.hard").Float(), e.isiClient.ClusterName, path)
		ch <- prometheus.MustNewConstMetric(pathAdvisoryQuota, prometheus.GaugeValue, gjson.Get(value.String(), "thresholds.advisory").Float(), e.isiClient.ClusterName, path)
		ch <- prometheus.MustNewConstMetric(pathLogicalUsed, prometheus.GaugeValue, gjson.Get(value.String(), "usage.logical").Float(), e.isiClient.ClusterName, path)
		ch <- prometheus.MustNewConstMetric(pathPhysicalUsed, prometheus.GaugeValue, gjson.Get(value.String(), "usage.physical").Float(), e.isiClient.ClusterName, path)
		return true
	})
	duration := float64(time.Since(start).Seconds())
	ch <- prometheus.MustNewConstMetric(isiCollectionDuration, prometheus.GaugeValue, duration, e.isiClient.ClusterName)
	ch <- prometheus.MustNewConstMetric(exporterUp, prometheus.GaugeValue, 1, e.isiClient.ClusterName)
	log.Debugf("Scrape of target '%s' took %f seconds", e.isiClient.ClusterName, duration)
	log.Infoln("Cluster exporter finished")
}

func arrayCount(r gjson.Result) (count float64) {
	r.ForEach(func(key, value gjson.Result) bool {
		count++
		return true
	})
	return
}

// Describe describes the metrics exported from this collector.
func (e *IsiClusterCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- exporterUp
	ch <- isiClusterInfo
	ch <- clusterSummaryCPU
	ch <- isiCollectionDuration
	ch <- clusterSummaryFTPthroughput
	ch <- clusterSummaryHTTPthroughput
	ch <- clusterSummaryHDFSthroughput
	ch <- clusterSummaryiSCSIthroughput
	ch <- clusterSummarySMBthroughput
	ch <- clusterSummaryNFSthroughput
	ch <- clusterSummaryNetInthroughput
	ch <- clusterSummaryNetOutthroughput
	ch <- clusterSummaryNetTotalthroughput
	ch <- clusterSummaryDiskInthroughput
	ch <- clusterSummaryDiskOutthroughput
	ch <- clusterIFSBytesAvail
	ch <- clusterIFSBytesFree
	ch <- clusterIFSBytesTotal
	ch <- clusterSSDIFSBytesAvail
	ch <- clusterSSDIFSBytesFree
	ch <- clusterSSDIFSBytesTotal
	ch <- alertsnumcritical
	ch <- alertsnumerror
	ch <- alertsnuminfo
	ch <- alertsnumwarning
	ch <- nodeDiskBusy
	ch <- nodeDiskAccessLatency
	ch <- nodeDiskBytesIn
	ch <- nodeDiskBytesOut
	ch <- pathHardQuota
	ch <- pathAdvisoryQuota
	ch <- pathLogicalUsed
	ch <- pathPhysicalUsed
}
