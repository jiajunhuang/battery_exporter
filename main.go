package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	BuildVersion = "0.1"
)

var (
	metricsPath    = flag.String("metrics_path", "/sys/class/power_supply/BAT0", "metrics path, usually it will be `/sys/class/power_supply/BAT0` or `/sys/class/power_supply/BAT1`")
	listenAt       = flag.String("listen_at", "0.0.0.0:9119", "web server listen address")
	cycleCountPath = flag.String("cycle_count_path", "/sys/class/power_supply/BAT0", "cycle count path, such as `/sys/class/power_supply/BAT0`, or `/sys/devices/platform/smapi/BAT0` for thinkpad")

	energyNow = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "battery_energy_now",
		Help: "Energy Now in mWh",
	})
	energyFull = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "battery_energy_full",
		Help: "Energy Full in mWh",
	})
	energyFullDesign = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "battery_energy_full_design",
		Help: "Energy Full in mWh By Design",
	})
	cycleCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "battery_cycle_count",
		Help: "Battery Charge Cycle Count",
	})
	MetricsMap = map[string]prometheus.Gauge{
		"energy_now":         energyNow,
		"energy_full":        energyFull,
		"energy_full_design": energyFullDesign,
	}
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(energyNow)
	prometheus.MustRegister(energyFull)
	prometheus.MustRegister(energyFullDesign)
	prometheus.MustRegister(cycleCount)
}

func readValue(path string) (float64, error) {
	m, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseFloat(strings.TrimRight(string(m), "\n"), 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func collectMetrics() {
	for k, v := range MetricsMap {
		metricFilePath := filepath.Join(*metricsPath, k)

		value, err := readValue(metricFilePath)
		if err != nil {
			log.Printf("failed to collect metric %s: %s", k, err)
			continue
		}

		v.Set(value)
	}

	// read cycle count, because ThinkPad needs special treat
	cycleCountFilePath := filepath.Join(*cycleCountPath, "cycle_count")
	value, err := readValue(cycleCountFilePath)
	if err != nil {
		log.Printf("failed to collect metric %s: %s", "cycle_count", err)
		return
	}
	cycleCount.Set(value)
}

func main() {
	flag.Parse()

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		collectMetrics()
		promhttp.HandlerFor(
			prometheus.DefaultGatherer, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError},
		).ServeHTTP(w, r)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
<head><title>Exporter v` + BuildVersion + `</title></head>
<body>
<h1>Redis Exporter ` + BuildVersion + `</h1>
<p><a href='/metrics'>Metrics</a></p>
</body>
</html>
`))
	})
	log.Printf("exporter listen at http://%s", *listenAt)
	log.Fatal(http.ListenAndServe(*listenAt, nil))
}
