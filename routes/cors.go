package routes

import (
	"net/http"

	"github.com/gorilla/handlers"
)

func LoadCors(r http.Handler) http.Handler {
	headers := handlers.AllowedHeaders([]string{"X-Request", "Content-Type", "AuthorizationX-Request", "Content-Type", "Authorization"}) // X-Request", "Content-Type", "AuthorizationX-Request", "Content-Type", "Authorization
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"}) // should change in the future when client is ready
	return handlers.CORS(headers, methods, origins)(r)
}
