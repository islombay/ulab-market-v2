package status

import "net/http"

var (
	StatusProductMainImageMaxSizeExceed = Status{
		Message:     "Main image size exceeds the limit",
		Description: "main image size is big",
		Code:        http.StatusRequestEntityTooLarge,
	}
	StatusProductVideoMaxSizeExceed = Status{
		Message:     "Product video size exceeds the limit",
		Description: "video size is big",
		Code:        http.StatusRequestEntityTooLarge,
	}

	StatusProductPhotoMaxCount = Status{
		Message:     "Product photo count is too many",
		Description: "product photo count is too many",
		Code:        http.StatusRequestEntityTooLarge,
	}

	StatusProductVideoMaxCount = Status{
		Message:     "Product video count is too many",
		Description: "product video count is too many",
		Code:        http.StatusRequestEntityTooLarge,
	}

	StatusProductStatusInvalid = Status{
		Message:     "Product status invalid",
		Description: "product status invalid",
		Code:        http.StatusBadRequest,
	}
)
