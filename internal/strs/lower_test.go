package strs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLowerName(t *testing.T) {
	cases := []struct {
		Args    string
		Expect  string
		Message string
	}{
		{
			"CaseToCase",
			"case_to_case",
			"failed 1",
		},
		{
			"CaseTOCase",
			"case_to_case",
			"failed 2",
		},
		{
			"caseT2CaseC",
			"case_t2_case_c",
			"failed 3",
		},
		{
			"case_t2_case_c100",
			"case_t2_case_c100",
			"failed 4",
		},
	}

	for _, _case := range cases {
		result := LowerName(_case.Args)
		assert.Equal(t, _case.Expect, result, _case.Message)
	}
}
