package auth_lib

import (
	"app/pkg/helper"
	"app/pkg/smtp"
	"app/storage"
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrUnimplemented = fmt.Errorf("unimplemented")
)

func SendVerificationCode(source, sourceType string, codeLength int,
	cache storage.CacheInterface,
	smtp smtp.SMTPInterface,
	codeExpirationTime time.Duration,
) (error, string) {
	oneTimeCode := GenerateRandomPassword(codeLength)
	oneTimeCodeExpireTime := time.Now().Add(codeExpirationTime)
	err := cache.Code().SetCode(
		context.Background(),
		source,
		oneTimeCode,
		oneTimeCodeExpireTime,
	)
	if err != nil {
		return err, ""
	}

	if sourceType == VerificationEmail {
		if err := smtp.Email().SendVerificationCode(source, oneTimeCode); err != nil {
			if errors.Is(err, helper.ErrInvalidEmail) {
				return helper.ErrInvalidEmail, ""
			}
			return err, ""
		}
	} else {
		// ignore
		// return ErrUnimplemented
	}
	return nil, oneTimeCode
}
