package paper

import (
	"log"
)

var (
	LogPrefix = "<===[rekaime]===>\n"
)

func Info(msg... any) {
	log.Printf("%s", LogPrefix)
	log.Print(msg...)
}

func Fatal(msg... any) {
	log.Printf("%s", LogPrefix)
	log.Fatal(msg...)
}

func Err(e error) bool {
	if e == nil {
		return false
	}
	Info(e)
	return true
}

func ErrFatal(e error) bool {
	if !Err(e) {
		return false
	}
	Fatal(e)
	return true
}