package wrap

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/wrap"
)

func AssertEqualWraps(t *testing.T, expected, actual wrap.Wrap) {
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Weight, actual.Weight)
	assert.Equal(t, expected.Cost, actual.Cost)
}
