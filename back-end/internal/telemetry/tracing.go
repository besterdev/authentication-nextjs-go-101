package telemetry

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func InitTracing(ctx context.Context, serviceName, endpoint string) (func(context.Context) error, error) {
	if endpoint == "" {
		return func(context.Context) error { return nil }, nil
	}

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("create otlp exporter: %w", err)
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	slog.Info("tracing enabled", "endpoint", endpoint)
	return tp.Shutdown, nil
}

func Tracing(serviceName string) fiber.Handler {
	tracer := otel.Tracer(serviceName)
	return func(c *fiber.Ctx) error {
		route := c.Route().Path
		if route == "" {
			route = c.Path()
		}

		ctx, span := tracer.Start(c.Context(), c.Method()+" "+route, oteltrace.WithSpanKind(oteltrace.SpanKindServer))
		defer span.End()

		c.SetUserContext(ctx)
		err := c.Next()
		span.SetAttributes(semconv.HTTPStatusCodeKey.Int(c.Response().StatusCode()))
		if err != nil {
			span.RecordError(err)
		}
		return err
	}
}
