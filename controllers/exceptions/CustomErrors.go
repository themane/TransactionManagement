package exceptions

type NoSuchCombinationError struct {
	Message string
}

func (m *NoSuchCombinationError) Error() string {
	if len(m.Message) > 0 {
		return m.Message
	}
	return "no such combination exists"
}

type AlreadyExistsError struct {
	Message string
}

func (m *AlreadyExistsError) Error() string {
	return m.Message
}
