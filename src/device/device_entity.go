package device

import (
	"errors"
	"sync"
	"time"
)

var InvalidArgument = errors.New("invalid argument")

type Device struct {
	id                    string
	heartBeatsCount       int
	earliestHeartBeatTime time.Time
	lastHeartBeatTime     time.Time
	uploadTimes           []int64
	heartBeatLock         sync.RWMutex
	statsLock             sync.RWMutex
}

func NewDevice(id string) *Device {
	return &Device{
		id:          id,
		uploadTimes: []int64{},
	}
}

func (d *Device) AddHeartBeat(t time.Time) {
	d.heartBeatLock.Lock()
	defer d.heartBeatLock.Unlock()
	d.heartBeatsCount += 1
	if d.lastHeartBeatTime.Before(t) {
		d.lastHeartBeatTime = t
	}
	if d.earliestHeartBeatTime.IsZero() || d.earliestHeartBeatTime.After(t) {
		d.earliestHeartBeatTime = t
	}
}

func (d *Device) AddStats(sentAt time.Time, value int64) error {
	if value <= 0 {
		return InvalidArgument
	}
	d.statsLock.Lock()
	defer d.statsLock.Unlock()
	d.uploadTimes = append(d.uploadTimes, value)
	return nil
}

func (d *Device) GetUpTime() float64 {
	d.heartBeatLock.RLock()
	defer d.heartBeatLock.RUnlock()
	sub := d.lastHeartBeatTime.Sub(d.earliestHeartBeatTime).Minutes()
	if sub < 1 {
		return 0
	}

	return (float64(d.heartBeatsCount) / sub) * 100
}

func (d *Device) GetUploadTimeAVG() time.Duration {
	d.statsLock.RLock()
	defer d.statsLock.RUnlock()
	if len(d.uploadTimes) == 0 {
		return 0
	}

	var total int64
	for _, v := range d.uploadTimes {
		total += v
	}

	avgNanoseconds := total / int64(len(d.uploadTimes))
	return time.Duration(avgNanoseconds)
}
