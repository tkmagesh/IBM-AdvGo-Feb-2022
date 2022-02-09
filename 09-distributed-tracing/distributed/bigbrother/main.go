package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"distributed-tracing/lib/tracing"
	"distributed-tracing/people"
)

var repo *people.Repository

func main() {
	tracing.Init("trace-demo")

	repo = people.NewRepository()
	defer repo.Close()

	http.HandleFunc("/getPerson/", handleGetPerson)

	log.Print("Listening on http://localhost:8081/")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func handleGetPerson(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("go-4-bigbrother")
	ctx, span := tracer.Start(r.Context(), "/getPerson")
	defer span.End()

	name := strings.TrimPrefix(r.URL.Path, "/getPerson/")
	person, err := repo.GetPerson(name, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	span.SetAttributes(
		attribute.String("name", person.Name),
		attribute.String("title", person.Title),
		attribute.String("description", person.Description),
	)

	bytes, _ := json.Marshal(person)
	w.Write(bytes)
}
