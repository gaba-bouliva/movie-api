package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrInvalidRuntimeFormat = errors.New("invalid runtime provided")
)

type Runtime struct{
	Value 		int32
}

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonStr := fmt.Sprintf("%d mins", r.Value)
	quotedJsonStr := strconv.Quote(jsonStr)
	return []byte(quotedJsonStr), nil
}

func (r *Runtime) UnmarshalJSON(data []byte) error {
	unquotedJsonStr, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	parts := strings.Split(unquotedJsonStr, " ") 

	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	parsedValue, err :=  strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	r.Value = int32(parsedValue)

	return nil
}