package captcha

import (
	ctx "context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type siteVerifyResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

type resolver struct {
	secret        string
	siteVerifyURL string
	minScore      float64
	enable        bool
}

func (c resolver) CheckRecaptcha(ctx ctx.Context, response *string) error {
	if !c.enable {
		return nil
	}

	if response == nil {
		log.Debug().Msg("No captcha token specified")
		return fmt.Errorf("no captcha token specified")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.siteVerifyURL, nil)
	if err != nil {
		log.Error().Err(err).Msgf("Error while creating Captcha verification request: %s", err.Error())
		return fmt.Errorf("error while creating Captcha verification request: %w", err)
	}

	q := req.URL.Query()
	q.Add("secret", c.secret)
	q.Add("response", *response)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msgf("Error while requesting Captcha verification: %s", err.Error())
		return fmt.Errorf("error while requesting Captcha verification: %w", err)
	}
	defer resp.Body.Close()

	var body siteVerifyResponse
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		log.Error().Err(err).Msgf("Error while decoding Captcha verification response: %s", err.Error())
		return fmt.Errorf("error while decoding Captcha verification response: %w", err)
	}

	// If success false, Captcha verification KO.
	if !body.Success {
		log.Debug().Msg("Captcha verification failed")
		return errors.New("captcha verification failed")
	}

	// If score is too low, verification KO.
	if body.Score < c.minScore {
		log.Debug().Float64("score", body.Score).Msg("Captcha verification failed: score is too low")
		return errors.New("captcha verification failed: score is too low")
	}

	return nil
}
