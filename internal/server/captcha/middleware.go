package captcha

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func VerificationMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the Captcha response token from default request body field 'g-recaptcha-response'.
			bodyBytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, fmt.Sprintf("Unauthorized: %s", err.Error()), http.StatusUnauthorized)
				return
			}

			// Unmarshal body into struct.
			var body siteVerifyRequest
			if err := json.Unmarshal(bodyBytes, &body); err != nil {
				http.Error(w, fmt.Sprintf("Unauthorized: %s", err.Error()), http.StatusUnauthorized)
				return
			}

			// Restore request body to read more than once.
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

			// Check and verify the Captcha response token.
			if err := checkRecaptcha(secret, body.RecaptchaResponse); err != nil {
				http.Error(w, fmt.Sprintf("Unauthorized: %s", err.Error()), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
