package service

import "go.opentelemetry.io/otel"

var (
	Tracer = otel.Tracer("message")
)
