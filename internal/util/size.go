package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var sizeRegex = regexp.MustCompile(`^(\d+)\s*([a-zA-Z]*)$`)

func ParseHumanSize(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty size")
	}

	matches := sizeRegex.FindStringSubmatch(s)
	if matches == nil {
		return 0, fmt.Errorf("invalid size format: %s", s)
	}

	value, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", matches[1])
	}

	unit := strings.ToUpper(matches[2])
	switch unit {
	case "", "B":
		return value, nil
	case "K", "KB":
		return value * 1024, nil
	case "M", "MB":
		return value * 1024 * 1024, nil
	case "G", "GB":
		return value * 1024 * 1024 * 1024, nil
	case "T", "TB":
		return value * 1024 * 1024 * 1024 * 1024, nil
	default:
		return 0, fmt.Errorf("unsupported unit: %s", unit)
	}
}
