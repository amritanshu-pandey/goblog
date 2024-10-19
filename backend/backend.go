package backend

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"amritanshu.in/goblog/md"
	"amritanshu.in/goblog/views"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func initMeter() (*sdkmetric.MeterProvider, error) {
	exp, err := prometheus.New()
	// exp, err := stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	// mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp, metric.WithInterval(time.Second * 5))))
	mp := sdkmetric.NewMeterProvider(metric.WithReader(exp))
	otel.SetMeterProvider(mp)
	return mp, nil
}

func RunServer(markdownPath string, assetsDir string, serverPort int, serverBindAddr string) error {
	mux := http.NewServeMux()
	markdownPosts, err := md.ActivePosts(markdownPath)
	if err != nil {
		return err
	}
	sortedTitles, err := md.SortedPostsByDate(markdownPath)
	if err != nil {
		return err
	}

	mp, err := initMeter()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := mp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down meter provider: %v", err)
		}
	}()

	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		handler := otelhttp.NewHandler(otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc)), "")
		fmt.Printf("Mux is handling connection for pattern: %s\n", pattern)
		mux.Handle(pattern, handler)
	}

	// Index Page
	handleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		views.Index(markdownPosts, sortedTitles).Render(r.Context(), w)
	})

	// Static assets
	staticFs := http.FileServer(http.Dir(assetsDir))
	handleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/assets", staticFs).ServeHTTP(w, r)
	})

	handleFunc("/article/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		views.Article(markdownPosts[slug]).Render(r.Context(), w)
	})

	handleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	fmt.Printf("Starting server on %s:%d\n", serverBindAddr, serverPort)
	http.ListenAndServe(fmt.Sprintf("%s:%d", serverBindAddr, serverPort), mux)
	return nil
}
