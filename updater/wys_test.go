package updater

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWYS(t *testing.T) {
	var args Args
	wys, err := ParseWYS("../test_files/widgetX.1.0.1.wys", args)
	assert.Nil(t, err)
	assert.Contains(t, wys.UpdateFileSite[0], "127.0.0.1")
}
