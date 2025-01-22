package guap

import (
	"reflect"
	"regexp"
)

type RouteManager struct {
	Routes []Route
}

func NewRouteManager() *RouteManager {
	return &RouteManager{
		Routes: make([]Route, 0),
	}
}

func (rm *RouteManager) RegisterRoute(method Method, path string, handler interface{}) {

	pattern := createRoutePattern(path)

	handlerValue := reflect.ValueOf(handler)
	if handlerValue.Kind() != reflect.Func {
		panic("handler must be a function")
	}

	rm.Routes = append(rm.Routes, Route{
		path:       path,
		pattern:    pattern,
		methodType: method,
		handler:    handlerValue,
	})

}

func createRoutePattern(route string) *regexp.Regexp {
	re := regexp.MustCompile(`\{(\w+)}`)
	regexPattern := "^" + re.ReplaceAllString(route, "(?P<$1>[^/]+)") + "$"
	return regexp.MustCompile(regexPattern)
}
