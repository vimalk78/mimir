// SPDX-License-Identifier: AGPL-3.0-only
// Provenance-includes-location: https://github.com/cortexproject/cortex/blob/master/pkg/util/log/experimental.go
// Provenance-includes-license: Apache-2.0
// Provenance-includes-copyright: The Cortex Authors.

package log

import (
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var experimentalFeaturesInUse = promauto.NewCounter(
	prometheus.CounterOpts{
		Namespace: "cortex",
		Name:      "experimental_features_used_total",
		Help:      "The number of experimental features in use.",
	},
)

// WarnExperimentalUse logs a warning and increments the experimental features metric.
func WarnExperimentalUse(feature string) {
	level.Warn(Logger).Log("msg", "experimental feature in use", "feature", feature)
	experimentalFeaturesInUse.Inc()
}
