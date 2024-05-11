package sinks

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

type sinkFactoryParam struct {
	parameter string
	factory   sinkFactory
}

func (s sinkFactoryParam) Invoke() (Sink, error) {
	return s.factory.CreateSink(s.parameter)
}

func (s sinkFactoryParam) String() string {
	return s.factory.Name()
}

type sinksCollectionVar []sinkFactoryParam

func (s *sinksCollectionVar) String() string {
	builder := strings.Builder{}
	hasValues := false
	for _, param := range *s {
		if hasValues {
			builder.WriteString(", ")
		}
		builder.WriteString(param.String())
		hasValues = true
	}
	return builder.String()
}

func (s *sinksCollectionVar) Set(value string) error {
	for _, factory := range allSinks {
		if factory.SupportsSink(value) {
			*s = append(*s, sinkFactoryParam{value, factory})
			return nil
		}
	}

	return fmt.Errorf("invalid sink to obtain: %s", value)
}

var sinksVar = sinksCollectionVar{sinkFactoryParam{"", stdOutSinkFactory{}}}

func init() {
	flag.Var(&sinksVar, "log-to", "Select the logging parts to write to. Can be applied multiple times")
}

func CreateSink() (io.Writer, error) {
	setup := newSetup()

	for _, factory := range sinksVar {
		s, err := factory.Invoke()
		setup.WithSink(s, err)
	}

	return setup.Build()
}
