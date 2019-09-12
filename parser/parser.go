package parser

import (
	"errors"
	"fmt"
	"strconv"
)

// Parser parse value to correct type
type Parser interface {
	Parse(value string) (interface{}, error)
	IsIndexed() bool
}

type parserWrapper struct {
	f       ParseFunc
	indexed bool
}

func (p *parserWrapper) Parse(value string) (interface{}, error) {
	return p.f(value)
}

func (p *parserWrapper) IsIndexed() bool {
	return p.indexed
}

// ParseFunc is wrapper function for Parser
type ParseFunc func(value string) (interface{}, error)

var typesConverter = make(map[string]Parser)

// Register new parser
func Register(_type string, p Parser) {
	typesConverter[_type] = p
}

// RegisterWithFunc register new parser
func RegisterWithFunc(_type string, f ParseFunc) {
	typesConverter[_type] = &parserWrapper{f: f, indexed: false}
}

// RegisterIndexedFunc register new indexed parser
func RegisterIndexedFunc(_type string, f ParseFunc) {
	typesConverter[_type] = &parserWrapper{f: f, indexed: true}
}

// Unregister a parser
func Unregister(_type string) {
	delete(typesConverter, _type)
}

// FindParser find parser for particular type
func FindParser(_type string) (Parser, error) {
	if p, ok := typesConverter[_type]; ok {
		return p, nil
	}
	return nil, createError(_type)
}

func createError(_type string) error {
	return fmt.Errorf("Cannot find parser for type: %v", _type)
}

func init() {
	// Simple
	RegisterWithFunc("Int", ParseFunc(func(value string) (interface{}, error) {
		return strconv.ParseInt(value, 10, 64)
	}))
	RegisterWithFunc("String", ParseFunc(func(value string) (interface{}, error) {
		return value, nil
	}))
	RegisterWithFunc("Float", ParseFunc(func(value string) (interface{}, error) {
		return strconv.ParseFloat(value, 64)
	}))
	RegisterWithFunc("Boolean", ParseFunc(func(value string) (interface{}, error) {
		return strconv.ParseBool(value)
	}))
	// Indexed
	RegisterWithFunc("Indexed", ParseFunc(func(value string) (interface{}, error) {
		return nil, errors.New("Unsupported operation")
	}))
}
