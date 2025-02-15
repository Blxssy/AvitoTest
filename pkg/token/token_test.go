package token

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewToken(t *testing.T) {
	cfg := TokenConfig{
		TokenKey: "testkey",
		TokenTTL: time.Hour,
	}
	tokenGen := NewTokenGen(cfg)

	token, err := tokenGen.NewToken("testuser")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestParseToken(t *testing.T) {
	cfg := TokenConfig{
		TokenKey: "testkey",
		TokenTTL: time.Hour,
	}
	tokenGen := NewTokenGen(cfg)

	token, err := tokenGen.NewToken("testuser")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	username, err := tokenGen.ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", username)
}

func TestParseTokenInvalid(t *testing.T) {
	cfg := TokenConfig{
		TokenKey: "testkey",
		TokenTTL: time.Hour,
	}
	tokenGen := NewTokenGen(cfg)

	_, err := tokenGen.ParseToken("invalidtoken")
	assert.Error(t, err)
}
