package tracer

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"go.opentelemetry.io/otel/trace"
)

const (
	targetFrameIndex = 3
	splitSize        = 2
	splitSuffixSize  = 2
	counterSize      = 2
)

func GetSpanID(ctx context.Context) string {
	spCtx := trace.SpanFromContext(ctx).SpanContext()
	if spCtx.HasSpanID() {
		return spCtx.SpanID().String()
	}

	return ""
}

func GetTraceID(ctx context.Context) string {
	spCtx := trace.SpanFromContext(ctx).SpanContext()
	if spCtx.HasTraceID() {
		return spCtx.TraceID().String()
	}

	return ""
}

func GetFunctionCallerName() string {
	fnName := getCallerName(targetFrameIndex)
	splitted := strings.Split(fnName, ".")

	if len(splitted) > splitSize {
		fnName = fmt.Sprintf("%s @ %s", strings.Join(splitted[(len(splitted)-splitSuffixSize):], "."),
			strings.Join(splitted[:(len(splitted)-splitSuffixSize)], "."),
		)
	}

	return fnName
}

func getCallerName(targetFrameIndex int) string {
	programCounters := make([]uintptr, targetFrameIndex+counterSize)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}

	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])

		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()

			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame.Function
}
