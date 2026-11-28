package handler_test

import (
	"testing"

	"github.com/aurum77/anahtar/internal/handler"
	"github.com/aurum77/anahtar/internal/resp"
	"github.com/aurum77/anahtar/internal/resp/responseTypes"
)

func TestSimpleString(t *testing.T) {
	expected := "+OK\r\n"
	handlerResponse := handler.Response{Type: responseTypes.SimpleString, Message: "OK"}

	got := resp.ConvertResponse(handlerResponse)

	if got != expected {
		t.Fatalf("got %v; expected %s", got, expected)
	}
}

func TestBulkString(t *testing.T) {
	expected := "$3\r\nfoo\r\n"
	handlerResponse := handler.Response{Type: responseTypes.BulkString, Message: "foo"}

	got := resp.ConvertResponse(handlerResponse)

	if got != expected {
		t.Fatalf("got %v; expected %s", got, expected)
	}
}

func TestEmptyBulkString(t *testing.T) {
	expected := "$0\r\n\r\n"
	handlerResponse := handler.Response{Type: responseTypes.BulkString, Message: ""}

	got := resp.ConvertResponse(handlerResponse)

	if got != expected {
		t.Fatalf("got %v; expected %s", got, expected)
	}
}

func TestNull(t *testing.T) {
	expected := "$-1\r\n"
	handlerResponse := handler.Response{Type: responseTypes.Null, Message: ""}

	got := resp.ConvertResponse(handlerResponse)

	if got != expected {
		t.Fatalf("got %v; expected %s", got, expected)
	}
}

func TestUnicode(t *testing.T) {
	//
	expected := "$5\r\nTürk\r\n"
	handlerResponse := handler.Response{Type: responseTypes.BulkString, Message: "Türk"}

	got := resp.ConvertResponse(handlerResponse)

	if got != expected {
		t.Fatalf("got %v; expected %s", got, expected)
	}
}

func TestError(t *testing.T) {
	expected := "-ERR something went wrong\r\n"
	handlerResponse := handler.Response{Type: responseTypes.Error, Message: "ERR something went wrong"}

	got := resp.ConvertResponse(handlerResponse)

	if got != expected {
		t.Fatalf("got %v; expected %s", got, expected)
	}
}
