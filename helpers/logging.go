package helpers

import (
	"fmt"
	"os"
)

var LogFile *os.File

func Println(args ...any) {
	args = append(args, "\n")
	fmt.Fprint(LogFile, args...)
}
