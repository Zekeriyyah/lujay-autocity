package pkg

import (
	"log"
	"os"
)

var (
	ErrorLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLog  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func Info(msg string) {
	InfoLog.Println(msg)
}

func Error(err error, msg string) {
	ErrorLog.Printf("%s: %v\n", msg, err)
}