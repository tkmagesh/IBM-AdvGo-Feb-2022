package main

import (
	"context"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"

	"distributed-tracing/lib/tracing"
)

func main() {
	tracing.Init("trace-demo")

	http.HandleFunc("/formatGreeting", handleFormatGreeting)

	log.Print("Listening on http://localhost:8082/")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func handleFormatGreeting(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("trace-demo")
	ctx, span := tracer.Start(r.Context(), "/formatGreeting")
	defer span.End()

	name := r.FormValue("name")
	title := r.FormValue("title")
	descr := r.FormValue("description")

	greeting := FormatGreeting(ctx, name, title, descr)
	w.Write([]byte(greeting))
}

// FormatGreeting combines information about a person into a greeting string.
func FormatGreeting(
	ctx context.Context,
	name, title, description string,
) string {
	tracer := otel.Tracer("trace-demo")
	_, span := tracer.Start(ctx, "format-greeting")
	defer span.End()

	response := "Hello, "
	if title != "" {
		response += title + " "
	}
	response += name + "!"
	if description != "" {
		response += " " + description
	}
	return response
}
