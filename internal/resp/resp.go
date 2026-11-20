package resp

import (
	"bufio"
	"errors"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/aurum77/anahtar/internal/handler"
	"github.com/aurum77/anahtar/internal/resp/responseTypes"
)

type ParseResult struct {
	Commands []string
	Err      error
}

func Parse(reader *bufio.Reader) ParseResult {
	var cmdArray []string

	commandCountByte, err := reader.ReadString('\n')

	if errors.Is(err, io.EOF) {
		if commandCountByte != "" {
			return ParseResult{Commands: nil, Err: ErrIncompleteRESPCommandLengthHeader}

		}
		return ParseResult{Commands: nil, Err: io.EOF}
	}

	if err != nil {
		return ParseResult{Commands: nil, Err: err}
	}

	if strings.HasPrefix(commandCountByte, "*") {
		commandCountByte = strings.TrimRight(commandCountByte, "\r\n")
		commandCountByte = strings.TrimPrefix(commandCountByte, "*")

		commandCount, err := strconv.Atoi(commandCountByte)

		if err != nil {
			return ParseResult{Commands: nil, Err: ErrIncompleteRESPCommandLengthHeader}
		}

		if commandCount <= 0 {
			return ParseResult{Commands: nil, Err: ErrCommandCountCannotBeZero}
		}

		if commandCount > 128 {
			return ParseResult{Commands: nil, Err: ErrCommandCountCannotExceedLimit}
		}

		for range commandCount {
			argumentLengthByte, err := reader.ReadString('\n')

			if err != nil {
				return ParseResult{Commands: nil, Err: ErrIncompleteRESPRequest}
			}

			if strings.HasPrefix(argumentLengthByte, "$") {
				argumentLengthByte = strings.TrimRight(argumentLengthByte, "\r\n")
				argumentLengthByte = strings.TrimPrefix(argumentLengthByte, "$")

				argumentLength, err := strconv.Atoi(argumentLengthByte)

				if err != nil {
					return ParseResult{Commands: nil, Err: ErrIncompleteRESPRequest}
				}

				if argumentLength < 0 {
					return ParseResult{Commands: nil, Err: ErrArgumentLengthCannotBeNegative}
				}

				if argumentLength > 1024*1024 {
					return ParseResult{Commands: nil, Err: ErrArgumentSizeCannotExceedSizeLimit}
				}

				argumentBytes := make([]byte, argumentLength)

				_, err = io.ReadFull(reader, argumentBytes)

				if err != nil {
					return ParseResult{Commands: nil, Err: ErrIncompleteRESPRequest}
				}

				crlf := make([]byte, 2)

				_, err = io.ReadFull(reader, crlf)

				if err != nil {
					return ParseResult{Commands: nil, Err: ErrIncompleteRESPRequest}
				}

				if string(crlf) != "\r\n" {
					return ParseResult{Commands: nil, Err: ErrInvalidRESPDivider}
				}

				argument := string(argumentBytes)

				cmdArray = append(cmdArray, argument)

				log.Printf("argument is %s", argument)
			} else {
				return ParseResult{Commands: nil, Err: ErrInvalidArgumentCountPrefix}

			}
		}
	} else {
		return ParseResult{Commands: nil, Err: ErrInvalidCommandCountPrefix}
	}

	log.Println(cmdArray)
	return ParseResult{Commands: cmdArray, Err: nil}
}

func ConvertResponse(response handler.Response) string {
	var responseBuilder strings.Builder

	switch response.Type {
	case responseTypes.SimpleString:
		responseBuilder.WriteString("+")
		responseBuilder.WriteString(response.Message)
		responseBuilder.WriteString("\r\n")

	case responseTypes.BulkString:
		responseBuilder.WriteString("$")
		responseBuilder.WriteString(strconv.Itoa(len(response.Message)))
		responseBuilder.WriteString("\r\n")
		responseBuilder.WriteString(response.Message)
		responseBuilder.WriteString("\r\n")

	case responseTypes.Error:
		responseBuilder.WriteString("-")
		responseBuilder.WriteString(response.Message)
		responseBuilder.WriteString("\r\n")

	case responseTypes.Null:
		responseBuilder.WriteString("$-1")
		responseBuilder.WriteString("\r\n")
	}

	return responseBuilder.String()
}
