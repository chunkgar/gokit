package options

import (
	"fmt"
	"time"

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
