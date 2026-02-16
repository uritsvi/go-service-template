package health

import "go-service-template/internal/models"

// OKResponse returns a health OK payload.
func OKResponse() models.Response {
	return models.Response{Status: "ok"}
}
