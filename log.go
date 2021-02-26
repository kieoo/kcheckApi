package main

import (
	"fmt"
	"log"
	"os"
)

func init() {

	logFile, err := os.OpenFile("./check.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0744)

	if err != nil {
		fmt.Println("open log file failed ", err)
		panic(1)
	}

	log.SetOutput(logFile)

	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
}
