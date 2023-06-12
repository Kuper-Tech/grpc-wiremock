package schematomock

import (
	"fmt"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
)

type statusLevel string

const (
	StatusCodeLevel1   statusLevel = "1xx"
	StatusCodeLevel2   statusLevel = "2xx"
	StatusCodeLevel3   statusLevel = "3xx"
	StatusCodeLevel4   statusLevel = "4xx"
	StatusCodeLevel5   statusLevel = "5xx"
	StatusCodeLevelAll statusLevel = "default"
)

func toStatusLevel(value string) statusLevel {
	return statusLevel(value)
}

var statusCodesByLevel = map[statusLevel][]int{
	StatusCodeLevel1: {100, 101, 102, 103},
	StatusCodeLevel2: {200, 201, 202, 203, 204, 205, 206, 207, 208, 226},
	StatusCodeLevel3: {300, 301, 302, 303, 304, 305, 306, 307, 308},
	StatusCodeLevel4: {400, 401, 402, 403, 404, 405, 406, 407, 408, 409, 410, 411,
		412, 413, 414, 415, 416, 417, 418, 421, 422, 423, 424, 425, 426, 428, 429, 431, 451},
	StatusCodeLevel5: {500, 501, 502, 503, 504, 505, 506, 507, 508, 510, 511},
	StatusCodeLevelAll: {
		100, 101, 102, 103,
		200, 201, 202, 203, 204, 205, 206, 207, 208, 226,
		300, 301, 302, 303, 304, 305, 306, 307, 308, 400, 401, 402, 403, 404, 405, 406, 407, 408, 409, 410, 411,
		412, 413, 414, 415, 416, 417, 418, 421, 422, 423, 424, 425, 426, 428, 429, 431, 451,
		500, 501, 502, 503, 504, 505, 506, 507, 508, 510, 511,
	},
}

func getDefinedStatusCodes(responses openapi3.Responses) map[int]struct{} {
	isIntF := func(value string) (int, bool) {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, false
		}
		return intValue, true
	}

	definedStatusCodes := make(map[int]struct{})

	for statusCodeRaw := range responses {
		if value, isInt := isIntF(statusCodeRaw); isInt {
			definedStatusCodes[value] = struct{}{}
		}
	}

	return definedStatusCodes
}

// getStatusCodes returns status codes will be saved as mock.
func getStatusCodes(definedStatuses map[int]struct{}, statusCodeRaw string) ([]int, error) {
	// Simple status code such as `200`, `400`, `500`, etc.
	statusCodeInt, err := strconv.Atoi(statusCodeRaw)
	if err == nil {
		return []int{statusCodeInt}, nil
	}

	// Skip `default` status code.
	if toStatusLevel(statusCodeRaw) == StatusCodeLevelAll {
		return nil, nil
	}

	// `1xx`, `2xx`, `3xx`, `4xx`, `5xx` status codes.
	// Means that all matched statuses have the same response.
	statusesByLevel, isExists := statusCodesByLevel[toStatusLevel(statusCodeRaw)]
	if !isExists {
		return nil, fmt.Errorf("provided status level doesn't exist: %s", statusCodeRaw)
	}

	var statuses []int

	for _, status := range statusesByLevel {
		if _, exists := definedStatuses[status]; !exists {
			statuses = append(statuses, status)
		}
	}

	return statuses, nil
}
