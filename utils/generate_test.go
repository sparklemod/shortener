package utils

import (
	"fmt"
	"testing"
)

func TestGenerateBase58String(t *testing.T) {
	t.Run("returns an base58 string", func(t *testing.T) {
		testCases := []struct {
			length   int
			expected int
		}{
			{length: 8, expected: 8},
			{length: 0, expected: 0},
		}

		for _, testCase := range testCases {
			got, err := GenerateBase58String(testCase.length)
			if err != nil {
				t.Error(err)
			}
			if len(got) != testCase.expected {
				t.Errorf("got %v, expected %v", got, testCase.expected)
			} else {
				fmt.Printf("output: \"%v\", %v\n", got, len(got))
			}
		}
	})

	t.Run("is randomness", func(t *testing.T) {
		results := make(map[string]bool)
		for i := 0; i < 100; i++ {
			s, err := GenerateBase58String(1)
			if err != nil {
				t.Errorf("unexpected error generating string: %v", err)
				continue
			}

			if results[s] {
				t.Errorf("duplicate string: \"%s\", step: %v", s, i)
			}
		}
	})
}
