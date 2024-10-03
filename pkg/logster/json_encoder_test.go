package logster

import (
	"bytes"
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestJsonEncoder_unJson(t *testing.T) {
	enc := newJSONEncoder(zapcore.EncoderConfig{}, false)
	enc.AddString("ololo", "{\"123\":123}")
	require.Equal(t, "\"ololo\":\"(\\\"123\\\":123)\"", enc.buf.String())

	enc = newJSONEncoder(zapcore.EncoderConfig{}, false)
	enc.AddByteString("ololo", []byte("{hello:123}")) // fmt struct looks like json for this crutch
	require.Equal(t, "\"ololo\":\"(hello:123)\"", enc.buf.String())
}

func TestJsonEncoder_TestReflected(t *testing.T) {
	enc := newJSONEncoder(zapcore.EncoderConfig{}, false)
	var in []interface{}
	b, err := StringArray{"A", "B", "C"}.Value()
	require.NoError(t, err)
	in = append(in, b)
	err = enc.AddReflected("test", in)
	require.NoError(t, err)
	require.Equal(t, `"test":"[{\"A\",\"B\",\"C\"}]"`, enc.buf.String())
}

// StringArray represents a one-dimensional array of the PostgreSQL character types.
type StringArray []string

// Value implements the driver.Valuer interface.
func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, 2*N bytes of quotes,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+3*n)
		b[0] = '{'

		b = appendArrayQuotedBytes(b, []byte(a[0]))
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = appendArrayQuotedBytes(b, []byte(a[i]))
		}

		return string(append(b, '}')), nil
	}

	return "{}", nil
}

func appendArrayQuotedBytes(b, v []byte) []byte {
	b = append(b, '"')
	for {
		i := bytes.IndexAny(v, `"\`)
		if i < 0 {
			b = append(b, v...)
			break
		}
		if i > 0 {
			b = append(b, v[:i]...)
		}
		b = append(b, '\\', v[i])
		v = v[i+1:]
	}
	return append(b, '"')
}
