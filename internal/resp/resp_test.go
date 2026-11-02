package resp_test

import (
	"bufio"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/aurum77/anahtar/internal/resp"
)

func TestValidInput(t *testing.T) {
	expected := []string{"SET", "foo", "bar"}

	input := strings.NewReader("*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n")
	reader := bufio.NewReader(input)

	got := resp.Parse(reader)

	if got.Err != nil {
		t.Fatalf("Got error: %s", got.Err)
	}

	if !reflect.DeepEqual(got.Commands, expected) {
		t.Fatalf("Parse = %v; want [SET foo bar]", got.Commands)
	}
}

func TestInvalidInput(t *testing.T) {
	input := strings.NewReader("?3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n")
	reader := bufio.NewReader(input)

	got := resp.Parse(reader)

	if got.Err == nil {
		t.Fatalf("Got no error on invalid input")
	}

	if !errors.Is(got.Err, resp.ErrInvalidCommandCountPrefix) {
		t.Fatalf("Got error: %s; want ErrInvalidCommandCountPrefix", got.Err)
	}
}

func TestTruncatedInput(t *testing.T) {
	input := strings.NewReader("*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$5\r\nba")
	reader := bufio.NewReader(input)

	got := resp.Parse(reader)

	if got.Err == nil {
		t.Fatalf("Got no error on truncated input")
	}

	if !errors.Is(got.Err, resp.ErrIncompleteRESPRequest) {
		t.Fatalf("Got error: %s; want ErrIncompleteRESPRequest", got.Err)
	}
}

func TestInvalidDivider(t *testing.T) {
	input := strings.NewReader("*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\n\r")
	reader := bufio.NewReader(input)

	got := resp.Parse(reader)

	if got.Err == nil {
		t.Fatalf("Got no error on invalid divider")
	}

	if !errors.Is(got.Err, resp.ErrInvalidRESPDivider) {
		t.Fatalf("Got error: %s; want ErrInvalidRESPDivider", got.Err)
	}
}

func TestCommandCountLimit(t *testing.T) {
	input := strings.NewReader("*129\r\n")
	reader := bufio.NewReader(input)

	got := resp.Parse(reader)

	if got.Err == nil {
		t.Fatalf("Got no error on command count")
	}

	if !errors.Is(got.Err, resp.ErrCommandCountCannotExceedLimit) {
		t.Fatalf("Got error: %s; want ErrCommandCountCannotExceedLimit", got.Err)
	}
}

func TestArgumentSizeLimit(t *testing.T) {
	input := strings.NewReader("*1\r\n$1048577\r\n")
	reader := bufio.NewReader(input)

	got := resp.Parse(reader)

	if got.Err == nil {
		t.Fatalf("Got no error on argument size")
	}

	if !errors.Is(got.Err, resp.ErrArgumentSizeCannotExceedSizeLimit) {
		t.Fatalf("Got error: %s; want ErrArgumentSizeCannotExceedSizeLimit", got.Err)
	}
}

func TestZeroCommandInput(t *testing.T) {
	input := strings.NewReader("*0\r\n")
	reader := bufio.NewReader(input)

	got := resp.Parse(reader)

	if got.Err == nil {
		t.Fatalf("Got no error on empty input")
	}

	if !errors.Is(got.Err, resp.ErrCommandCountCannotBeZero) {
		t.Fatalf("Got error: %s; want ErrCommandCountCannotBeZero", got.Err)
	}
}

func TestNegativeArgumentLength(t *testing.T) {
	input := strings.NewReader("*1\r\n$-1\r\n")
	reader := bufio.NewReader(input)

	got := resp.Parse(reader)

	if got.Err == nil {
		t.Fatalf("Got no error on negative argument length")
	}

	if !errors.Is(got.Err, resp.ErrArgumentLengthCannotBeNegative) {
		t.Fatalf("Got error: %s; want ErrArgumentLengthCannotBeNegative", got.Err)
	}
}

func TestZeroLengthArgument(t *testing.T) {
	expected := []string{""}

	input := strings.NewReader("*1\r\n$0\r\n\r\n")
	reader := bufio.NewReader(input)

	got := resp.Parse(reader)

	if got.Err != nil {
		t.Fatalf("Got error: %s", got.Err)
	}

	if !reflect.DeepEqual(got.Commands, expected) {
		t.Fatalf("Parse = %v; want [\"\"]", got.Commands)
	}
}
