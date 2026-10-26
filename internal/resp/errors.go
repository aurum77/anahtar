package resp

import "errors"

var (
	ErrIncompleteRESPCommandLengthHeader = errors.New("Incomplete RESP command length header")
	ErrIncompleteRESPRequest             = errors.New("Incomplete RESP request")

	ErrInvalidRESPDivider         = errors.New("Invalid RESP divider")
	ErrInvalidArgumentCountPrefix = errors.New("Invalid argument count byte prefix")
	ErrInvalidCommandCountPrefix  = errors.New("Invalid command count byte prefix")

	ErrCommandCountCannotBeZero      = errors.New("Command count can not be negative or zero")
	ErrCommandCountCannotExceedLimit = errors.New("Command count can not exceed 128 commands")

	ErrArgumentLengthCannotBeNegative    = errors.New("Argument length can not be negative")
	ErrArgumentSizeCannotExceedSizeLimit = errors.New("Argument size can not exceed 1 MiB")
)
