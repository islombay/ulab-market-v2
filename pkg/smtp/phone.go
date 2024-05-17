package smtp

import (
	"app/config"
	"app/pkg/helper"
	"app/pkg/logs"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type phoneSMTPService struct {
	cfg *config.SMTPEmailConfig
	log logs.LoggerInterface
}

func NewPhoneSMTPService(cfg *config.SMTPEmailConfig, log logs.LoggerInterface) *phoneSMTPService {
	return &phoneSMTPService{
		cfg: cfg,
		log: log,
	}
}

func (se *phoneSMTPService) SendVerificationCode(to, code string) error {
	if !helper.IsValidPhone(to) {
		return helper.ErrInvalidPhone
	}

	url := "https://lms-back.nvrbckdown.uz/lms/api/v1/send-otp"
	payload := map[string]interface{}{
		"code":         code,
		"phone_number": to,
		"sender_id":    os.Getenv("otp_id"),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		se.log.Error("Failed to marshal JSON: %v", logs.Error(err))
		return err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		se.log.Error("Failed to create request: %v", logs.Error(err))
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		se.log.Error("Failed to make request: %v", logs.Error(err))
		return err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		se.log.Error("Failed to read response body: %v", logs.Error(err))
		return err
	}

	// Print the response status and body
	if response.StatusCode != http.StatusOK {
		se.log.Error("could not send otp code", logs.String("body", string(body)))
		return err
	}

	return nil
}
