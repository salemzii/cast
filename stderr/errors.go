package stderr

import "log"

type StdError struct {
	Error   error
	Message string
}

func (err StdError) HandleError() {
	log.Printf("%s:%s", err.Message, err.Error)
}
