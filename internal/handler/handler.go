package handler

import (
	"log"
	"strings"

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

	command := stringArray[0]

	log.Printf("COMMAND: %s", command)

	switch command {
	case "PING":
		return Response{true, PONG}

	case "GET":
		if len(stringArray) != 2 {
			return Response{false, WRONG_ARGUMENT_COUNT}
		}

		key := stringArray[1]

		value, ok := h.store.Get(key)

		if !ok {
			return Response{false, KEY_NOT_FOUND}
		}

		return Response{true, value}

	case "SET":
		if len(stringArray) != 3 {
			return Response{false, WRONG_ARGUMENT_COUNT}
		}

		key := stringArray[1]

		value := stringArray[2]

		h.store.Set(key, value,)

		return Response{true, OK}

	case "DEL":
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
