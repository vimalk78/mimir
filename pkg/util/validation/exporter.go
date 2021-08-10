// SPDX-License-Identifier: AGPL-3.0-only
// Provenance-includes-location: https://github.com/cortexproject/cortex/blob/master/pkg/util/validation/exporter.go
// Provenance-includes-license: Apache-2.0
// Provenance-includes-copyright: The Cortex Authors.

package validation

import (
	"github.com/prometheus/client_golang/prometheus"
)

// OverridesExporter exposes per-tenant resource limit overrides as Prometheus metrics
type OverridesExporter struct {
	tenantLimits TenantLimits
	description  *prometheus.Desc
}

// NewOverridesExporter creates an OverridesExporter that reads updates to per-tenant
// limits using the provided function.
func NewOverridesExporter(tenantLimits TenantLimits) *OverridesExporter {
	return &OverridesExporter{
		tenantLimits: tenantLimits,
		description: prometheus.NewDesc(
			"cortex_overrides",
			"Resource limit overrides applied to tenants",
			[]string{"limit_name", "user"},
			nil,
		),
	}
}

func (oe *OverridesExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- oe.description
}

func (oe *OverridesExporter) Collect(ch chan<- prometheus.Metric) {
	allLimits := oe.tenantLimits.AllByUserID()
	for tenant, limits := range allLimits {
		// Write path limits
		ch <- prometheus.MustNewConstMetric(oe.description, prometheus.GaugeValue, limits.IngestionRate, "ingestion_rate", tenant)
		ch <- prometheus.MustNewConstMetric(oe.description, prometheus.GaugeValue, float64(limits.IngestionBurstSize), "ingestion_burst_size", tenant)
		ch <- prometheus.MustNewConstMetric(oe.description, prometheus.GaugeValue, float64(limits.MaxLocalSeriesPerUser), "max_local_series_per_user", tenant)
		ch <- prometheus.MustNewConstMetric(oe.description, prometheus.GaugeValue, float64(limits.MaxLocalSeriesPerMetric), "max_local_series_per_metric", tenant)
		ch <- prometheus.MustNewConstMetric(oe.description, prometheus.GaugeValue, float64(limits.MaxGlobalSeriesPerUser), "max_global_series_per_user", tenant)
		ch <- prometheus.MustNewConstMetric(oe.description, prometheus.GaugeValue, float64(limits.MaxGlobalSeriesPerMetric), "max_global_series_per_metric", tenant)

		// Read path limits
		ch <- prometheus.MustNewConstMetric(oe.description, prometheus.GaugeValue, float64(limits.MaxFetchedSeriesPerQuery), "max_fetched_series_per_query", tenant)
		ch <- prometheus.MustNewConstMetric(oe.description, prometheus.GaugeValue, float64(limits.MaxFetchedChunkBytesPerQuery), "max_fetched_chunk_bytes_per_query", tenant)
		ch <- prometheus.MustNewConstMetric(oe.description, prometheus.GaugeValue, float64(limits.MaxSeriesPerQuery), "max_series_per_query", tenant)
		ch <- prometheus.MustNewConstMetric(oe.description, prometheus.GaugeValue, float64(limits.MaxSamplesPerQuery), "max_samples_per_query", tenant)

		// Ruler limits
		ch <- prometheus.MustNewConstMetric(oe.description, prometheus.GaugeValue, float64(limits.RulerMaxRulesPerRuleGroup), "ruler_max_rules_per_rule_group", tenant)
		ch <- prometheus.MustNewConstMetric(oe.description, prometheus.GaugeValue, float64(limits.RulerMaxRuleGroupsPerTenant), "ruler_max_rule_groups_per_tenant", tenant)
	}
}
