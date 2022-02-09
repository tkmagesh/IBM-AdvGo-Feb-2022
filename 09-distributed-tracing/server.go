package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"distributed-tracing/lib/tracing"
	"distributed-tracing/people"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var repo *people.Repository
var tracer trace.Tracer

func main() {
	tracer = tracing.Init("hello-server")
	repo = people.NewRepository()
	defer repo.Close()

	http.HandleFunc("/sayHello/", handleSayHello)

	log.Print("Listening on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleSayHello(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "sayHello")
	defer span.End()

	name := strings.TrimPrefix(r.URL.Path, "/sayHello/")
	greeting, err := SayHello(name, ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	span.AddEvent("Sending Response :", trace.WithAttributes(attribute.String("greeting", greeting)))
	w.Write([]byte(greeting))
}

// SayHello creates a greeting for the named person.
func SayHello(name string, ctx context.Context) (string, error) {
	//accessing the span from the context
	span := trace.SpanFromContext(ctx)
	person, err := repo.GetPerson(name, ctx)
	if err != nil {
		return "", err
	}
	span.SetAttributes(
		attribute.String("name", person.Name),
		attribute.String("title", person.Title),
		attribute.String("desc", person.Description),
	)
	return FormatGreeting(
		person.Name,
		person.Title,
		person.Description,
		ctx,
	), nil
}

// FormatGreeting combines information about a person into a greeting string.
func FormatGreeting(name, title, description string, ctx context.Context) string {
	_, childSpan := tracer.Start(ctx, "format-greeting")
	defer childSpan.End()

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
