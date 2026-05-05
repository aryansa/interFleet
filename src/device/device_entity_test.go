package device

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DeviceSuite struct {
	suite.Suite
}

func TestDeviceSuite(t *testing.T) {
	suite.Run(t, new(DeviceSuite))
}

func (s *DeviceSuite) TestNewDevice_InitialState() {
	d := NewDevice("dev-1")

	s.Require().NotNil(d)
	s.Equal("dev-1", d.id)
	s.Equal(0, d.heartBeatsCount)
	s.True(d.earliestHeartBeatTime.IsZero())
	s.True(d.lastHeartBeatTime.IsZero())
	s.Empty(d.uploadTimes)
}

func (s *DeviceSuite) TestAddHeartBeat_FirstHeartbeatSetsEarliestAndLast() {
	d := NewDevice("dev-1")
	t1 := time.Date(2026, 5, 4, 10, 0, 0, 0, time.UTC)

	d.AddHeartBeat(t1)

	s.Equal(1, d.heartBeatsCount)
	s.True(d.earliestHeartBeatTime.Equal(t1))
	s.True(d.lastHeartBeatTime.Equal(t1))
}

func (s *DeviceSuite) TestAddHeartBeat_OutOfOrderUpdatesEarliestAndLastCorrectly() {
	d := NewDevice("dev-1")

	t2 := time.Date(2026, 5, 4, 10, 2, 0, 0, time.UTC)
	t1 := time.Date(2026, 5, 4, 10, 0, 0, 0, time.UTC)
	t3 := time.Date(2026, 5, 4, 10, 5, 0, 0, time.UTC)

	d.AddHeartBeat(t2)
	d.AddHeartBeat(t1)
	d.AddHeartBeat(t3)

	s.Equal(3, d.heartBeatsCount)
	s.True(d.earliestHeartBeatTime.Equal(t1))
	s.True(d.lastHeartBeatTime.Equal(t3))
}

func (s *DeviceSuite) TestGetUpTime_NoHeartbeats_ReturnsZero() {
	d := NewDevice("dev-1")
	s.Equal(0.0, d.GetUpTime())
}

func (s *DeviceSuite) TestGetUpTime_OneHeartbeat_ReturnsZero() {
	d := NewDevice("dev-1")
	d.AddHeartBeat(time.Date(2026, 5, 4, 10, 0, 0, 0, time.UTC))

	s.Equal(0.0, d.GetUpTime())
}

func (s *DeviceSuite) TestGetUpTime_MultipleHeartbeatsAcrossMinutes() {
	d := NewDevice("dev-1")

	d.AddHeartBeat(time.Date(2026, 5, 4, 10, 0, 0, 0, time.UTC))
	d.AddHeartBeat(time.Date(2026, 5, 4, 10, 1, 0, 0, time.UTC))
	d.AddHeartBeat(time.Date(2026, 5, 4, 10, 2, 0, 0, time.UTC))

	s.Equal(150.0, d.GetUpTime())
}

func (s *DeviceSuite) TestAddStats_ValidValue_AppendsUploadTime() {
	d := NewDevice("dev-1")
	sentAt := time.Date(2026, 5, 4, 10, 0, 0, 0, time.UTC)

	err := d.AddStats(sentAt, 100)

	s.NoError(err)
	s.Len(d.uploadTimes, 1)
	s.Equal(int64(100), d.uploadTimes[0])
}

func (s *DeviceSuite) TestAddStats_ZeroValue_ReturnsInvalidArgument() {
	d := NewDevice("dev-1")

	err := d.AddStats(time.Now(), 0)

	s.Error(err)
	s.ErrorIs(err, InvalidArgument)
	s.Empty(d.uploadTimes)
}

func (s *DeviceSuite) TestAddStats_NegativeValue_ReturnsInvalidArgument() {
	d := NewDevice("dev-1")

	err := d.AddStats(time.Now(), -1)

	s.Error(err)
	s.ErrorIs(err, InvalidArgument)
	s.Empty(d.uploadTimes)
}

func (s *DeviceSuite) TestGetUploadTimeAVG_NoStats_ReturnsZero() {
	d := NewDevice("dev-1")
	s.Equal(time.Duration(0), d.GetUploadTimeAVG())
}

func (s *DeviceSuite) TestGetUploadTimeAVG_MultipleStats_ReturnsAverage() {
	d := NewDevice("dev-1")

	require.NoError(s.T(), d.AddStats(time.Now(), 100))
	require.NoError(s.T(), d.AddStats(time.Now(), 200))
	require.NoError(s.T(), d.AddStats(time.Now(), 300))

	s.Equal(time.Duration(200), d.GetUploadTimeAVG())
}
