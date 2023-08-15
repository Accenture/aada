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
}

type SignedInformation struct {
	Information Information `cbor:"1,keyasint,omitempty"`
	Signature   []byte      `cbor:"2,keyasint,omitempty"`
}

// DecodeFromString unpacks the base64 encoded state and unmarshals into the struct
func (si *SignedInformation) DecodeFromString(str string) error {
	block, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		return errors.Wrap(err, "failed to decode state")
	}

	err = cbor.Unmarshal(block, si)
	if err != nil {
		return errors.Wrap(err, "failed to unpack data")
	}

	return nil
}

// Validate asks KMS to compute/validate the signature on the state
func (si *SignedInformation) Validate(ctx context.Context) error {
	mode, err := cbor.CoreDetEncOptions().EncMode()
	if err != nil {
		return errors.Wrap(err, "failed to create serialization options")
	}

	encodedBlob, err := mode.Marshal(&si.Information)
	if err != nil {
		return errors.Wrap(err, "failed to marshal information structure")
	}

	signer := kms.NewFromConfig(awsConfig)
	vOut, err := signer.Verify(ctx, &kms.VerifyInput{
		KeyId:            &kmsKeyArn,
		Message:          encodedBlob,
		Signature:        si.Signature,
		SigningAlgorithm: types.SigningAlgorithmSpecEcdsaSha256,
	})
	if err != nil {
		return errors.Wrap(err, "failed to call key verification")
	}

	if vOut.SignatureValid {
		return nil // Correct case
	}

	return errors.New("signature was not valid")
}
