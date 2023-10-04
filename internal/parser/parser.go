package parser

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type UserPass struct {
	User string
	Pass string
}

func GetQueryBool(c echo.Context, s string) (bool, error) {
	rest := false
	restRaw := c.QueryParam(s)

	if restRaw != "" {
		var err error

		rest, err = strconv.ParseBool(restRaw)
		if err != nil {
			return false, err
		}
	}

	return rest, nil
}

func GetAuthorizationBasic(r *http.Request) (*UserPass, error) {
	headerValue := r.Header.Get("authorization")
	if headerValue == "" {
		return nil, fmt.Errorf("authorization header not found")
	}

	components := strings.SplitN(headerValue, " ", 2)

	if len(components) != 2 || !strings.EqualFold(components[0], "Basic") {
		return nil, fmt.Errorf("schema Basic not found")
	}

	decodedValue, err := base64.StdEncoding.DecodeString(components[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse base64 basic credentials")
	}

	credential := strings.SplitN(string(decodedValue), ":", 2)
	if len(credential) != 2 {
		return nil, fmt.Errorf("failed to parse basic credentials")
	}

	return &UserPass{
		User: credential[0],
		Pass: credential[1],
	}, nil
}
