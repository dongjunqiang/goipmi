/*
Copyright (c) 2014 EOITek, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ipmi

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSDRRecType_McDeviceLocator(t *testing.T) {

	deviceId := "AST2400"
	rec, err := NewSDRMcDeviceLocator(uint16(1), deviceId)
	assert.Nil(t, err)
	fmt.Println("DeviceId: ", rec.DeviceId())
	fmt.Println("RecordId: ", rec.RecordId())

	rec.sdrMcDeviceLocatorFields = sdrMcDeviceLocatorFields{
		DeviceSlaveAddr: 0x20,
		ChannelNumber:   0x00,
		PSNGI:           0x00,
		DeviceCap:       0xff,
		reserved:        [3]byte{0, 0, 0},
		EntityId:        0x00,
		EntityIns:       0x01,
		OEM:             0x00,
	}

	data, err := rec.MarshalBinary()
	assert.Nil(t, err)
	assert.NotNil(t, data)

	flen := 10 //sizeof sdrmcdevicelocatorfields
	assert.Equal(t, byte(0x01), data[0])
	assert.Equal(t, byte(0x00), data[1])
	assert.Equal(t, byte(SDR_RECORD_TYPE_MC_DEVICE_LOCATOR), data[3])
	assert.Equal(t, byte(flen+1+len(deviceId)), data[4])
	assert.Equal(t, byte(len(deviceId)), data[flen+1+4])
	assert.Equal(t, deviceId, string(data[flen+1+4+1:]))
}

func TestSDRRecType_FullSensor(t *testing.T) {
	deviceId := "Ambient Temp"
	rec, err := NewSDRFullSensor(uint16(321), deviceId)
	assert.Nil(t, err)
	fmt.Println("DeviceId: ", rec.DeviceId())
	fmt.Println("RecordId: ", rec.RecordId())

	rec.sdrFullSensorFields = sdrFullSensorFields{
		SensorOwnerId:        0x20,
		SensorOwnerLUN:       0x00,
		SensorNumber:         0x01,
		EntityId:             0x37,
		EntityIns:            0x01,
		SensorInit:           0x5f,
		SensorCap:            0x68,
		SensorType:           SDRSensorType(0x01),
		ReadingType:          SDRSensorReadingType(0x01),
		AssertionEventMask:   0x3285,
		DeassertionEventMask: 0x3285,
		DiscreteReadingMask:  0x1b1b,
		Unit:                 0x80,
		BaseUnit:             0x01,
		ModifierUnit:         0x00,
		Linearization:        0x00,
		MTol:                 0x0201,
		Bacc:                 0x0200,
		Acc:                  0x30,
		RBexp:                0x00,
		AnalogFlag:           0x07,
		NominalReading:       0x97,
		NormalMax:            0xc5,
		NormalMin:            0x8b,
		SensorMax:            0x7f,
		SensorMin:            0x80,
		U_NR:                 0x29,
		U_C:                  0x27,
		U_NC:                 0x25,
		L_NR:                 0x00,
		L_C:                  0x03,
		L_NC:                 0x08,
		PositiveHysteresis:   0x02,
		NegativeHysteresis:   0x02,
		reserved:             [2]byte{0, 0},
		OEM:                  0x00,
	}

	data, err := rec.MarshalBinary()
	assert.Nil(t, err)
	assert.NotNil(t, data)

	flen := 42
	assert.Equal(t, byte(0x41), data[0])
	assert.Equal(t, byte(0x01), data[1])
	assert.Equal(t, byte(SDR_RECORD_TYPE_FULL_SENSOR), data[3])
	assert.Equal(t, byte(flen+1+len(deviceId)), data[4])
	assert.Equal(t, byte(len(deviceId)), data[flen+1+4])
	assert.Equal(t, deviceId, string(data[flen+1+4+1:]))

	//test M,B
	assert.Equal(t, byte(0x01), data[24])
	assert.Equal(t, byte(0x02), data[25])
	assert.Equal(t, byte(0x00), data[26])
	assert.Equal(t, byte(0x02), data[27])
	assert.Equal(t, byte(0x30), data[28])
	assert.Equal(t, byte(0x00), data[29])

}
