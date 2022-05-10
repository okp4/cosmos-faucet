package captcha

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const siteVerifyURL = "https://www.google.com/recaptcha/api/siteverify"

type siteVerifyResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

type siteVerifyRequest struct {
	RecaptchaResponse string `json:"g-recaptcha-response"`
}

func checkRecaptcha(secret, response string) error {
	req, err := http.NewRequest(http.MethodPost, siteVerifyURL, nil)
	if err != nil {
		log.Error().Msgf("Error while creating Captcha verification request: %s", err.Error())
		return fmt.Errorf("error while creating Captcha verification request: %s", err.Error())
	}

	q := req.URL.Query()
	q.Add("secret", secret)
	q.Add("response", response)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Msgf("Error while requesting Captcha verification: %s", err.Error())
		return fmt.Errorf("error while requesting Captcha verification: %s", err.Error())
	}
	defer resp.Body.Close()

	var body siteVerifyResponse
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		log.Error().Msgf("Error while decoding Captcha verification response: %s", err.Error())
		return fmt.Errorf("error while decoding Captcha verification response: %s", err.Error())
	}

	// If success false, Captcha verification KO.
	if !body.Success {
		log.Debug().Msg("Captcha verification failed")
		return errors.New("captcha verification failed")
	}

	// If score is too low, verification KO.
	if body.Score > 0.5 {
		log.Debug().Msg("Captcha verification failed: score is too low")
		return errors.New("captcha verification failed: score is too low")
	}

	return nil
}
