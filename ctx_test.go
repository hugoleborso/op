package op

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContext_PathParam(t *testing.T) {
	t.Run("can read path param", func(t *testing.T) {
		t.Skip("TODO: coming in go1.22")
		r := httptest.NewRequest("GET", "http://example.com/foo/123", nil)

		c := NewContext[any](r, readOptions{})

		param := c.PathParam("id")
		require.NotEmpty(t, param)
		require.Equal(t, param, "123")
	})
}

func TestContext_QueryParam(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.com/foo/123?id=456&other=hello", nil)

	c := NewContext[any](r, readOptions{})

	param := c.QueryParam("id")
	require.NotEmpty(t, param)
	require.Equal(t, param, "456")

	param = c.QueryParam("other")
	require.NotEmpty(t, param)
	require.Equal(t, param, "hello")

	param = c.QueryParam("notfound")
	require.Empty(t, param)
}

func TestContext_QueryParams(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.com/foo/123?id=456&other=hello", nil)

	c := NewContext[any](r, readOptions{})

	params := c.QueryParams()
	require.NotEmpty(t, params)
	require.Equal(t, params["id"], "456")
	require.Equal(t, params["other"], "hello")
	require.Empty(t, params["notfound"])
}

type testStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type testStructNormalizable struct {
	Name string `json:"name" validate:"required,min=3,max=10"`
	Age  int    `json:"age" validate:"min=18"`
}

func (b *testStructNormalizable) Normalize() error {
	b.Name = "normalized " + b.Name
	b.Age *= 2
	return nil
}

type testStructNormalizableWithError struct {
	Name string `json:"name" validate:"required,min=3,max=10"`
	Age  int    `json:"age" validate:"min=18"`
}

func (b *testStructNormalizableWithError) Normalize() error {
	return errors.New("error")
}

func TestContext_Body(t *testing.T) {
	t.Run("can read JSON body", func(t *testing.T) {
		// Create new Reader
		a := strings.NewReader(`{"name":"John","age":30}`)

		// Test an http request
		r := httptest.NewRequest("GET", "http://example.com/foo", a)

		c := NewContext[testStruct](r, readOptions{})

		body, err := c.Body()
		require.NoError(t, err)
		require.Equal(t, body.Name, "John")
		require.Equal(t, body.Age, 30)
	})

	t.Run("can read JSON body twice", func(t *testing.T) {
		a := strings.NewReader(`{"name":"John","age":30}`)

		r := httptest.NewRequest("GET", "http://example.com/foo", a)

		c := NewContext[testStruct](r, readOptions{})

		body, err := c.Body()
		require.NoError(t, err)
		require.Equal(t, body.Name, "John")
		require.Equal(t, body.Age, 30)

		body, err = c.Body()
		require.NoError(t, err)
		require.Equal(t, body.Name, "John")
		require.Equal(t, body.Age, 30)
	})

	t.Run("can read and validate valid JSON body", func(t *testing.T) {
		type testStruct struct {
			Name string `json:"name" validate:"required,min=3,max=10"`
			Age  int    `json:"age" validate:"min=18"`
		}

		reqBody := strings.NewReader(`{"name":"John","age":30}`)
		c := NewContext[testStruct](
			httptest.NewRequest("GET", "http://example.com/foo", reqBody),
			readOptions{})

		body, err := c.Body()
		require.NoError(t, err)
		require.Equal(t, body.Name, "John")
		require.Equal(t, body.Age, 30)
	})

	t.Run("can read and validate invalid JSON body", func(t *testing.T) {
		type testStruct struct {
			Name string `json:"name" validate:"required,min=3,max=10"`
			Age  int    `json:"age" validate:"min=18"`
		}

		reqBody := strings.NewReader(`{"name":"VeryLongName","age":12}`)
		c := NewContext[testStruct](
			httptest.NewRequest("GET", "http://example.com/foo", reqBody),
			readOptions{})

		body, err := c.Body()
		require.Error(t, err)
		require.Equal(t, body.Name, "VeryLongName")
		require.Equal(t, body.Age, 12)
	})

	t.Run("can normalize JSON body with custom method", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name":"John","age":30}`)
		c := NewContext[testStructNormalizable](
			httptest.NewRequest("GET", "http://example.com/foo", reqBody),
			readOptions{})

		body, err := c.Body()
		require.NoError(t, err)
		require.Equal(t, body.Name, "normalized John")
		require.Equal(t, body.Age, 60)
	})

	t.Run("can normalize JSON body with custom method returning error", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name":"John","age":30}`)
		c := NewContext[testStructNormalizableWithError](
			httptest.NewRequest("GET", "http://example.com/foo", reqBody),
			readOptions{})

		body, err := c.Body()
		require.Error(t, err)
		require.Equal(t, body.Name, "John")
		require.Equal(t, body.Age, 30)
	})

	t.Run("can read string body", func(t *testing.T) {
		// Create new Reader
		a := strings.NewReader("Hello World")

		// Test an http request
		r := httptest.NewRequest("GET", "http://example.com/foo", a)

		c := NewContext[string](r, readOptions{})

		_, err := c.Body()
		require.NoError(t, err)
	})
}

func FuzzContext_Body(f *testing.F) {
	f.Add("Hello Fuzz")

	f.Fuzz(func(t *testing.T, s string) {
		// Create new Reader
		a := strings.NewReader(s)

		// Test an http request
		r := httptest.NewRequest("GET", "http://example.com/foo", a)

		c := NewContext[string](r, readOptions{})

		_, err := c.Body()
		require.NoError(t, err)
	})
}

func BenchmarkContext_Body(b *testing.B) {
	b.Run("valid JSON body", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reqBody := strings.NewReader(`{"name":"John","age":30}`)
			c := NewContext[testStruct](
				httptest.NewRequest("GET", "http://example.com/foo", reqBody),
				readOptions{})
			_, err := c.Body()
			if err != nil {
				b.Fatal(err, "iteration", i)
			}
		}
	})

	// This test does make really sense because the body will not be accessed millions of times.
	// It however does show that caching the body works.
	// See [Body] for more information.
	b.Run("valid JSON body cache", func(b *testing.B) {
		reqBody := strings.NewReader(`{"name":"John","age":30}`)
		c := NewContext[testStruct](
			httptest.NewRequest("GET", "http://example.com/foo", reqBody),
			readOptions{})
		for i := 0; i < b.N; i++ {
			_, err := c.Body()
			if err != nil {
				b.Fatal(err, "iteration", i)
			}
		}
	})

	b.Run("invalid JSON body", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reqBody := strings.NewReader(`{"name":"John","age":30}`)
			c := NewContext[testStruct](
				httptest.NewRequest("GET", "http://example.com/foo", reqBody),
				readOptions{})
			_, err := c.Body()
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("string body", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reqBody := strings.NewReader(`{"name":"John","age":30}`)
			c := NewContext[testStruct](
				httptest.NewRequest("GET", "http://example.com/foo", reqBody),
				readOptions{})
			_, err := c.Body()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func TestContext_MustBody(t *testing.T) {
	t.Run("can read JSON body", func(t *testing.T) {
		// Create new Reader
		a := strings.NewReader(`{"name":"John","age":30}`)

		// Test an http request
		r := httptest.NewRequest("GET", "http://example.com/foo", a)

		c := NewContext[testStruct](r, readOptions{})

		body := c.MustBody()
		require.Equal(t, body.Name, "John")
		require.Equal(t, body.Age, 30)
	})

	t.Run("cannot read invalid JSON body", func(t *testing.T) {
		type testStruct struct {
			Name string `json:"name" validate:"required,min=3,max=10"`
			Age  int    `json:"age" validate:"min=18"`
		}

		reqBody := strings.NewReader(`{"name":"VeryLongName","age":12}`)
		c := NewContext[testStruct](
			httptest.NewRequest("GET", "http://example.com/foo", reqBody),
			readOptions{})

		require.Panics(t, func() {
			c.MustBody()
		})
	})
}
