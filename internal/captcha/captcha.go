package captcha

type Resolver interface {
	CheckRecaptcha(string) error
}

type resolver struct {
	secret string
}

func NewCaptchaResolver(secret string) Resolver {
	return resolver{
		secret: secret,
	}
}
