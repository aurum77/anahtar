package resp

import (
	"bufio"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/aurum77/anahtar/internal/handler"
)

func Parse(reader *bufio.Reader) []string {
	var cmdArray []string

	commandCountByte, err := reader.ReadString('\n')

	if err == io.EOF {
		return nil
	}

	if err != nil {
		log.Println("Read failure")
		return nil
	}

	if strings.HasPrefix(commandCountByte, "*") {
		commandCountByte = strings.TrimRight(commandCountByte, "\r\n")
		commandCountByte = strings.TrimPrefix(commandCountByte, "*")

		commandCount, err := strconv.Atoi(commandCountByte)

		if err != nil {
			log.Println("Read failure")
			return nil
		}

		for range commandCount {
			argumentLengthByte, err := reader.ReadString('\n')

			if err != nil {
				log.Println("Read failure")
				return nil
			}

			if strings.HasPrefix(argumentLengthByte, "$") {
				argumentLengthByte = strings.TrimRight(argumentLengthByte, "\r\n")
				argumentLengthByte = strings.TrimPrefix(argumentLengthByte, "$")

				argumentLength, err := strconv.Atoi(argumentLengthByte)

				if err != nil {
					log.Println("Conversion failure")
					return nil
				}

				argumentBytes := make([]byte, argumentLength)

				_, err = io.ReadFull(reader, argumentBytes)

				if err != nil {
					log.Println("Read failure")
					return nil
				}

				crlf := make([]byte, 2)

				_, err = io.ReadFull(reader, crlf)

				if err != nil {
					log.Println("Read failure")
					return nil
				}

				argument := string(argumentBytes)

				cmdArray = append(cmdArray, argument)

				log.Printf("argument is %s", argument)
			}
		}
	}

	log.Println(cmdArray)
	return cmdArray
}

func Response(response handler.Response) string {
	var responseBuilder strings.Builder

	if response.Ok {
		responseBuilder.WriteString("+")
	} else {
		responseBuilder.WriteString("-")
	}

	responseBuilder.WriteString(response.Message)
	responseBuilder.WriteString("\r\n")

	return responseBuilder.String()
}
