package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"regexp"
	"strconv"
	"time"
)

const (
	DisablePaginationHeader = "X-Disable-Pagination"
	ErrorLogFormat          = "got err: %v, context: %s - %s"
	ShutDownEvent           = "ShutDownEvent"
)

func IsValidSanitizeSQL(queryParam string) bool {
	regexQueryParam := regexp.MustCompile(`^[\w ]+$`)
	return regexQueryParam.MatchString(queryParam)
}

func Contains(elems []string, elem string) bool {
	for _, e := range elems {
		if elem == e {
			return true
		}
	}
	return false
}

func ContainsError(err error, errorTargetSlice []error) bool {
	if len(errorTargetSlice) > 0 {
		for _, errSingle := range errorTargetSlice {
			if errors.Is(err, errSingle) {
				return true
			}
		}
	}
	return false
}

func StringUnitToDuration(input string) time.Duration {
	durationMapping := map[string]time.Duration{
		"second": time.Second,
		"minute": time.Minute,
		"hour":   time.Hour,
		"day":    24 * time.Hour,
		"week":   7 * 24 * time.Hour,
		"month":  30 * 24 * time.Hour,
		"year":   365 * 24 * time.Hour,
	}

	switch input {
	case "second":
		return durationMapping["second"]
	case "minute":
		return durationMapping["minute"]
	case "hour":
		return durationMapping["hour"]
	case "week":
		return durationMapping["week"]
	case "month":
		return durationMapping["month"]
	case "year":
		return durationMapping["year"]
	default:
		return time.Second
	}
}

func IsDisablePagination(ctx *gin.Context) bool {
	disablePaginationHeaderVal := ctx.Request.Header.Get(DisablePaginationHeader)
	if disablePaginationHeaderVal == "" {
		return false
	}

	boolValue, err := strconv.ParseBool(disablePaginationHeaderVal)
	if err != nil {
		fmt.Println("Error on IsDisablePagination, strconv.ParseBool:", err)
		return false
	}
	return boolValue
}
