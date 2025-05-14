package utils

import (
    "fmt"
    "log"
    "time"
    "regexp"
	"strconv"
)

type MongoQueryValues struct {
    FromDateTime time.Time
    ToDateTime   time.Time
    SortOrder    int
}

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

func StringToUint(strVal string) (uint, error) {
    unitVal, err := strconv.ParseUint(strVal, 10, 64)
    return uint(unitVal), err
}

func ExtractBinSizeAndUnit(s string) (int, string, error) {
	re := regexp.MustCompile(`^(\d+)([a-zA-Z]+)$`)

	matches := re.FindStringSubmatch(s)
	if len(matches) != 3 {
		return 0, "", fmt.Errorf("string does not match expected pattern")
	}

	num, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, "", err
	}

	return num, matches[2], nil
}

func ConvertFromTimeAndToTimeAndSortToMongoQueryValues(from, to, sort string) (MongoQueryValues, error) {
    mongoQueryValues := MongoQueryValues{}
    var err error

    mongoQueryValues.FromDateTime, err = time.Parse(time.RFC3339, from)
    if err != nil {
        return MongoQueryValues{}, err
    }
    mongoQueryValues.ToDateTime, err = time.Parse(time.RFC3339, to)
    if err != nil {
        return MongoQueryValues{}, err
    }

    if sort == "asc" {
        mongoQueryValues.SortOrder = 1
    } else if sort == "desc" {
        mongoQueryValues.SortOrder = -1
    } else {
        return MongoQueryValues{}, err
    }

    return mongoQueryValues, nil
}
