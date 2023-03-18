package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandInt(t *testing.T) {
	assert := assert.New(t)
	str := GenerateRandInt(10)
	assert.Len(str, 10)
}
