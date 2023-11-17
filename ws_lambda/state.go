package main

import (
	"context"
	"encoding/base64"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/fxamacker/cbor/v2"
	"github.com/pkg/errors"
)

const (
	ModeUnknown       = 0
	ModeAccess        = 1
	ModeConfiguration = 2
)

type Information struct {
	ConnectionId     string `cbor:"1,keyasint,omitempty"`
	ConnectMode      uint8  `cbor:"2,keyasint,omitempty"`
	ProfileName      string `cbor:"3,keyasint,omitempty"`
	ClientVersion    string `cbor:"4,keyasint,omitempty"`
	ApiId            string `cbor:"5,keyasint,omitempty"`
	AWSRegion        string `cbor:"6,keyasint,omitempty"`
	Duration         int    `cbor:"7,keyasint,omitempty"`
}

type SignedInformation struct {
	Information Information `cbor:"1,keyasint,omitempty"`
	Signature   []byte      `cbor:"2,keyasint,omitempty"`
}

func (i *Information) Sign(ctx context.Context) (*SignedInformation, error) {
	mode, err := cbor.CoreDetEncOptions().EncMode()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create serialization options")
	}

	encodedBlob, err := mode.Marshal(i)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal information structure")
	}

	signer := kms.NewFromConfig(awsConfig)
	sigOut, err := signer.Sign(ctx, &kms.SignInput{
		KeyId:            &kmsKeyArn,
		Message:          encodedBlob,
		SigningAlgorithm: types.SigningAlgorithmSpecEcdsaSha256,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign information structure")
	}

	signed := &SignedInformation{
		Information: *i, // Copy of information, not pointer to it
		Signature:   sigOut.Signature,
	}
	return signed, nil
}

func (si *SignedInformation) EncodeToString() (string, error) {
	mode, err := cbor.CoreDetEncOptions().EncMode()
	if err != nil {
		return "", err
	}

	encodedBlob, err := mode.Marshal(si)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(encodedBlob), nil
}
