package scalar

import (
	"errors"
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/cosmos/cosmos-sdk/types"
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
	if _, err := types.AccAddressFromBech32(value); err != nil {
		return "", err
	}
	return value, nil
}

func MarshalUInt64(i uint64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, fmt.Sprintf("%d", i))
	})
}

func UnmarshalUInt64(v interface{}) (uint64, error) {
	switch v := v.(type) {
	case int:
		return uint64(v), nil
	default:
		return 0, fmt.Errorf("%T is not an uint64", v)
	}
}
