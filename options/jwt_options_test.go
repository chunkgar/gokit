package options

import (
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestNewJWTOptions(t *testing.T) {
	opts := NewJWTOptions()
	assert.Equal(t, "", opts.Issuer)
	assert.Equal(t, "", opts.Realm)
	assert.Equal(t, "", opts.Key)
	assert.Equal(t, 1*time.Hour, opts.Timeout)
	assert.Equal(t, 12*time.Hour, opts.MaxRefresh)
}

func TestJwtOptions_Validate(t *testing.T) {
	// 测试有效的选项
	validOpts := &JwtOptions{
		Issuer:     "test-issuer",
		Realm:      "test-realm",
		Key:        "test-key",
		Timeout:    1 * time.Hour,
		MaxRefresh: 12 * time.Hour,
	}
	errs := validOpts.Validate()
	assert.Empty(t, errs)

	// 测试无效的选项 - 缺少 Issuer
	invalidOpts1 := &JwtOptions{
		Issuer:     "",
		Realm:      "test-realm",
		Key:        "test-key",
		Timeout:    1 * time.Hour,
		MaxRefresh: 12 * time.Hour,
	}
	errs = invalidOpts1.Validate()
	assert.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "jwt.issuer is required")

	// 测试无效的选项 - 缺少 Key
	invalidOpts2 := &JwtOptions{
		Issuer:     "test-issuer",
		Realm:      "test-realm",
		Key:        "",
		Timeout:    1 * time.Hour,
		MaxRefresh: 12 * time.Hour,
	}
	errs = invalidOpts2.Validate()
	assert.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "jwt.key is required")

	// 测试无效的选项 - 同时缺少 Issuer 和 Key
	invalidOpts3 := &JwtOptions{
		Issuer:     "",
		Realm:      "test-realm",
		Key:        "",
		Timeout:    1 * time.Hour,
		MaxRefresh: 12 * time.Hour,
	}
	errs = invalidOpts3.Validate()
	assert.Len(t, errs, 2)
}

func TestJwtOptions_AddFlags(t *testing.T) {
	opts := NewJWTOptions()
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	opts.AddFlags(fs)

	// 验证所有标志是否已添加
	assert.NotNil(t, fs.Lookup("jwt.issuer"))
	assert.NotNil(t, fs.Lookup("jwt.realm"))
	assert.NotNil(t, fs.Lookup("jwt.key"))
	assert.NotNil(t, fs.Lookup("jwt.timeout"))
	assert.NotNil(t, fs.Lookup("jwt.max-refresh"))

	// 测试设置标志值
	err := fs.Parse([]string{
		"--jwt.issuer=test-issuer",
		"--jwt.realm=test-realm",
		"--jwt.key=test-key",
		"--jwt.timeout=2h",
		"--jwt.max-refresh=24h",
	})
	assert.NoError(t, err)

	assert.Equal(t, "test-issuer", opts.Issuer)
	assert.Equal(t, "test-realm", opts.Realm)
	assert.Equal(t, "test-key", opts.Key)
	assert.Equal(t, 2*time.Hour, opts.Timeout)
	assert.Equal(t, 24*time.Hour, opts.MaxRefresh)
}

func TestJwtOptions_NewJwt(t *testing.T) {
	opts := &JwtOptions{
		Issuer:     "test-issuer",
		Realm:      "test-realm",
		Key:        "test-key",
		Timeout:    1 * time.Hour,
		MaxRefresh: 12 * time.Hour,
	}

	jwt := opts.NewJwt()
	assert.Equal(t, "test-issuer", jwt.Issuer)
	assert.Equal(t, "test-realm", jwt.Realm)
	assert.Equal(t, []byte("test-key"), jwt.Key)
	assert.Equal(t, 1*time.Hour, jwt.Timeout)
	assert.Equal(t, 12*time.Hour, jwt.MaxRefresh)
}

func TestNewJWT(t *testing.T) {
	opts := &JwtOptions{
		Issuer:     "test-issuer",
		Realm:      "test-realm",
		Key:        "test-key",
		Timeout:    1 * time.Hour,
		MaxRefresh: 12 * time.Hour,
	}

	jwt := NewJWT(opts)
	assert.Equal(t, "test-issuer", jwt.Issuer)
	assert.Equal(t, "test-realm", jwt.Realm)
	assert.Equal(t, []byte("test-key"), jwt.Key)
	assert.Equal(t, 1*time.Hour, jwt.Timeout)
	assert.Equal(t, 12*time.Hour, jwt.MaxRefresh)
}

func TestJWT_GenerateJWT_ParseJWT(t *testing.T) {
	jwt := &JWT{
		Issuer:     "test-issuer",
		Realm:      "test-realm",
		Key:        []byte("test-key"),
		Timeout:    1 * time.Hour,
		MaxRefresh: 12 * time.Hour,
	}

	// 测试生成JWT
	userId := "user-123"
	role := "admin"
	token, expireTime, err := jwt.GenerateJWT(userId, role)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.True(t, expireTime.After(time.Now()))
	assert.True(t, expireTime.Before(time.Now().Add(1*time.Hour+5*time.Second)))

	// 测试解析JWT
	claims, err := jwt.ParseJWT(token)
	assert.NoError(t, err)
	assert.Equal(t, userId, claims.Subject)
	assert.Equal(t, role, claims.Role)
	assert.Equal(t, "test-issuer", claims.Issuer)

	// 测试解析无效的JWT
	_, err = jwt.ParseJWT("invalid.token.here")
	assert.Error(t, err)

	// 测试过期的token
	expiredJwt := &JWT{
		Issuer:     "test-issuer",
		Realm:      "test-realm",
		Key:        []byte("test-key"),
		Timeout:    -1 * time.Hour, // 过期的超时设置
		MaxRefresh: 12 * time.Hour,
	}
	expiredToken, _, err := expiredJwt.GenerateJWT(userId, role)
	assert.NoError(t, err)
	_, err = jwt.ParseJWT(expiredToken)
	assert.Error(t, err)
}
