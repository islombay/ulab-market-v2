package status

import "net/http"

var (
	StatusBrandNotFound = Status{
		Message:     "Brand not found",
		Description: "brand not found",
		Code:        http.StatusNotFound,
	}

	StatusBranchNotFound = Status{
		Message:     "Branch not found",
		Description: "branch not found",
		Code:        http.StatusNotFound,
	}
)
