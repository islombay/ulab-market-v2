package helper

import (
	"strconv"
	"strings"
)

func ReplaceQueryParams(namedQuery string, params map[string]interface{}) (string, []interface{}) {
	var (
		i    = 1
		args []interface{}
	)

	for k, v := range params {
		if k != "" {
			oldSize := len(namedQuery)
			namedQuery = strings.ReplaceAll(namedQuery, "@"+k, "$"+strconv.Itoa(i))

			if oldSize != len(namedQuery) {
				args = append(args, v)
				i++
			}
		}
	}

	return namedQuery, args
}
