package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSafeFilePath(t *testing.T) {
	base := "/data"

	testCases := []struct {
		name        string
		parts       []string
		expected    string
		expectError bool
	}{
		{
			name:     "valid-path",
			parts:    []string{"project-a", "firmware.bin"},
			expected: "/data/project-a/firmware.bin",
		},
		{
			name:     "nested-valid",
			parts:    []string{"project-a", "sub", "fw.bin"},
			expected: "/data/project-a/sub/fw.bin",
		},
		{
			name:        "dotdot-in-project",
			parts:       []string{"..", "passwd"},
			expectError: true,
		},
		{
			name:        "dotdot-in-filename",
			parts:       []string{"project-a", "../../etc/passwd"},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := safeFilePath(base, tc.parts...)
			if tc.expectError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, p)
		})
	}
}

func TestParseVersionHeader(t *testing.T) {

	testCases := []struct {
		name           string
		hdrs           []string
		expected       map[string]string
		expectedErrMsg string
	}{
		{"nil", nil, map[string]string{}, ""},
		{"kv", []string{"fw:1.2.3 hv:1.0"}, map[string]string{"fw": "1.2.3", "hv": "1.0"}, ""},
		{"tasmota", []string{"13.3.0.3(tasmota-4M)"}, map[string]string{"ver": "13.3.0.3(tasmota-4M)"}, ""}, // see #35
		{"multiple-kv", []string{"fw:1.2.3", "hv:1.0"}, map[string]string{"fw": "1.2.3", "hv": "1.0"}, ""},  // not actually sent by esp, just check corner case
		{"malicious-space", []string{"foo:1.2.3 bar"}, nil, "failed to parse version:"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)

			vmap, err := parseVersionHeader(tc.hdrs)
			if tc.expectedErrMsg == "" {
				assert.NoError(err)
				assert.Equal(tc.expected, vmap)
			} else {
				assert.ErrorContains(err, tc.expectedErrMsg)
			}

		})
	}
}
