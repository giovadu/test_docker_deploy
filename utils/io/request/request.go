package request

import (
	"encoding/json"
	"net/http"
)

func Json(r *http.Request, obj interface{}) error {
	return json.NewDecoder(r.Body).Decode(&obj)
}
