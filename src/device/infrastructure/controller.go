package infrastructure

import (
	"errors"
	"interFleet/src/device"
	"interFleet/src/device/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Handler handlers.CommandHandler
	Repo    device.IDeviceReadRepository
}

func (h *Controller) RegisterRoutes(registry func(method string, url string, handler gin.HandlerFunc)) {
	registry(http.MethodPost, "/devices/:device_id/heartbeat", h.addHeartBeat)
	registry(http.MethodPost, "/devices/:device_id/stats", h.addStats)
	registry(http.MethodGet, "/devices/:device_id/stats", h.getStats)
}

func (h *Controller) errorResponse(c *gin.Context, err error) {
	statusCode := http.StatusInternalServerError
	if errors.Is(err, NotFoundErr) {
		statusCode = http.StatusNotFound
	}
	c.JSON(statusCode, map[string]string{"msg": err.Error()})
}

func (h *Controller) addHeartBeat(context *gin.Context) {
	command := handlers.HeartbeatCommand{DeviceID: context.Param("device_id")}
	if err := context.ShouldBindJSON(&command); err != nil {
		h.errorResponse(context, err)
		return
	}
	if err := h.Handler.HandleNewHeartBeatCommand(command); err != nil {
		h.errorResponse(context, err)
		return
	}
	context.Status(http.StatusNoContent)
}

func (h *Controller) addStats(context *gin.Context) {
	command := handlers.UploadStatsCommand{DeviceID: context.Param("device_id")}
	if err := context.ShouldBindJSON(&command); err != nil {
		h.errorResponse(context, err)
		return
	}
	if err := h.Handler.HandleUploadStatCommand(command); err != nil {
		h.errorResponse(context, err)
		return
	}
	context.Status(http.StatusNoContent)
}

func (h *Controller) getStats(c *gin.Context) {
	deviceID := c.Param("device_id")

	dev, err := h.Repo.GetById(deviceID)
	if err != nil {
		h.errorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"avg_upload_time": dev.GetUploadTimeAVG().String(),
		"uptime":          dev.GetUpTime(),
	})
}
