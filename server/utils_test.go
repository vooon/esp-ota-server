package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
