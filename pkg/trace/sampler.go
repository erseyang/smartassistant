package trace

import (
	"encoding/binary"
	"fmt"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func CustomSamplerOpt(debug bool) sdktrace.TracerProviderOption {
	if debug {
		return sdktrace.WithSampler(CustomSampler(1))
	}
	return sdktrace.WithSampler(CustomSampler(0.01))
}

func CustomSampler(rate float64) sdktrace.Sampler {
	// debug模式下全量采集
	if rate >= 1 {
		return sdktrace.AlwaysSample()
	}

	if rate <= 0 {
		rate = 0
	}

	return &customSampler{
		traceIDUpperBound: uint64(rate * (1 << 63)),
		description:       fmt.Sprintf("TraceIDRatioBased{%g}", rate),
	}
}

type customSampler struct {
	traceIDUpperBound uint64
	description       string
}

func (cs customSampler) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
	psc := trace.SpanContextFromContext(p.ParentContext)
	if psc.TraceFlags() == trace.FlagsSampled {
		return sdktrace.SamplingResult{
			Decision:   sdktrace.RecordAndSample,
			Tracestate: psc.TraceState(),
		}
	}
	x := binary.BigEndian.Uint64(p.TraceID[0:8]) >> 1
	if x < cs.traceIDUpperBound {
		return sdktrace.SamplingResult{
			Decision:   sdktrace.RecordAndSample,
			Tracestate: psc.TraceState(),
		}
	}
	return sdktrace.SamplingResult{
		Decision:   sdktrace.Drop,
		Tracestate: psc.TraceState(),
	}
}

func (as customSampler) Description() string {
	return "customSampler"
}
