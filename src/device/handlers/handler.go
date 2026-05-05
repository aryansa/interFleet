package handlers

import (
	"interFleet/src/device"
	"time"
)

type CommandHandler struct {
	Repo device.IDeviceReadRepository
}

type HeartbeatCommand struct {
	SentAt   string `json:"sent_at" binding:"required"`
	DeviceID string `json:"device_id" binding:"required"`
}

type UploadStatsCommand struct {
	DeviceID   string `json:"device_id" binding:"required"`
	SentAt     string `json:"sent_at" binding:"required"`
	UploadTime int64  `json:"upload_time" binding:"required"`
}

func (h *CommandHandler) HandleNewHeartBeatCommand(command HeartbeatCommand) error {
	deviceEntity, err := h.Repo.GetById(command.DeviceID)
	if err != nil {
		return err
	}
	sentAt, err := time.Parse(time.RFC3339, command.SentAt)
	if err != nil {
		return err
	}
	deviceEntity.AddHeartBeat(sentAt)
	return nil
}

func (h *CommandHandler) HandleUploadStatCommand(command UploadStatsCommand) error {
	deviceEntity, err := h.Repo.GetById(command.DeviceID)
	if err != nil {
		return err
	}
	sentAt, err := time.Parse(time.RFC3339, command.SentAt)
	if err != nil {
		return err
	}

	return deviceEntity.AddStats(sentAt, command.UploadTime)
}
