package guap

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
)

type Middleware func(http.Handler) http.HandlerFunc

type APIServer struct {
	Addr         string
	RouteManager *RouteManager
}

type APIOptions struct {
	Middleware Middleware
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		Addr:         addr,
		RouteManager: NewRouteManager(),
	}
}

func (s *APIServer) Start(options *APIOptions) error {

	http.HandleFunc("/", s.routeHandler)

	log.Printf("Starting API server at %s", s.Addr)

	return http.ListenAndServe(s.Addr, nil)
}

func (s *APIServer) routeHandler(w http.ResponseWriter, r *http.Request) {
	methodAllowed := true
	for _, route := range s.RouteManager.Routes {
		if r.Method == route.methodType.String() && route.pattern.MatchString(r.URL.Path) {

			params := extractParameters(route.pattern, r.URL.Path)

			invokeHandler(route.handler, params, w)
			return
		} else if r.Method != route.methodType.String() && route.pattern.MatchString(r.URL.Path) {
			methodAllowed = false
		}
	}

	if !methodAllowed {
		http.Error(w, "405 Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.NotFound(w, r)
}

func extractParameters(pattern *regexp.Regexp, path string) map[string]string {
	match := pattern.FindStringSubmatch(path)
	params := map[string]string{}
	if match == nil {
		return params
	}

	for i, name := range pattern.SubexpNames() {
		if i > 0 && name != "" {
			params[name] = match[i]
		}
	}
	return params
}

func invokeHandler(handler reflect.Value, params map[string]string, w http.ResponseWriter) {
	handlerType := handler.Type()
	var args []reflect.Value

	i := 0
	for _, value := range params {
		argType := handlerType.In(i)

		if argType.Kind() == reflect.String {
			args = append(args, reflect.ValueOf(value))
		} else if argType.Kind() == reflect.Int {
			val, _ := strconv.Atoi(value)
			args = append(args, reflect.ValueOf(val))
		} else {
			log.Fatalf("Invalid argument type: %s", argType.String())
		}
		i++
	}

	results := handler.Call(args)
	if len(results) > 0 {
		_, err := fmt.Fprint(w, results[0].Interface())
		if err != nil {
			return
		}
	}
}
