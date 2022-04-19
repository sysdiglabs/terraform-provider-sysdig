package monitor

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/jmespath/go-jmespath"
	"github.com/spf13/cast"
)

func errorFromResponse(response *http.Response) error {
	var data interface{}
	err := json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return errors.New(response.Status)
	}

	search, err := jmespath.Search("[error, message, errors[].[reason, message]][][] | join(', ', @)", data)
	if err != nil {
		return errors.New(response.Status)
	}

	if searchArray, ok := search.([]interface{}); ok {
		return errors.New(strings.Join(cast.ToStringSlice(searchArray), ", "))
	}

	toString := cast.ToString(search)
	if toString == "" {
		return errors.New(response.Status)
	}
	return errors.New(toString)
}
