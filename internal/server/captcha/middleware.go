package captcha

import (
	"fmt"
	"net/http"
)

func VerificationMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the Captcha response token from header request field 'g-recaptcha-response'.
			if err := checkRecaptcha(secret, r.Header.Get("g-recaptcha-response")); err != nil {
				http.Error(w, fmt.Sprintf("Unauthorized: %s", err.Error()), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
