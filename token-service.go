package main

import (
	"encoding/hex"
	"errors"
	"time"

	// jwt "github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/v4"
)

const HmacSecret = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

type FileTokenParams struct {
	Bid           string `json:"bid"`
	FileSize      uint64 `json:"file_size"`
	FileHash      string `json:"file_hash"`
	Timeperiod    uint64 `json:"timeperiod"`
	SegmentsCount uint64 `json:"segments_count"`
}
type JWTFileTokenClaims struct {
	FileTokenParams
	jwt.StandardClaims
}

type JWTFileService struct{}

func (j JWTFileService) CreateFileToken(params FileTokenParams) (string, error) {
	hmacSecretBytes, err := hex.DecodeString(HmacSecret) //
	if err != nil {
		return "", err
	}
	claims := JWTFileTokenClaims{
		params,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "storage-node",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString(hmacSecretBytes)
	if err != nil {
		return "", err
	}

	return str, nil
}

func (j JWTFileService) ParseToken(tokenString string) (*FileTokenParams, error) {
	hmacSecretBytes, err := hex.DecodeString(HmacSecret) //
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(tokenString, &JWTFileTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return hmacSecretBytes, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTFileTokenClaims); ok && token.Valid {
		return &claims.FileTokenParams, nil
	} else {
		return nil, errors.New("unexpected error while parsing file token")
	}
}
