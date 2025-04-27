package http

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, implement proper origin checks
	},
}

// HandleActivityStatsWS handles WebSocket connections for real-time activity stats
// @Summary Get real-time activity stats
// @Description Stream real-time activity statistics via WebSocket
// @Tags stats
// @Accept json
// @Produce json
// @Success 101 "Switching to WebSocket protocol"
// @Failure 400 {object} apperrors.Error
// @Router /stats/activity/ws [get]
func (h *UserHandler) HandleActivityStatsWS(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	// Get initial stats
	stats, statsIterator, err := h.useCase.GetActivityStats(c.Request().Context())
	if err != nil {
		return err
	}

	if err := ws.WriteJSON(stats); err != nil {
		return err
	}

	// Stream updates
	for stat := range statsIterator {
		if err := ws.WriteJSON(stat); err != nil {
			break
		}
	}

	return nil
}
