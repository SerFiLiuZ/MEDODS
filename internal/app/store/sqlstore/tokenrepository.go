package sqlstore

import (
	"crypto/rand"
	"time"

	"github.com/SerFiLiuZ/MEDODS/internal/app/models"
	"github.com/SerFiLiuZ/MEDODS/internal/app/utils"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
)

type TokenRepository struct {
	store *Store
}

type claims struct {
	UserID string `json:"uid"`
	jwt.StandardClaims
}

func (r *TokenRepository) GetAccessRefreshTokens(guid string, jwtKey []byte) (*models.AuthUser, error) {
	logger := *utils.NewLogger()
	//logger.EnableDebug()

	logger.Debugf("TokenRepository: GetAccessRefreshTokens: start GetAccessRefreshTokens")

	var au models.AuthUser

	rtExpiration := time.Now().Add(5 * time.Hour)
	au.RtExpiration = rtExpiration

	logger.Debugf("TokenRepository: GetAccessRefreshTokens: au.RtExpiration: %v", au.RtExpiration)

	linkBytes := make([]byte, 32)
	_, err := rand.Read(linkBytes)
	if err != nil {
		return nil, err
	}

	for i, b := range linkBytes {
		linkBytes[i] = letters[b%byte(len(letters))]
	}

	rtClaims := &claims{
		UserID: guid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: rtExpiration.Unix(),
			Id:        string(linkBytes),
		},
	}

	logger.Debugf("TokenRepository: GetAccessRefreshTokens: start NewWithClaims")

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, rtClaims)
	rtSigned, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		logger.Errorf("TokenRepository: GetAccessRefreshTokens: SignedString: " + err.Error())
		return nil, err
	}

	au.RtSigned = rtSigned

	logger.Debugf("TokenRepository: GetAccessRefreshTokens: au.RtSigned: %v", au.RtSigned)
	logger.Debugf("TokenRepository: GetAccessRefreshTokens: au.RtSigned (size): %v", len([]byte(au.RtSigned)))

	atExpiration := time.Now().Add(5 * time.Minute)
	au.AtExpiration = atExpiration

	logger.Debugf("TokenRepository: GetAccessRefreshTokens: au.AtExpiration: %v", au.AtExpiration)

	atClaims := &claims{
		UserID: guid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: atExpiration.Unix(),
			Id:        string(linkBytes),
		},
	}

	logger.Debugf("TokenRepository: GetAccessRefreshTokens: start NewWithClaims")

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	atSigned, err := accessToken.SignedString(jwtKey)
	if err != nil {
		logger.Errorf("TokenRepository: GetAccessRefreshTokens: SignedString: " + err.Error())
		return nil, err
	}

	au.AtSigned = atSigned

	logger.Debugf("TokenRepository: GetAccessRefreshTokens: au.AtSigned: %v", au.AtSigned)

	return &au, nil
}

func (r *TokenRepository) GetHashedToken(token string) (string, error) {
	logger := *utils.NewLogger()
	logger.EnableDebug()

	logger.Debugf("TokenRepository: start GetHashedToken")
	logger.Debugf("TokenRepository: GetHashedToken: token: %v", token)
	logger.Debugf("TokenRepository: GetHashedToken: token (size): %v", len([]byte(token)))

	// TODO: эта функция работает некорректно
	hashed, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.MaxCost)
	if err != nil {
		logger.Errorf("TokenRepository: GetHashedToken: GenerateFromPassword: " + err.Error())
		return "", err
	}

	return string(hashed), nil
}

// func reverse(s string) string {
// 	runes := []rune(s)
// 	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
// 		runes[i], runes[j] = runes[j], runes[i]
// 	}
// 	return string(runes)
// }
