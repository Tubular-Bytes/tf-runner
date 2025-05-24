package tofu_test

import (
	"testing"

	"github.com/Tubular-Bytes/tf-runner/pkg/tofu"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	stateURL := "http://localhost:8080/state"
	lockURL := "http://localhost:8080/lock"
	unlockURL := "http://localhost:8080/unlock"

	expected := `terraform {
  backend "http" {
    address = "http://localhost:8080/state"
    lock_address = "http://localhost:8080/lock"
    unlock_address = "http://localhost:8080/unlock"
  }
}`

	result, err := tofu.Render(stateURL, lockURL, unlockURL)
	require.NoError(t, err)

	require.Equal(t, expected, string(result))
}
