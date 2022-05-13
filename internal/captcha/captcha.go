package captcha

import "github.com/rs/zerolog/log"

type Resolver interface {
	CheckRecaptcha(string) error
}

type resolver struct {
	secret string
}

func NewCaptchaResolver(secret string) Resolver {
	if secret == "" {
		log.Fatal().Msg("Required Captcha secret not set")
	}
	return resolver{
		secret: secret,
	}
}
