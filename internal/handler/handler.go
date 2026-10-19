package handler

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aurum77/anahtar/internal/resp/responseTypes"
	"github.com/aurum77/anahtar/internal/store"
)

type Handler struct {
	store *store.Store
}

type Response struct {
	Type    int
	Message string
}

func New(s *store.Store) *Handler {
	return &Handler{store: s}
}

func (h *Handler) HandleCommand(commandArray []string) Response {
	log.Println("Handling command")

	command := strings.ToLower(commandArray[0])

	log.Printf("COMMAND: %s", command)

	switch command {
	case "ping":
		return Response{responseTypes.SimpleString, PONG}

	case "get":
		if len(commandArray) != 2 {
			return Response{responseTypes.Error, WRONG_ARGUMENT_COUNT}
		}

		key := commandArray[1]

		value, ok := h.store.Get(key)

		if !ok {
			return Response{responseTypes.Null, ""}
		}

		return Response{responseTypes.BulkString, value}

	case "set":
		var expiry time.Time

		if len(commandArray) != 3 && len(commandArray) != 5 {
			return Response{responseTypes.Error, WRONG_ARGUMENT_COUNT}
		}

		key := commandArray[1]

		value := commandArray[2]

		if len(commandArray) == 5 {
			var multiplier time.Duration

			// ex(seconds) or px(milliseconds)e
			multiplierType := strings.ToLower(commandArray[3])

			switch multiplierType {
			case "ex":
				multiplier = time.Second
			case "px":
				multiplier = time.Millisecond
			default:
				return Response{responseTypes.Error, INVALID_COMMAND}
			}

			amount, err := strconv.Atoi(commandArray[4])

			if err != nil {
				log.Printf("ERR: %s", err)
				return Response{responseTypes.Error, INVALID_COMMAND}
			}

			if amount <= 0 {
				return Response{responseTypes.Error, INVALID_EXPIRE_TIME_IN_SET}
			}

			expiry = time.Now().Add(time.Duration(amount * int(multiplier)))
		}

		h.store.Set(key, value, expiry)

		return Response{responseTypes.SimpleString, OK}

	case "del":
		if len(commandArray) != 2 {
			return Response{responseTypes.Error, WRONG_ARGUMENT_COUNT}
		}

		key := commandArray[1]

		ok := h.store.Del(key)

		if !ok {
			return Response{responseTypes.Error, KEY_NOT_FOUND}
		}

		return Response{responseTypes.SimpleString, OK}

	default:
		return Response{responseTypes.Error, INVALID_COMMAND}
	}
}
