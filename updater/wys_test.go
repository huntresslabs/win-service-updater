package updater

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWYS(t *testing.T) {
	var args Args
	wys, err := ParseWys("../test_files/compressed.wys", args)
	assert.Nil(t, err)
	assert.Contains(t, wys.UpdateFileSite, "127.0.0.1")
}
