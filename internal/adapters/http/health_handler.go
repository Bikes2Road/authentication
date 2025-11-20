package http

import (
	"net/http"

	"github.com/bikes2road/authentication/internal/ports"
	"github.com/gin-gonic/gin"
)

type healthHandler struct{}

// NewHealthHandler crea una nueva instancia del handler de health check
func NewHealthHandler() ports.HealthHandler {
	return &healthHandler{}
}

// Health godoc
// @Summary      Health check
// @Description  Verifica que el servicio esté funcionando
// @Tags         health
// @Produce      json
// @Success      200 {object} HealthResponse "Servicio funcionando correctamente"
// @Router       /health [get]
func (h *healthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status: "OK",
	})
}

// HealthResponse representa la respuesta del health check
type HealthResponse struct {
	Status string `json:"status"`
}
