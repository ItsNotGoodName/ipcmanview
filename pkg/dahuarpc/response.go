package dahuarpc

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Response[T any] struct {
	ID      int             `json:"id"`
	Session ResponseSession `json:"session"`
	Error   *Error          `json:"error"`
	Params  T               `json:"params"`
	Result  ResponseResult  `json:"result"`
}

// ---------- ResponseSession

type ResponseSession string

func (s *ResponseSession) UnmarshalJSON(data []byte) error {
	// string -> string
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*s = ResponseSession(str)
		return nil
	}

	// int64 -> string
	var num int64
	if err := json.Unmarshal(data, &num); err == nil {
		*s = ResponseSession(strconv.FormatInt(num, 10))
		return nil
	}

	return fmt.Errorf("session is not a string or number")
}

func (s ResponseSession) String() string {
	return string(s)
}

// ---------- ResponseResult

type ResponseResult int64

func (s *ResponseResult) UnmarshalJSON(data []byte) error {
	// int64 -> int64
	var num int64
	if err := json.Unmarshal(data, &num); err == nil {
		*s = ResponseResult(num)
		return nil
	}

	// bool -> int64
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		if b {
			*s = 1
		}
		return nil
	}

	return fmt.Errorf("result is not a number or boolean")
}

func (s ResponseResult) Integer() int64 {
	return int64(s)
}

func (s ResponseResult) Bool() bool {
	return s == 1
}
