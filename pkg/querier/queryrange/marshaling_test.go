// SPDX-License-Identifier: AGPL-3.0-only
// Provenance-includes-location: https://github.com/cortexproject/cortex/blob/master/pkg/querier/queryrange/marshaling_test.go
// Provenance-includes-license: Apache-2.0
// Provenance-includes-copyright: The Cortex Authors.

package queryrange

import (
	"bytes"
	"context"
	"io/ioutil"
	"math/rand"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/grafana/mimir/pkg/cortexpb"
)

func BenchmarkPrometheusCodec_DecodeResponse(b *testing.B) {
	const (
		numSeries           = 1000
		numSamplesPerSeries = 1000
	)

	// Generate a mocked response and marshal it.
	res := mockPrometheusResponse(numSeries, numSamplesPerSeries)
	encodedRes, err := json.Marshal(res)
	require.NoError(b, err)
	b.Log("test prometheus response size:", len(encodedRes))

	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		_, err := PrometheusCodec.DecodeResponse(context.Background(), &http.Response{
			StatusCode:    200,
			Body:          ioutil.NopCloser(bytes.NewReader(encodedRes)),
			ContentLength: int64(len(encodedRes)),
		}, nil)
		require.NoError(b, err)
	}
}

func BenchmarkPrometheusCodec_EncodeResponse(b *testing.B) {
	const (
		numSeries           = 1000
		numSamplesPerSeries = 1000
	)

	// Generate a mocked response and marshal it.
	res := mockPrometheusResponse(numSeries, numSamplesPerSeries)

	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		_, err := PrometheusCodec.EncodeResponse(context.Background(), res)
		require.NoError(b, err)
	}
}

func mockPrometheusResponse(numSeries, numSamplesPerSeries int) *PrometheusResponse {
	stream := make([]SampleStream, numSeries)
	for s := 0; s < numSeries; s++ {
		// Generate random samples.
		samples := make([]cortexpb.Sample, numSamplesPerSeries)
		for i := 0; i < numSamplesPerSeries; i++ {
			samples[i] = cortexpb.Sample{
				Value:       rand.Float64(),
				TimestampMs: int64(i),
			}
		}

		// Generate random labels.
		lbls := make([]cortexpb.LabelAdapter, 10)
		for i := range lbls {
			lbls[i].Name = "a_medium_size_label_name"
			lbls[i].Value = "a_medium_size_label_value_that_is_used_to_benchmark_marshalling"
		}

		stream[s] = SampleStream{
			Labels:  lbls,
			Samples: samples,
		}
	}

	return &PrometheusResponse{
		Status: "success",
		Data: PrometheusData{
			ResultType: "vector",
			Result:     stream,
		},
	}
}
