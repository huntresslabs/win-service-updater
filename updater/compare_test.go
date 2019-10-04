package updater

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompare_CompareVersions(t *testing.T) {
	type versionTest struct {
		a        string
		b        string
		expected int
	}

	var versionTests = []versionTest{
		{"", "", A_EQUAL_TO_B},
		{"0.5.2", "0.6.2", A_LESS_THAN_B},
		{"0.5.2", "0.5.2", A_EQUAL_TO_B},
		{"1.0.0.1", "1.0.0.2", A_LESS_THAN_B},
		{"100.0.0.1", "200.0.0.2", A_LESS_THAN_B},
		{"0.0.0.5", "0.0.0.4", A_GREATER_THAN_B},
		{"10000.0.0.1", "20000.0.0.2", A_LESS_THAN_B},
		{"1.1.0beta1", "1.1.0beta2", A_LESS_THAN_B},
		{"1.1.1 beta 1", "1.1.0 beta 2", A_GREATER_THAN_B},
		{"1.1.1 beta1", "1.1.1beta 2", A_LESS_THAN_B},
		{"1.1.1 alpha", "1.1.1 beta", A_LESS_THAN_B},
		{"1.1.1.1.1", "1.1.1", A_GREATER_THAN_B},
		{"1.1.1.1.1", "1.1.2", A_LESS_THAN_B},
		{"1.01", "1.1", A_LESS_THAN_B},
	}

	for _, tt := range versionTests {
		// fmt.Printf("Compare \"%s\" and \"%s\"\n", tt.a, tt.b)
		actual := CompareVersions(tt.a, tt.b)
		assert.Equal(t, tt.expected, actual, fmt.Sprintf("a = %s; b = %s", tt.a, tt.b))
	}
}
