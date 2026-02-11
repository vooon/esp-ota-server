package server

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thanhpk/randstr"
)

func TestGetEspHeader(t *testing.T) {
	testCases := []struct {
		name       string
		header     http.Header
		key        string
		expected   []string
		expectedOK bool
	}{
		{
			name: "esp8266",
			header: http.Header{
				"X-Esp8266-Mode": []string{"sketch"},
			},
			key:        "mode",
			expected:   []string{"sketch"},
			expectedOK: true,
		},
		{
			name: "esp32",
			header: http.Header{
				"X-Esp32-Mode": []string{"sketch"},
			},
			key:        "mode",
			expected:   []string{"sketch"},
			expectedOK: true,
		},
		{
			name:       "not-found",
			header:     http.Header{},
			key:        "mode",
			expected:   nil,
			expectedOK: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, ok := getEspHeader(tc.header, tc.key)
			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expected, val)
		})
	}
}

func TestGetBinaryFile(t *testing.T) {
	const (
		project  = "project-a"
		filename = "firmware.bin"
	)
	content := randstr.String(128)

	dataDir := t.TempDir()
	t.Chdir(dataDir)
	projectDir := filepath.Join(dataDir, project)
	require.NoError(t, os.MkdirAll(projectDir, 0o755))

	binPath := filepath.Join(projectDir, filename)
	require.NoError(t, os.WriteFile(binPath, []byte(content), 0o600))

	sum := md5.Sum([]byte(content))
	md5sum := hex.EncodeToString(sum[:])

	s := server{
		config: Config{
			DataDirPath: ".",
		},
	}

	testCases := []struct {
		name          string
		headers       http.Header
		expectedCode  int
		expectedBody  string
		expectMD5     bool
		expectSHA512  bool
	}{
		{
			name:         "missing-mode-header",
			headers:      http.Header{},
			expectedCode: http.StatusBadRequest,
			expectedBody: "bad request",
		},
		{
			name: "not-modified-when-md5-matches",
			headers: http.Header{
				"X-Esp8266-Mode":       []string{"sketch"},
				"X-Esp8266-Sketch-Md5": []string{md5sum},
			},
			expectedCode: http.StatusNotModified,
			expectedBody: "",
			expectMD5:    true,
			expectSHA512: true,
		},
		{
			name: "send-file-when-md5-differs",
			headers: http.Header{
				"X-Esp8266-Mode":       []string{"sketch"},
				"X-Esp8266-Sketch-Md5": []string{"deadbeef"},
			},
			expectedCode: http.StatusOK,
			expectedBody: content,
			expectMD5:    true,
			expectSHA512: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/bin/"+project+"/"+filename, nil)
			req.Header = tc.headers.Clone()
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath("/bin/:project/:file")
			c.SetPathValues(echo.PathValues{
				{Name: "project", Value: project},
				{Name: "file", Value: filename},
			})

			err := s.getBinaryFile(c)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedBody, rec.Body.String())
			if tc.expectMD5 {
				assert.Equal(t, []string{md5sum}, rec.Header()["x-MD5"])
			}
			if tc.expectSHA512 {
				assert.NotEmpty(t, rec.Header().Get("x-SHA512"))
			}
		})
	}
}
