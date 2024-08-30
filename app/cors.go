package app

import (
	"net/http"

	"github.com/rs/cors"
)

var Cors = cors.New(cors.Options{
	AllowedOrigins: []string{"*"},
	AllowedMethods: []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodOptions,
	},
	AllowedHeaders: []string{
		"Accept",
		"content-type",
		"Content-Length",
		"Accept-Encoding",
		"X-CSRF-Token",
		"Authorization",
	},
	AllowCredentials: true,
})
