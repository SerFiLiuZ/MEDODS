package sqlstore

import (
	"crypto/rand"
	"time"

	"github.com/SerFiLiuZ/MEDODS/internal/app/models"
	"github.com/dgrijalva/jwt-go"
)

const (
	letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
)

type TokenRepository struct {
	store *Store
}

type сlaims struct {
	UserID string `json:"uid"`
	jwt.StandardClaims
}

func (r *TokenRepository) GetAccessRefreshTokens(guid, jwtKey string) (*models.AuthUser, error) {
	var au models.AuthUser

	rtExpiration := time.Now().Add(5 * time.Hour)
	au.RtExpiration = rtExpiration

	linkBytes := make([]byte, 32)
	_, err := rand.Read(linkBytes)
	if err != nil {
		return nil, err
	}
	for i, b := range linkBytes {
		linkBytes[i] = letters[b%byte(len(letters))]
	}
	rtClaims := &сlaims{
		UserID: guid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: rtExpiration.Unix(),
			Id:        string(linkBytes),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, rtClaims)
	rtSigned, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	au.RtSigned = rtSigned

	atExpiration := time.Now().Add(5 * time.Minute)
	au.AtExpiration = atExpiration

	atClaims := &сlaims{
		UserID: guid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: atExpiration.Unix(),
			Id:        string(linkBytes),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	atSigned, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	au.AtSigned = atSigned

	return &au, nil
}
