package helper

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
)

var (
	ErrInvalidEmail     = fmt.Errorf("invalid_email")
	ErrInvalidPhone     = fmt.Errorf("invalid_phone")
	ErrInvalidImageType = fmt.Errorf("invalid_image_type")
	ErrInvalidVideoType = fmt.Errorf("invalid_video_type")
	ErrInvalidIconType  = fmt.Errorf("invalid_icon_type")
)

var (
	validImageTypes = []string{
		"image/jpg",
		"image/webp",
		"image/png",
		"image/jpeg",
	}
	validVideoTypes = []string{
		"video/mp4",
		"video/x-msvideo",
		"video/quicktime",
	}
	validIconTypes = []string{
		"image/svg+xml",
		"text/plain",
		"text/plain; charset=utf-8",
		"image/png",
	}
)

func IsValidPhone(phone string) bool {
	if containsPlus := strings.Contains(phone, "+"); containsPlus {
		return false
	}

	// phone = strings.Replace(phone, "+", "", -1)
	r := regexp.MustCompile(`^998[0-9]{2}[0-9]{7}$`)
	return r.MatchString(phone)
}

func IsValidEmail(email string) bool {
	r := regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
	return r.MatchString(email)
}

func IsValidLogin(login string) bool {
	r := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{5,29}$`)
	return r.MatchString(login)
}

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func IsValidPassword(s string) bool {
	hasLetter := regexp.MustCompile(`[a-zA-Z]`)
	hasSpecial := regexp.MustCompile(`[:\-+_=%#@!^&<*,.]`)
	hasLength := len(s) >= 6
	return hasLength && hasSpecial.MatchString(s) && hasLetter.MatchString(s)
}

func IsValidImage(header *multipart.FileHeader) (bool, error, string) {
	file, err := header.Open()
	if err != nil {
		return false, err, ""
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return false, err, ""
	}

	contentType := http.DetectContentType(buffer)
	found := false
	for _, s := range validImageTypes {
		if s == contentType {
			found = true
			break
		}
	}

	if !found {
		return false, ErrInvalidImageType, contentType
	}
	return true, nil, contentType
}

func IsValidIcon(header *multipart.FileHeader) (bool, error, string) {
	file, err := header.Open()
	if err != nil {
		return false, err, ""
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return false, err, ""
	}

	contentType := http.DetectContentType(buffer)
	found := false
	for _, s := range validIconTypes {
		if s == contentType {
			found = true
			break
		}
	}

	if !found {
		return false, ErrInvalidIconType, contentType
	}
	return true, nil, contentType
}

func IsValidVideo(header *multipart.FileHeader) (bool, error, string) {
	file, err := header.Open()
	if err != nil {
		return false, err, ""
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return false, err, ""
	}

	contentType := http.DetectContentType(buffer)
	for _, validType := range validVideoTypes {
		if contentType == validType {
			return true, nil, contentType
		}
	}

	return false, ErrInvalidVideoType, contentType
}
