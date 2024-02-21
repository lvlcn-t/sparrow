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

package health

import (
	"github.com/caas-team/sparrow/internal/helper"
	"github.com/prometheus/client_golang/prometheus"
)

// metrics contains the metric collectors for the Health check
type metrics struct {
	*prometheus.GaugeVec
}

// newMetrics initializes metric collectors of the health check
func newMetrics() metrics {
	return metrics{
		GaugeVec: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "sparrow_health_up",
				Help: "Health of targets",
			},
			[]string{
				"target",
			},
		),
	}
}

// GetCollectors returns all metric collectors of check
func (m *metrics) GetCollectors() []prometheus.Collector {
	return []prometheus.Collector{m}
}

func (m *metrics) Set(target string, status float64) {
	m.WithLabelValues(target).Set(status)
}

// removeObsolete removes metrics for targets not in the new config.
func (m *metrics) RemoveObsolete(oldTars, newTars []string) {
	for _, t := range oldTars {
		if !helper.SliceContains(newTars, t) {
			m.DeleteLabelValues(t)
		}
	}
}
