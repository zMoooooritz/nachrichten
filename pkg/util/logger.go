package util

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

var Logger *log.Logger

func init() {
	Logger = log.New(io.Discard, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func SetLogFile(filePath string) error {
	directory := filepath.Dir(filePath)
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		log.Fatalln("Error dir")
		return err
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Error file")
		return err
	}

	Logger.SetOutput(file)
	return nil
}
