package captcha

import (
	"context"

	"github.com/rs/zerolog/log"
)

type Resolver interface {
	CheckRecaptcha(context.Context, *string) error
}

type ResolverConfig struct {
	Secret    string  `mapstructure:"captcha-secret"`
	VerifyURL string  `mapstructure:"captcha-verify-url"`
	MinScore  float64 `mapstructure:"captcha-min-score"`
	Enable    bool    `mapstructure:"captcha"`
}

func NewCaptchaResolver(config ResolverConfig) Resolver {
	if config.Enable && config.Secret == "" {
		log.Error().Msg("Required Captcha secret not set")
	}
	return resolver{
		secret:        config.Secret,
		siteVerifyURL: config.VerifyURL,
		minScore:      config.MinScore,
		enable:        config.Enable,
	}
}
