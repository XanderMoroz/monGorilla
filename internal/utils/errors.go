package utils

type Error string

const (
	ERR0303 Error = "ERR0303"
	ERR0304 Error = "ERR0304"
	ERR0401 Error = "ERR0401"
	ERR0402 Error = "ERR0402"
	ERR0403 Error = "ERR0403"
	ERR0404 Error = "ERR0404"
	ERR0405 Error = "ERR0405"
	ERR0406 Error = "ERR0406"
	ERR0407 Error = "ERR0407"
)

func (er Error) ToDescription() string {
	switch er {
	case ERR0303:
		return "Credentials are not provided!"
	case ERR0304:
		return "Invalid credentials provided!"
	case ERR0401:
		return "Authentication credentials are not provided!"
	case ERR0402:
		return "User not found!"
	case ERR0403:
		return "Incorrect password!"
	case ERR0404:
		return "Registered user found!"
	case ERR0405:
		return "Error occurred while creating Authentication Credential!"
	case ERR0406:
		return "Password's are not same!"
	case ERR0407:
		return "Error occurred while creating account!"
	default:
		return ""
	}
}

func ExceptionToString(err error) string {
	return err.Error()
}
