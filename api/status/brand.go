package status

import "net/http"

var (
	StatusBrandNotFound = Status{
		Message:     "Brand not found",
		Description: "brand not found",
		Code:        http.StatusNotFound,
	}
)
