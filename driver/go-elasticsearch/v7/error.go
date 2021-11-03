package elasticsearch7

const codeUnknown errCode = 0

type errCode = int

type err struct {
	code    errCode
	message string
}

func (e err) Error() string {
	return e.message
}

func newError(code errCode, message string) *err {
	return &err{
		code:    code,
		message: message,
	}
}
