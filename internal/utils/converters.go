package utils

import (
	"fmt"
	"gqlgen-starter/internal/oops"
	"net/http"
	"strconv"
)

func ID64(in string) (int64, error) {
	out, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		return 0, &oops.CodedError{
			HumanMessage: fmt.Sprintf("%s not a valid number", in),
			Context:      fmt.Sprintf("%s not parsable as int64", in),
			HttpStatus:   http.StatusBadRequest,
			Err:          err,
		}
	}
	return out, nil
}
