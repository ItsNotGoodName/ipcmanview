package ptz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getSeq(t *testing.T) {
	actual := getSeq("s2c7dqQZa29bc00a2728116234", 0)
	assert.Equal(t, 860883456, actual)

	actual = getSeq("s2c7dqQZa29bc00a2728116234", 1)
	assert.Equal(t, 860883457, actual)
}
