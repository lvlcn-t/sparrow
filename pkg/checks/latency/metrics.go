// sparrow
// (C) 2024, Deutsche Telekom IT GmbH
//
// Deutsche Telekom IT GmbH and all other contributors /
// copyright owners license this file to you under the Apache
// License, Version 2.0 (the "License"); you may not use this
// file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package latency

import (
	"github.com/caas-team/sparrow/internal/helper"
	"github.com/prometheus/client_golang/prometheus"
)

// metrics defines the metric collectors of the latency check
type metrics struct {
	duration  *prometheus.GaugeVec
	count     *prometheus.CounterVec
	histogram *prometheus.HistogramVec
}

// newMetrics initializes metric collectors of the latency check
func newMetrics() metrics {
	return metrics{
		duration: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "sparrow_latency_duration_seconds",
				Help: "Latency with status information of targets",
			},
			[]string{
				"target",
				"status",
			},
		),
		count: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "sparrow_latency_count",
				Help: "Count of latency checks done",
			},
			[]string{
				"target",
			},
		),
		histogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "sparrow_latency_duration",
				Help: "Latency of targets in seconds",
			},
			[]string{
				"target",
			},
		),
	}
}

// GetCollectors returns all metric collectors
func (m *metrics) GetCollectors() []prometheus.Collector {
	return []prometheus.Collector{
		m.duration,
		m.count,
		m.histogram,
	}
}

func (m *metrics) Set(target, status string, total float64) {
	m.duration.WithLabelValues(target, status).Set(total)
	m.histogram.WithLabelValues(target).Observe(total)
	m.count.WithLabelValues(target).Inc()
}

// RemoveObsolete removes metrics for targets not in the new config.
func (m *metrics) RemoveObsolete(oldTars, newTars []string) {
	for _, t := range oldTars {
		if !helper.SliceContains(newTars, t) {
			m.duration.DeleteLabelValues(t)
			m.histogram.DeleteLabelValues(t)
			m.count.DeleteLabelValues(t)
		}
	}
}
