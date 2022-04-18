package errors

import (
	"fmt"
	"log"
	"runtime"
)

func HandleError(err error) error {
	_, filename, line, _ := runtime.Caller(1)
	err = fmt.Errorf("\n[error] %s:%d %w\n", filename, line, err)

	log.Println(err)

	return err
}
