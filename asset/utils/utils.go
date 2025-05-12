package utils

import (
    "log"
    "strconv"
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

func StringToUint(strVal string) uint {
    unitVal, err := strconv.ParseUint(strVal, 10, 64)
    FailOnError(err, "Failed to convert string to uint.")

    return uint(unitVal)
}
