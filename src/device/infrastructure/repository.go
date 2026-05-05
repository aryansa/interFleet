package infrastructure

import (
	"encoding/csv"
	"errors"
	"fmt"
	"interFleet/src/device"
	"os"
)

type DeviceRepository struct {
	deviceCache map[string]*device.Device
}

var NotFoundErr = errors.New("device not found")

func NewDeviceRepository(csvAddress string) (*DeviceRepository, error) {
	file, err := os.Open(csvAddress)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 2. Create a new reader
	reader := csv.NewReader(file)

	// 3. Read all records at once
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	deviceRepository := &DeviceRepository{
		deviceCache: make(map[string]*device.Device),
	}
	records = records[1:]
	// 4. Iterate through the records
	for _, record := range records {
		if len(record) < 1 {
			fmt.Println("Skipping empty CSV file")
			continue
		}
		deviceRepository.deviceCache[record[0]] = device.NewDevice(record[0])
	}
	fmt.Printf("Successfully found %d devices\n", len(deviceRepository.deviceCache))
	return deviceRepository, nil
}

func (d *DeviceRepository) GetById(id string) (*device.Device, error) {
	tempDevice, ok := d.deviceCache[id]
	if !ok {
		return nil, NotFoundErr
	}
	return tempDevice, nil
}
