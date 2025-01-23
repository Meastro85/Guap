package guap

import (
	"reflect"
	"regexp"
)

type Method string

const (
	Get    Method = "GET"
	Post   Method = "POST"
	Delete Method = "DELETE"
	Put    Method = "PUT"
)

func (m Method) String() string {
	return string(m)
}

type BasicRoute struct {
	method  Method
	path    string
	handler reflect.Value
}

type Route struct {
	path       string
	pattern    *regexp.Regexp
	methodType Method
	handler    reflect.Value
}
