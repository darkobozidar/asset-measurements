package utils

import (
    "log"
)

func FailOnError(err error, msg string) {
    if err != nil {
        log.Panicf("%s: %s", msg, err)
    }
}

func LogOnError(err error, msg string) {
    if err != nil {
        log.Printf("%s: %s", msg, err)
    }
}
