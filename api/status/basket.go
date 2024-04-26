package status

import "net/http"

var (
	StatusProductQuantityTooMany = Status{
		Message: "Product quantity too many",
		Code:    http.StatusRequestEntityTooLarge,
	}
)
