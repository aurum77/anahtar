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

	if err == io.EOF {
		return ParseResult{Commands: nil, Err: err}
	}

	if err != nil {
		return ParseResult{Commands: nil, Err: err}
	}

	if strings.HasPrefix(commandCountByte, "*") {
		commandCountByte = strings.TrimRight(commandCountByte, "\r\n")
		commandCountByte = strings.TrimPrefix(commandCountByte, "*")

		commandCount, err := strconv.Atoi(commandCountByte)

		if err != nil {
			return ParseResult{Commands: nil, Err: errors.New("Error: Parser read failure")}
		}

		for range commandCount {
			argumentLengthByte, err := reader.ReadString('\n')

			if err != nil {
				return ParseResult{Commands: nil, Err: errors.New("Error: Parser read failure")}
			}

			if strings.HasPrefix(argumentLengthByte, "$") {
				argumentLengthByte = strings.TrimRight(argumentLengthByte, "\r\n")
				argumentLengthByte = strings.TrimPrefix(argumentLengthByte, "$")

				argumentLength, err := strconv.Atoi(argumentLengthByte)

				if err != nil {
					return ParseResult{Commands: nil, Err: errors.New("Error: Parser conversion failure")}
				}

				argumentBytes := make([]byte, argumentLength)

				_, err = io.ReadFull(reader, argumentBytes)

				if err != nil {
					return ParseResult{Commands: nil, Err: errors.New("Error: Parser read failure")}
				}

				crlf := make([]byte, 2)

				_, err = io.ReadFull(reader, crlf)

				if err != nil {
					return ParseResult{Commands: nil, Err: errors.New("Error: Parser read failure")}
				}

				argument := string(argumentBytes)

				cmdArray = append(cmdArray, argument)

				log.Printf("argument is %s", argument)
			}
		}
	}

	log.Println(cmdArray)
	return ParseResult{Commands: cmdArray, Err: nil}
}

func Response(response handler.Response) string {
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
