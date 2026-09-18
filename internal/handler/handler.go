package handler

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aurum77/anahtar/internal/store"
)

type Handler struct {
	store *store.Store
}

type Response struct {
	Ok      bool
	Message string
}

func New(s *store.Store) *Handler {
	return &Handler{store: s}
}

func (h *Handler) HandleCommand(requestString string) Response {
	log.Println("Handling command")

	stringArray := strings.Fields(requestString)

	if len(stringArray) == 0 {
		return Response{false, INVALID_COMMAND}
	}

	command := strings.ToLower(stringArray[0])

	log.Printf("COMMAND: %s", command)

	switch command {
	case "ping":
		return Response{true, PONG}

	case "get":
		if len(stringArray) != 2 {
			return Response{false, WRONG_ARGUMENT_COUNT}
		}

		key := stringArray[1]

		value, ok := h.store.Get(key)

		if !ok {
			return Response{false, KEY_NOT_FOUND}
		}

		return Response{true, value}

	case "set":
		var expiry time.Time

		if len(stringArray) != 3 && len(stringArray) != 5 {
			return Response{false, WRONG_ARGUMENT_COUNT}
		}

		key := stringArray[1]

		value := stringArray[2]

		if len(stringArray) == 5 {
			var multiplier time.Duration

			// ex(seconds) or px(milliseconds)e
			multiplierType := strings.ToLower(stringArray[3])

			switch multiplierType {
			case "ex":
				multiplier = time.Second
			case "px":
				multiplier = time.Millisecond
			default:
				return Response{false, INVALID_COMMAND}
			}

			amount, err := strconv.Atoi(stringArray[4])

			if err != nil {
				log.Printf("ERR: %s", err)
				return Response{false, INVALID_COMMAND}
			}

			if amount <= 0 {
				return Response{false, INVALID_EXPIRE_TIME_IN_SET}
			}

			expiry = time.Now().Add(time.Duration(amount * int(multiplier)))
		}

		h.store.Set(key, value, expiry)

		return Response{true, OK}

	case "del":
		if len(stringArray) != 2 {
			return Response{false, WRONG_ARGUMENT_COUNT}
		}

		key := stringArray[1]

		ok := h.store.Del(key)

		if !ok {
			return Response{false, KEY_NOT_FOUND}
		}

		return Response{true, OK}

	default:
		return Response{false, INVALID_COMMAND}
	}
}
