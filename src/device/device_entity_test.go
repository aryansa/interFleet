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
	s.Equal(int64(0), d.uploadTimesCount)
	s.Equal(0.0, d.uploadTimeSum)
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

func (s *DeviceSuite) TestAddStats_ValidValue_UpdatesCountAndSum() {
	d := NewDevice("dev-1")
	sentAt := time.Date(2026, 5, 4, 10, 0, 0, 0, time.UTC)

	err := d.AddStats(sentAt, 100)

	s.NoError(err)
	s.Equal(int64(1), d.uploadTimesCount)
	s.Equal(100.0, d.uploadTimeSum)
}

func (s *DeviceSuite) TestAddStats_ZeroValue_ReturnsInvalidArgument() {
	d := NewDevice("dev-1")

	err := d.AddStats(time.Now(), 0)

	s.Error(err)
	s.ErrorIs(err, InvalidArgument)
	s.Equal(int64(0), d.uploadTimesCount)
	s.Equal(0.0, d.uploadTimeSum)
}

func (s *DeviceSuite) TestAddStats_NegativeValue_ReturnsInvalidArgument() {
	d := NewDevice("dev-1")

	err := d.AddStats(time.Now(), -1)

	s.Error(err)
	s.ErrorIs(err, InvalidArgument)
	s.Equal(int64(0), d.uploadTimesCount)
	s.Equal(0.0, d.uploadTimeSum)
}

func (s *DeviceSuite) TestAddStats_MultipleValues_AccumulatesCountAndSum() {
	d := NewDevice("dev-1")

	require.NoError(s.T(), d.AddStats(time.Now(), 100))
	require.NoError(s.T(), d.AddStats(time.Now(), 200))
	require.NoError(s.T(), d.AddStats(time.Now(), 300))

	s.Equal(int64(3), d.uploadTimesCount)
	s.Equal(600.0, d.uploadTimeSum)
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

func (s *DeviceSuite) TestGetUploadTimeAVG_SingleStat_ReturnsThatValue() {
	d := NewDevice("dev-1")

	require.NoError(s.T(), d.AddStats(time.Now(), 123456))

	s.Equal(time.Duration(123456), d.GetUploadTimeAVG())
}
