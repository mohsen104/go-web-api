package services

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mohsen104/web-api/config"
	"github.com/mohsen104/web-api/pkg/logging"
	"github.com/mohsen104/web-api/pkg/service_errors"
)

type TokenService struct {
	logger logging.Logger
	cfg    *config.Config
}

type TokenDto struct {
	UserId    string
	FirstName string
	LastName  string
	Username  string
	Email     string
	Roles     []string
}

type TokenDetails struct {
	AccessToken            string `json:"accessToken"`
	RefreshToken           string `json:"refreshToken"`
	AccessTokenExpireTime  int64  `json:"accessTokenExpireTime"`
	RefreshTokenExpireTime int64  `json:"refreshTokenExpireTime"`
}

func NewTokenService(cfg *config.Config) *TokenService {
	logger := logging.NewLogger(cfg)
	return &TokenService{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *TokenService) GenerateToken(token *TokenDto) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AccessTokenExpireTime = int64(time.Now().Add(s.cfg.JWT.AccessTokenExpireDuration * time.Minute).Unix())
	td.RefreshTokenExpireTime = int64(time.Now().Add(s.cfg.JWT.RefreshTokenExpireDuration * time.Minute).Unix())

	atc := jwt.MapClaims{}

	atc["userId"] = token.UserId
	atc["firstName"] = token.FirstName
	atc["lastName"] = token.LastName
	atc["username"] = token.Username
	atc["email"] = token.Email
	atc["roles"] = token.Roles
	atc["exp"] = td.AccessTokenExpireTime

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atc)

	var err error
	td.AccessToken, err = at.SignedString([]byte(s.cfg.JWT.Secret))

	if err != nil {
		return nil, err
	}

	rtc := jwt.MapClaims{}

	rtc["userId"] = token.UserId
	rtc["exp"] = td.RefreshTokenExpireTime

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtc)

	td.RefreshToken, err = rt.SignedString([]byte(s.cfg.JWT.RefreshSecret))

	if err != nil {
		return nil, err
	}

	return td, nil
}

func (s *TokenService) VerifyToken(token string) (*jwt.Token, error) {
	at, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, &service_errors.ServiceError{EndUserMessages: service_errors.UnExpectedError}
		}
		return []byte(s.cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	return at, nil
}

func (s *TokenService) GetClaims(token string) (claimMap map[string]interface{}, err error) {
	claimMap = map[string]interface{}{}

	verifyToken, err := s.VerifyToken(token)
	if err != nil {
		return nil, err
	}
	claims, ok := verifyToken.Claims.(jwt.MapClaims)
	if ok && verifyToken.Valid {
		for k, v := range claims {
			claimMap[k] = v
		}
		return claimMap, nil
	}
	return nil, &service_errors.ServiceError{EndUserMessages: service_errors.ClaimsNotFound}
}
