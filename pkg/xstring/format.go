package xstring

import "strconv"

func FormatID(v int64) string {
	if v == 0 {
		return ""
	}

	return strconv.FormatInt(v, 10)
}

func ParseID(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}

	return strconv.ParseInt(s, 10, 64)
}
