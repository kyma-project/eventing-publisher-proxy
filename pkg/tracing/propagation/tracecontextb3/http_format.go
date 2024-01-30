package tracecontextb3

import (
	"github.com/kyma-project/eventing-publisher-proxy/pkg/tracing/propagation"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	ocpropagation "go.opencensus.io/trace/propagation"
)

// TraceContextEgress returns a propagation.HTTPFormat that reads both TraceContext and B3 tracing
// formats, preferring TraceContext. It always writes TraceContext format exclusively.
func TraceContextEgress() *propagation.HTTPFormatSequence {
	return &propagation.HTTPFormatSequence{
		Ingress: []ocpropagation.HTTPFormat{
			&tracecontext.HTTPFormat{},
			&b3.HTTPFormat{},
		},
		Egress: []ocpropagation.HTTPFormat{
			&tracecontext.HTTPFormat{},
		},
	}
}
