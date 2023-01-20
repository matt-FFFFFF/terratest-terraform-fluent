package check

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCheckErrorF(t *testing.T) {
	err := newCheckErrorf("test %s", "error")
	assert.Equal(t, "test error", err.Error())
}

func TestNewCheckErrorAsError(t *testing.T) {
	f := func(e error) string {
		return e.Error()
	}
	ce := newCheckError("test error")
	assert.Equal(t, "test error", f(ce))
}

func TestNewCheckErrorNill(t *testing.T) {
	ce := func() *CheckError {
		return nil
	}()
	assert.Nil(t, ce)
}

func TestCheckErrorContains(t *testing.T) {
	ce := newCheckError("test error")
	ce.ErrorContains(t, "test")
}

func TestCheckErrorNotContains(t *testing.T) {
	ce := newCheckError("test error")
	ce.ErrorContains(t, "notcontained")
}

func TestCheckErrorIfNotFailWithNil(t *testing.T) {
	var ce *CheckError = nil
	ce.IfNotFail(t)
}
