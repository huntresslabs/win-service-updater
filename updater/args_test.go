package updater

import (
	"testing"

	"src/github.com/stretchr/testify/assert"
)

func TestArgs(t *testing.T) {
	argv := []string{"/noerr", "-logfile=foo"}
	args := ParseArgs(argv)
	assert.True(t, args.Noerr)
	assert.Equal(t, args.Logfile, "foo")
}
