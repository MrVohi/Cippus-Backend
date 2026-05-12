package services

import (
"net/url"
"net/http"
"encoding/json"
)

type RecaptchaResponse struct{
	Success bool `json:"success"`
}

func (s *AuthService) VerifyRecaptcha(token string) bool {	
	resp, err := http.PostForm(
		"https://www.google.com/recaptcha/api/siteverify",
		url.Values{ "secret": {s.recaptchaSecret}, "response": {token}},
	)
	if err != nil {
		return false
	}

	body := RecaptchaResponse{}
	json.NewDecoder(resp.Body).Decode(&body)
	return body.Success
}