package captcha

import (
	"context"

	"github.com/rs/zerolog/log"
)

type Resolver interface {
	CheckRecaptcha(context.Context, string) error
}

func NewCaptchaResolver(secret, verifyURL string, minScore float64) Resolver {
	if secret == "" {
		log.Error().Msg("Required Captcha secret not set")
	}
	return resolver{
		secret:        secret,
		siteVerifyURL: verifyURL,
		minScore:      minScore,
	}
}
