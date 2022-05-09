package pkg

import (
	"errors"
	"io"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalAddress(a string) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(a))
	})
}

func UnmarshalAddress(v interface{}) (string, error) {
	value, ok := v.(string)
	if !ok {
		return "", errors.New("address must be a string")
	}
	return value, nil
}
