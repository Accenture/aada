package main

type AssumptionResponse struct {
	Status string `json:"status"`

	Profile         string `json:"profile,omitempty"`
	Version         int    `json:"version,omitempty"`
	AccessKeyId     string `json:"access_key_id,omitempty"`
	SecretAccessKey string `json:"secret_access_key,omitempty"`
	SessionToken    string `json:"session_token,omitempty"`
	Expiration      string `json:"expiration,omitempty"`

	ProfileList map[string]string `json:"profile_list,omitempty"`
}
