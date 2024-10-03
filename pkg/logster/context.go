package logster

import (
	"context"
	"github.com/google/uuid"
)

type ctxKey int

const (
	fieldsKey      ctxKey = 0
	spanKey               = "span"
	traceKey              = "trace"
	parentTraceKey        = "parent_trace"
)

func Log(ctx context.Context, logger Logger) Logger {
	if f, ok := FieldsFromContext(ctx); ok {
		return logger.WithFields(f)
	}
	return logger
}

func Span(ctx context.Context, spanPart string) context.Context {
	return span(ctx, spanPart, false)
}

// SpanChild copies passed context's behaviour but rewrites trace with new value
// use it when you need to keep same cancel for graceful shutdown but need differ trace
func SpanChild(ctx context.Context, spanPart string) context.Context {
	return span(ctx, spanPart, true)
}

func span(ctx context.Context, span string, updateTrace bool) context.Context {
	if span == "" {
		return ctx
	}

	fields, _ := FieldsFromContext(ctx)

	if parentSpan, ok := fields[spanKey].(string); ok && parentSpan != "" {
		span = parentSpan + ":" + span
	}
	fields[spanKey] = span

	if _, ok := fields[traceKey]; !ok {
		// case than SpanChild calls first time. Creates only one trace
		// Next SpanChild would create parentTrace and new trace
		fields[traceKey] = uuid.NewString() // possible panic btw
		return context.WithValue(ctx, fieldsKey, fields)
	}

	if updateTrace {
		parentTrace := fields[traceKey]
		fields[parentTraceKey] = parentTrace
		fields[traceKey] = uuid.NewString() // possible panic btw
	}

	return context.WithValue(ctx, fieldsKey, fields)
}

func FieldsFromContext(ctx context.Context) (Fields, bool) {
	fields, ok := ctx.Value(fieldsKey).(Fields)
	if ok {
		fieldsCopy := make(map[string]interface{}, len(fields))

		for k, v := range fields {
			fieldsCopy[k] = v
		}

		return fieldsCopy, true
	}

	return Fields{}, false
}

func ContextWithFields(ctx context.Context, fields Fields) context.Context {
	var putFields = fields
	// enrich with fields from parent context
	if oldFields, ok := FieldsFromContext(ctx); ok {
		for k, v := range fields {
			// NB! new fields overwrites old ones
			oldFields[k] = v
		}
		putFields = oldFields
	}
	if _, ok := putFields[traceKey]; !ok {
		putFields[traceKey] = uuid.NewString() // possible panic btw
	}
	return context.WithValue(ctx, fieldsKey, putFields)
}
