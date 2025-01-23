package main

import (
	"Guap/guap"
	"fmt"
	"log"
	"net/http"
)

func main() {

	server := guap.NewAPIServer(":8080")

	server.RouteManager.RegisterRoute(guap.Get, "/test", testCall)

	server.RouteManager.RegisterRoute(guap.Post, "/test", testPostCall)

	server.RouteManager.RegisterRoute(guap.Get, "/test/{id}", testCallWithId)

	server.RouteManager.RegisterRoute(guap.Post, "/test/{id}", testPostCallWithId)

	server.RouteManager.RegisterRoute(guap.Get, "/test/{id}/{text}", testCallWithIdAndText)

	options := guap.APIOptions{
		Middleware: LoggingMiddleware,
	}

	err := server.Start(&options)
	if err != nil {
		panic(err)
	}
}

func testCall() string {
	return fmt.Sprintf("Test call")
}

func testCallWithId(id int) string {
	return fmt.Sprintf("Test call with id: %d", id)
}

func testCallWithIdAndText(id int, text string) string {
	return fmt.Sprintf("Test call with id: %d, text: %s", id, text)
}

type postTest struct {
	Text   string
	Number int
}

func testPostCall(body postTest) string {
	return fmt.Sprintf("Test post call with text: %s, number: %d", body.Text, body.Number)
}

func testPostCallWithId(id int, body postTest) string {
	return fmt.Sprintf("Test post call with text: %s, number: %d, id: %d", body.Text, body.Number, id)
}

func LoggingMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Response completed for %s %s", r.Method, r.URL.Path)
	}
}
