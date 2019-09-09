package updater

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWYS(t *testing.T) {
	wys := ParseWys("../test_files/wys_uncompressed.bin")
	assert.Contains(t, wys.UpdateFileSite, "127.0.0.1")
}
