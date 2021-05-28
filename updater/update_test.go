package updater

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdate_CompareVersions(t *testing.T) {
	type versionTest struct {
		a        string
		b        string
		expected int
	}

	var versionTests = []versionTest{
		{"0.5.2", "0.6.2", A_LESS_THAN_B},
		{"0.5.2", "0.5.2", A_EQUAL_TO_B},
		{"2.2.2", "2.2.2.1", A_LESS_THAN_B},
		{"3.3.3.1", "3.3.3", A_GREATER_THAN_B},
		{"1.0.0.1", "1.0.0.2", A_LESS_THAN_B},
		{"100.0.0.1", "200.0.0.2", A_LESS_THAN_B},
		{"0.0.0.5", "0.0.0.4", A_GREATER_THAN_B},
		{"10000.0.0.1", "20000.0.0.2", A_LESS_THAN_B},
	}

	for _, tt := range versionTests {
		actual := CompareVersions(tt.a, tt.b)
		assert.Equal(t, tt.expected, actual, fmt.Sprintf("a = %s; b = %s", tt.a, tt.b))
	}
}
