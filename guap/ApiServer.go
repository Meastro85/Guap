package guap

import (
	"encoding/json"
	"fmt"
	"io"
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

			invokeHandler(route, params, w, r)
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

func parseParam(argType reflect.Type, value string) reflect.Value {

	if argType.Kind() == reflect.String {
		return reflect.ValueOf(value)
	} else if argType.Kind() == reflect.Int {
		val, _ := strconv.Atoi(value)
		return reflect.ValueOf(val)
	} else {
		log.Fatalf("Invalid argument type: %s", argType.String())
		return reflect.Value{}
	}
}

func getParams(handler reflect.Value, params map[string]string, r *http.Request) []reflect.Value {
	handlerType := handler.Type()
	var args []reflect.Value
	i := 0
	paramCount := handlerType.NumIn()
	for _, value := range params {
		argType := handlerType.In(i)

		param := parseParam(argType, value)
		if param.IsValid() {
			args = append(args, param)
		}
		i++

	}

	if i == paramCount-1 {
		argType := handlerType.In(i)
		contentType := r.Header.Get("Content-Type")
		val, err := getBody(r, argType, contentType)
		if err != nil {
			log.Fatalf("Invalid argument type: %s", argType.String())
		}
		if val.IsValid() {
			args = append(args, val)
		}
	}

	return args
}

func getBody(r *http.Request, argType reflect.Type, contentType string) (reflect.Value, error) {

	if r.Body == nil {
		return reflect.Value{}, nil
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return reflect.Value{}, err
	}

	argValue := reflect.New(argType).Interface()

	if contentType == "application/json" {
		if err := json.Unmarshal(body, argValue); err != nil {
			return reflect.Value{}, err
		}
	}

	return reflect.ValueOf(argValue).Elem(), nil
}

func invokeHandler(route Route, params map[string]string, w http.ResponseWriter, r *http.Request) {
	args := getParams(route.handler, params, r)

	results := route.handler.Call(args)
	if len(results) > 0 {
		_, err := fmt.Fprint(w, results[0].Interface())
		if err != nil {
			return
		}
	}
}
