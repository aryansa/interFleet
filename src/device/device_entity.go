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

	//uploadTimes   []int64
	uploadTimesCount int64
	uploadTimeSum    float64
	heartBeatLock    sync.RWMutex
	uploadTimeLock   sync.RWMutex
}

func NewDevice(id string) *Device {
	return &Device{
		id: id,
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
	d.uploadTimeLock.Lock()
	defer d.uploadTimeLock.Unlock()
	d.uploadTimesCount += 1
	d.uploadTimeSum += float64(value)
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
	d.uploadTimeLock.RLock()
	defer d.uploadTimeLock.RUnlock()
	if d.uploadTimesCount == 0 {
		return 0
	}

	avgNanoseconds := d.uploadTimeSum / float64(d.uploadTimesCount)
	return time.Duration(avgNanoseconds)
}
