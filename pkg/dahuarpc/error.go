package dahuarpc

import "encoding/json"

type ErrorType string

var (
	ErrorTypeInvalidLogin      ErrorType = "InvalidLogin"
	ErrorTypeInvalidSession    ErrorType = "InvalidSession"
	ErrorTypeInvalidRequest    ErrorType = "InvalidRequest"
	ErrorTypeMethodNotFound    ErrorType = "MethodNotFound"
	ErrorTypeInterfaceNotFound ErrorType = "InterfaceNotFound"
	ErrorTypeNoData            ErrorType = "NoData"
	ErrorTypeUnknown           ErrorType = "Unknown"
)

const (
	LoginMessageUserOrPasswordNotValid = "User or password not valid"
	LoginMessageUserNotValid           = "User not valid"
	LoginMessagePasswordNotValid       = "Password not valid"
	LoginMessageInBlackList            = "User in blackList"
	LoginMessageHasBeedUsed            = "User has be used"
	LoginMessageHasBeenLocked          = "User locked"
)

func newError(code int, message string) (ErrorType, string) {
	// Login
	switch code {
	case 268632085:
		return ErrorTypeInvalidLogin, LoginMessageUserOrPasswordNotValid
	case 268632081:
		return ErrorTypeInvalidLogin, LoginMessageHasBeenLocked
	}
	switch message {
	case "UserNotValidt":
		return ErrorTypeInvalidLogin, LoginMessageUserNotValid
	case "PasswordNotValid":
		return ErrorTypeInvalidLogin, LoginMessagePasswordNotValid
	case "InBlackList":
		return ErrorTypeInvalidLogin, LoginMessageInBlackList
	case "HasBeedUsed":
		return ErrorTypeInvalidLogin, LoginMessageHasBeedUsed
	case "HasBeenLocked":
		return ErrorTypeInvalidLogin, LoginMessageHasBeenLocked
	}

	// Default
	switch code {
	case 268894209:
		return ErrorTypeInvalidRequest, message
	case 268894210:
		return ErrorTypeMethodNotFound, message
	case 268632064:
		return ErrorTypeInterfaceNotFound, message
	case 285409284:
		return ErrorTypeNoData, message
	case 287637505, 287637504:
		return ErrorTypeInvalidSession, message
	default:
		return ErrorTypeUnknown, message
	}
}

type Error struct {
	Method  string
	Code    int
	Message string
	Type    ErrorType
}

func (r *Error) Error() string {
	return r.Message
}

func (r *Error) UnmarshalJSON(data []byte) error {
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	r.Code = res.Code
	r.Type, r.Message = newError(res.Code, res.Message)

	return nil
}
