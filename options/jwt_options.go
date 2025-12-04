package options

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/pflag"
)

type JwtOptions struct {
	Issuer     string        `json:"issuer"      mapstructure:"issuer"`
	Realm      string        `json:"realm"       mapstructure:"realm"`
	Key        string        `json:"key"         mapstructure:"key"`
	Timeout    time.Duration `json:"timeout"     mapstructure:"timeout"`
	MaxRefresh time.Duration `json:"max-refresh" mapstructure:"max-refresh"`
}

func NewJWTOptions() *JwtOptions {
	return &JwtOptions{
		Issuer:     "",
		Realm:      "",
		Key:        "",
		Timeout:    1 * time.Hour,
		MaxRefresh: 12 * time.Hour,
	}
}

func (o *JwtOptions) Validate() []error {
	var errs []error

	if o.Issuer == "" {
		errs = append(errs, fmt.Errorf("jwt.issuer is required"))
	}

	if o.Key == "" {
		errs = append(errs, fmt.Errorf("jwt.key is required"))
	}

	return errs
}

func (o *JwtOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Issuer, "jwt.issuer", o.Issuer, "Issuer name.")
	fs.StringVar(&o.Realm, "jwt.realm", o.Realm, "Realm name to display to the user.")
	fs.StringVar(&o.Key, "jwt.key", o.Key, "Private key used to sign jwt token.")
	fs.DurationVar(&o.Timeout, "jwt.timeout", o.Timeout, "JWT token timeout.")

	fs.DurationVar(&o.MaxRefresh, "jwt.max-refresh", o.MaxRefresh, ""+
		"This field allows clients to refresh their token until MaxRefresh has passed.")
}

func (o *JwtOptions) NewJwt() *JWT {
	return &JWT{
		Issuer:     o.Issuer,
		Realm:      o.Realm,
		Key:        []byte(o.Key),
		Timeout:    o.Timeout,
		MaxRefresh: o.MaxRefresh,
	}
}

type JWTClaims struct {
	jwt.RegisteredClaims // 用户ID在 Subject 字段中

	Role string `json:"role"` // 可能没有，可能 "anon" 或 "service_role"
}

type JWT struct {
	Issuer     string
	Realm      string
	Key        []byte
	Timeout    time.Duration
	MaxRefresh time.Duration
}

func NewJWT(opts *JwtOptions) *JWT {
	return &JWT{
		Issuer:     opts.Issuer,
		Realm:      opts.Realm,
		Key:        []byte(opts.Key),
		Timeout:    opts.Timeout,
		MaxRefresh: opts.MaxRefresh,
	}
}

// 生成 JWT token
func (j *JWT) GenerateJWT(userId string, role string) (string, time.Time, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(j.Timeout)

	claims := JWTClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			Issuer:    j.Issuer,
			Subject:   userId,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(j.Key)

	return token, expireTime, err
}

// 解析 JWT token
func (j *JWT) ParseJWT(token string) (*JWTClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.Key, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := tokenClaims.Claims.(*JWTClaims); ok && tokenClaims.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
