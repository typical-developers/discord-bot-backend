package httpx

import (
	"net/http"
)

func GetQueryParam(r *http.Request, key string, defaultValue ...string) string {
	param := r.URL.Query().Get(key)
	if param == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return ""
	}

	return param
}
