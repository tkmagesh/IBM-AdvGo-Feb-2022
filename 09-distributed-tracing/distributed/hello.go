package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	/* opentracing "github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
	*/
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	xhttp "distributed-tracing/lib/http"
	"distributed-tracing/lib/model"
	"distributed-tracing/lib/tracing"
)

func main() {
	tracing.Init("trace-demo")

	http.HandleFunc("/sayHello/", handleSayHello)

	log.Print("Listening on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleSayHello(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("trace-demo")
	ctx, span := tracer.Start(r.Context(), "say-hello")
	defer span.End()

	name := strings.TrimPrefix(r.URL.Path, "/sayHello/")
	greeting, err := SayHello(ctx, name)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	span.SetAttributes(attribute.String("response", greeting))
	w.Write([]byte(greeting))
}

// SayHello creates a greeting for the named person.
func SayHello(ctx context.Context, name string) (string, error) {
	person, err := getPerson(ctx, name)
	if err != nil {
		return "", err
	}

	return formatGreeting(ctx, person)
}

func getPerson(ctx context.Context, name string) (*model.Person, error) {
	res, err := xhttp.Get("http://localhost:8081/getPerson/" + name)
	if err != nil {
		return nil, err
	}
	var person model.Person
	if err := json.Unmarshal(res, &person); err != nil {
		return nil, err
	}
	return &person, nil
}

func formatGreeting(ctx context.Context, person *model.Person) (string, error) {
	v := url.Values{}
	v.Set("name", person.Name)
	v.Set("title", person.Title)
	v.Set("description", person.Description)
	url := "http://localhost:8082/formatGreeting?" + v.Encode()
	res, err := xhttp.Get(url)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
