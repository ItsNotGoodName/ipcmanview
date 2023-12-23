package dahuarpc

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type ErrorType string

var (
	ErrorTypeInvalidSession    ErrorType = "invalid session"
	ErrorTypeInvalidRequest    ErrorType = "InvalidRequest"
	ErrorTypeMethodNotFound    ErrorType = "MethodNotFound"
	ErrorTypeInterfaceNotFound ErrorType = "InterfaceNotFound"
	ErrorTypeNoData            ErrorType = "NoData"
	ErrorTypeUnknown           ErrorType = "Unknown"
)

func errorTypeFromCode(code int) ErrorType {
	switch code {
	case 268894209:
		return ErrorTypeInvalidRequest
	case 268894210:
		return ErrorTypeMethodNotFound
	case 268632064:
		return ErrorTypeInterfaceNotFound
	case 285409284:
		return ErrorTypeNoData
	case 287637505, 287637504:
		return ErrorTypeInvalidSession
	default:
		return ErrorTypeUnknown
	}
}

type Response[T any] struct {
	ID      int             `json:"id"`
	Session ResponseSession `json:"session"`
	Error   *ResponseError  `json:"error"`
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

// ---------- ResponseError

type ResponseError struct {
	Method  string
	Code    int
	Message string
	Type    ErrorType
}

func (r *ResponseError) Error() string {
	return r.Message
}

func (r *ResponseError) UnmarshalJSON(data []byte) error {
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	r.Code = res.Code
	r.Message = res.Message
	r.Type = errorTypeFromCode(r.Code)

	return nil
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
