package updater

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgs(t *testing.T) {
	argv := []string{"/noerr", "-logfile=foo"}
	args := ParseArgs(argv)
	assert.True(t, args.Noerr)
	assert.Equal(t, args.Logfile, "foo")

	argv = []string{"/justcheck", "/outputinfo=foo"}
	args = ParseArgs(argv)
	assert.True(t, args.Justcheck)
	assert.Equal(t, args.Outputinfo, "foo")

	argv = []string{"/fromservice", "/quickcheck"}
	args = ParseArgs(argv)
	assert.True(t, args.Fromservice)
	assert.True(t, args.Quickcheck)
}
