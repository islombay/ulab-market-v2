package status

import "net/http"

var (
	StatusProductArticulTooLong = Status{
		Message: "Articul too long",
		Code:    http.StatusRequestEntityTooLarge,
	}
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
	StatusProductNotFount = Status{
		Message:     "Product not found",
		Description: "product not found",
		Code:        http.StatusNotFound,
	}

	StatusProductPriceInvalid = Status{
		Message:     "Product price invalid",
		Description: "product price invalid",
		Code:        http.StatusBadRequest,
	}
)
