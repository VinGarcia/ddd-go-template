package rest

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultipartStream(t *testing.T) {
	t.Run("should stream normal readers correctly", func(t *testing.T) {
		stream, contentType, err := newMultipartStream(map[string]io.Reader{
			"item1": strings.NewReader(`{"fake":"json"}`),
			"item2": strings.NewReader(`================ other payload ==================`),
		})
		assert.Equal(t, nil, err)

		boundary := stream.multipartWriter.Boundary()
		assert.Equal(t, true, strings.Contains(contentType, `multipart/form-data;`))
		assert.Equal(t, true, strings.Contains(contentType, `boundary=`+boundary))

		// Reading he payload little by little
		// to make sure we are processing
		// the pauses correctly:
		var payload string
		buf := make([]byte, 10)
		var n int
		for err == nil {
			n, err = stream.Read(buf)
			payload += string(buf[:n])
		}
		assert.Equal(t, io.EOF, err)
		assert.Equal(t, 1, strings.Count(payload, `{"fake":"json"}`))
		assert.Equal(t, 1, strings.Count(payload, `================ other payload ==================`))
		assert.Equal(t, 3, strings.Count(payload, boundary))
		assert.Equal(t, 1, strings.Count(payload, `name="item1"`))
		assert.Equal(t, 1, strings.Count(payload, `name="item2"`))
		assert.Equal(t, 0, strings.Count(payload, `Content-Type:`))
	})

	t.Run("should stream items with Content-Type correctly", func(t *testing.T) {
		stream, contentType, err := newMultipartStream(map[string]io.Reader{
			"item1": MultipartItem(strings.NewReader(`{"fake":"json"}`), "application/json"),
			"item2": strings.NewReader(`================ other payload ==================`),
		})
		assert.Equal(t, nil, err)

		boundary := stream.multipartWriter.Boundary()
		assert.Equal(t, true, strings.Contains(contentType, `multipart/form-data;`))
		assert.Equal(t, true, strings.Contains(contentType, `boundary=`+boundary))

		// Reading he payload little by little
		// to make sure we are processing
		// the pauses correctly:
		var payload string
		buf := make([]byte, 10)
		var n int
		for err == nil {
			n, err = stream.Read(buf)
			payload += string(buf[:n])
		}
		assert.Equal(t, io.EOF, err)
		assert.Equal(t, 1, strings.Count(payload, `{"fake":"json"}`))
		assert.Equal(t, 1, strings.Count(payload, `================ other payload ==================`))
		assert.Equal(t, 3, strings.Count(payload, boundary))
		assert.Equal(t, 1, strings.Count(payload, `name="item1"`))
		assert.Equal(t, 1, strings.Count(payload, `name="item2"`))
		assert.Equal(t, 1, strings.Count(payload, `Content-Type:`))
		assert.Equal(t, 1, strings.Count(payload, `Content-Type: application/json`))
	})

	t.Run("should stream files correctly", func(t *testing.T) {
		stream, contentType, err := newMultipartStream(map[string]io.Reader{
			"item1": strings.NewReader(`{"fake":"json"}`),
			"item2": MultipartFile(strings.NewReader(`================ other payload ==================`), "fake-filename"),
		})
		assert.Equal(t, nil, err)

		boundary := stream.multipartWriter.Boundary()
		assert.Equal(t, true, strings.Contains(contentType, `multipart/form-data;`))
		assert.Equal(t, true, strings.Contains(contentType, `boundary=`+boundary))

		// Reading he payload little by little
		// to make sure we are processing
		// the pauses correctly:
		var payload string
		buf := make([]byte, 10)
		var n int
		for err == nil {
			n, err = stream.Read(buf)
			payload += string(buf[:n])
		}
		assert.Equal(t, io.EOF, err)
		assert.Equal(t, 1, strings.Count(payload, `{"fake":"json"}`))
		assert.Equal(t, 1, strings.Count(payload, `================ other payload ==================`))
		assert.Equal(t, 3, strings.Count(payload, boundary))
		assert.Equal(t, 1, strings.Count(payload, `name="item1"`))
		assert.Equal(t, 1, strings.Count(payload, `name="item2"`))
		assert.Equal(t, 1, strings.Count(payload, `filename="fake-filename"`))
		assert.Equal(t, 1, strings.Count(payload, `Content-Type:`))
		assert.Equal(t, 1, strings.Count(payload, `Content-Type: application/octet-stream`))
	})
}
