package main

import (
	"Guap/guap"
	"fmt"
)

func main() {

	server := guap.NewAPIServer(":8080")

	server.RouteManager.RegisterRoute(guap.Get, "/test", testCall)

	server.RouteManager.RegisterRoute(guap.Post, "/test", testCall)

	server.RouteManager.RegisterRoute(guap.Get, "/test/{id}", testCallWithId)

	server.RouteManager.RegisterRoute(guap.Get, "/test/{id}/{text}", testCallWithIdAndText)

	err := server.Start(nil)
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
