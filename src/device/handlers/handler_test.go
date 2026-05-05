package handlers

import (
	"errors"
	"testing"
	"time"

	"interFleet/src/device"
	"interFleet/src/device/infrastructure/mocks"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type CommandHandlerSuite struct {
	suite.Suite
	ctrl    *gomock.Controller
	repo    *mocks.MockIDeviceReadRepository
	handler *CommandHandler
}

func TestCommandHandlerSuite(t *testing.T) {
	suite.Run(t, new(CommandHandlerSuite))
}

func (s *CommandHandlerSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.repo = mocks.NewMockIDeviceReadRepository(s.ctrl)
	s.handler = &CommandHandler{
		Repo: s.repo,
	}
}

func (s *CommandHandlerSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *CommandHandlerSuite) TestHandleNewHeartBeatCommand_RepoError() {
	expectedErr := errors.New("not found")
	cmd := HeartbeatCommand{
		DeviceID: "missing-device",
		SentAt:   "2026-05-04T10:00:00Z",
	}

	s.repo.EXPECT().
		GetById("missing-device").
		Return((*device.Device)(nil), expectedErr)

	err := s.handler.HandleNewHeartBeatCommand(cmd)

	s.Error(err)
	s.ErrorIs(err, expectedErr)
}

func (s *CommandHandlerSuite) TestHandleNewHeartBeatCommand_InvalidTimestamp() {
	dev := device.NewDevice("dev-1")
	cmd := HeartbeatCommand{
		DeviceID: "dev-1",
		SentAt:   "not-a-timestamp",
	}

	s.repo.EXPECT().
		GetById("dev-1").
		Return(dev, nil)

	err := s.handler.HandleNewHeartBeatCommand(cmd)

	s.Error(err)
}

func (s *CommandHandlerSuite) TestHandleNewHeartBeatCommand_Success_ChangesObservableState() {
	dev := device.NewDevice("dev-1")
	dev.AddHeartBeat(time.Date(2026, 5, 4, 10, 0, 0, 0, time.UTC))

	cmd := HeartbeatCommand{
		DeviceID: "dev-1",
		SentAt:   "2026-05-04T10:02:00Z",
	}

	s.repo.EXPECT().
		GetById("dev-1").
		Return(dev, nil)

	err := s.handler.HandleNewHeartBeatCommand(cmd)

	s.NoError(err)
	s.Equal(100.0, dev.GetUpTime())
}

func (s *CommandHandlerSuite) TestHandleNewHeartBeatCommand_Success_FirstHeartbeat() {
	dev := device.NewDevice("dev-1")

	cmd := HeartbeatCommand{
		DeviceID: "dev-1",
		SentAt:   "2026-05-04T10:02:00Z",
	}

	s.repo.EXPECT().
		GetById("dev-1").
		Return(dev, nil)

	err := s.handler.HandleNewHeartBeatCommand(cmd)

	s.NoError(err)
	s.Equal(0.0, dev.GetUpTime()) // first heartbeat only
}

func (s *CommandHandlerSuite) TestHandleNewHeartBeatCommand_Success_EarlierHeartbeat() {
	dev := device.NewDevice("dev-1")
	dev.AddHeartBeat(time.Date(2026, 5, 4, 10, 5, 0, 0, time.UTC))

	cmd := HeartbeatCommand{
		DeviceID: "dev-1",
		SentAt:   "2026-05-04T10:00:00Z",
	}

	s.repo.EXPECT().
		GetById("dev-1").
		Return(dev, nil)

	err := s.handler.HandleNewHeartBeatCommand(cmd)

	s.NoError(err)
	s.Equal(40.0, dev.GetUpTime()) // 2 heartbeats / 5 minutes * 100
}

func (s *CommandHandlerSuite) TestHandleNewHeartBeatCommand_NilDeviceFromRepo_Panics() {
	cmd := HeartbeatCommand{
		DeviceID: "dev-1",
		SentAt:   "2026-05-04T10:02:00Z",
	}

	s.repo.EXPECT().
		GetById("dev-1").
		Return(nil, nil)

	s.Panics(func() {
		_ = s.handler.HandleNewHeartBeatCommand(cmd)
	})
}
