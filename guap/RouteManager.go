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
	firstChar := path[0]
	if firstChar != '/' {
		panic("First character must be /")
	}

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

func (rm *RouteManager) RegisterRoutes(routes []BasicRoute) {
	for _, route := range routes {
		rm.RegisterRoute(route.method, route.path, route.handler)
	}
}

func createRoutePattern(route string) *regexp.Regexp {
	re := regexp.MustCompile(`\{(\w+)}`)
	regexPattern := "^" + re.ReplaceAllString(route, "(?P<$1>[^/]+)") + "$"
	return regexp.MustCompile(regexPattern)
}
