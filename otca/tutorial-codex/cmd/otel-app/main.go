package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/example/otca-lab/internal/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/trace"
)

const (
	defaultHTTPPort = "8080"
)

var (
	inventoryRequestCounter   instrument.Int64Counter
	inventoryLatencyHistogram instrument.Float64Histogram
)

func main() {
	rand.Seed(time.Now().UnixNano())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	shutdown, err := telemetry.InitProviders(ctx)
	if err != nil {
		log.Fatalf("failed to initialise telemetry: %v", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := shutdown(shutdownCtx); err != nil {
			log.Printf("telemetry shutdown error: %v", err)
		}
	}()

	if err := configureInstruments(ctx); err != nil {
		log.Fatalf("failed to configure instruments: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", otelhttp.NewHandler(http.HandlerFunc(rootHandler), "root"))
	mux.Handle("/healthz", http.HandlerFunc(healthHandler))
	mux.Handle("/inventory", otelhttp.NewHandler(http.HandlerFunc(inventoryHandler), "inventory"))

	port := envOr("APP_PORT", defaultHTTPPort)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		log.Printf("server listening on :%s", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down HTTP server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
}

func configureInstruments(ctx context.Context) error {
	meter := otel.Meter("otca-lab/app")

	var err error
	inventoryRequestCounter, err = meter.Int64Counter(
		"lab.inventory.requests",
		instrument.WithUnit("1"),
		instrument.WithDescription("number of inventory lookups"),
	)
	if err != nil {
		return fmt.Errorf("create inventory counter: %w", err)
	}

	inventoryLatencyHistogram, err = meter.Float64Histogram(
		"lab.inventory.latency",
		instrument.WithUnit("ms"),
		instrument.WithDescription("inventory lookup latency"),
	)
	if err != nil {
		return fmt.Errorf("create latency histogram: %w", err)
	}

	// Prime the instruments to surface export errors during bootstrap instead of at first request.
	inventoryLatencyHistogram.Record(ctx, 0)
	return nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("otca-lab/app")
	_, span := tracer.Start(ctx, "renderRoot")
	defer span.End()

	inventoryTotal, err := simulateInventory(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "failed to load inventory", http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"message":         "OTCA practice lab ready",
		"inventory_total": inventoryTotal,
		"timestamp":       time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "encode error", http.StatusInternalServerError)
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func inventoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	attrs := attribute.NewSet(
		attribute.String("http.method", r.Method),
		attribute.String("http.route", "/inventory"),
	)

	start := time.Now()
	total, err := simulateInventory(ctx)
	inventoryLatencyHistogram.Record(ctx, float64(time.Since(start).Milliseconds()), metric.WithAttributeSet(attrs))

	if err != nil {
		inventoryRequestCounter.Add(ctx, 1, metric.WithAttributeSet(attrs), metric.WithAttributes(attribute.String("status", "error")))
		http.Error(w, "inventory error", http.StatusInternalServerError)
		return
	}

	inventoryRequestCounter.Add(ctx, 1, metric.WithAttributeSet(attrs), metric.WithAttributes(attribute.String("status", "success")))

	response := map[string]any{
		"inventory_total": total,
		"items": []map[string]any{
			{"sku": "KUBE-C101", "qty": rand.Intn(25) + 1},
			{"sku": "OTEL-T403", "qty": rand.Intn(25) + 1},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "encode error", http.StatusInternalServerError)
	}
}

func simulateInventory(ctx context.Context) (int, error) {
	tracer := otel.Tracer("otca-lab/inventory")
	_, span := tracer.Start(ctx, "simulateInventory", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	// Add a small random delay to produce measurable spans.
	simulatedLatency := time.Duration(50+rand.Intn(120)) * time.Millisecond
	select {
	case <-time.After(simulatedLatency):
		span.SetAttributes(
			attribute.Int("inventory.latency_ms", int(simulatedLatency.Milliseconds())),
		)
	case <-ctx.Done():
		return 0, ctx.Err()
	}

	// Simulate occasional failure to mirror troubleshooting exercises.
	if rand.Intn(100) < 5 {
		err := fmt.Errorf("inventory backend timeout")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}

	total := rand.Intn(50) + 50
	span.SetAttributes(attribute.Int("inventory.total", total))
	return total, nil
}

func envOr(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
