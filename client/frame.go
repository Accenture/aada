package main

import "encoding/json"

type Frame struct {
	Nonce           string            `json:"nonce,omitempty"`
	Profile         string            `json:"profile,omitempty"`
	State           string            `json:"state,omitempty"`
	Context         string            `json:"context,omitempty"`
	Status          string            `json:"status,omitempty"`
	Mode            string            `json:"mode,omitempty"`
	Version         int               `json:"version,omitempty"`
	ClientVersion   string            `json:"client_version,omitempty"`
	AccessKeyId     string            `json:"access_key_id,omitempty"`
	SecretAccessKey string            `json:"secret_access_key,omitempty"`
	SessionToken    string            `json:"session_token,omitempty"`
	Expiration      string            `json:"expiration,omitempty"`
	Message         string            `json:"message,omitempty"`
	ProfileList     map[string]string `json:"profile_list,omitempty"`
}

type CredentialStruct struct {
	Version         int    `json:"Version,omitempty"`
	AccessKeyId     string `json:"AccessKeyId,omitempty"`
	SecretAccessKey string `json:"SecretAccessKey,omitempty"`
	SessionToken    string `json:"SessionToken,omitempty"`
	Expiration      string `json:"Expiration,omitempty"`
}

func (frame *Frame) ToCredentialString() string {
	cs := &CredentialStruct{
		Version:         frame.Version,
		AccessKeyId:     frame.AccessKeyId,
		SecretAccessKey: frame.SecretAccessKey,
		SessionToken:    frame.SessionToken,
		Expiration:      frame.Expiration,
	}
	raw, _ := json.Marshal(cs)
	return string(raw)
}
