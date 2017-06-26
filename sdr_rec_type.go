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
	"bytes"
	"encoding/binary"
	"errors"
)

var (
	ErrDeviceIdMustLess16 = errors.New("Device Id must be less or equal to 16 bytes length")
)

type SDRRecord interface {
	DeviceId() string
	RecordId() uint16
	RecordType() SDRRecordType
}

type SDRRecordHeader struct {
	recordId   uint16
	SDRVersion uint8
	rtype      SDRRecordType
}

// section 43.9
type sdrMcDeviceLocatorFields struct { //size 10
	DeviceSlaveAddr uint8
	ChannelNumber   uint8

	PSNGI     uint8
	DeviceCap uint8
	reserved  [3]byte
	EntityId  uint8
	EntityIns uint8
	OEM       uint8
}

type SDRMcDeviceLocator struct {
	SDRRecordHeader
	sdrMcDeviceLocatorFields
	deviceId string
}

func NewSDRMcDeviceLocator(id uint16, name string) (*SDRMcDeviceLocator, error) {
	if len(name) > 16 {
		return nil, ErrDeviceIdMustLess16
	}
	r := &SDRMcDeviceLocator{}
	r.recordId = id
	r.rtype = SDR_RECORD_TYPE_MC_DEVICE_LOCATOR
	r.SDRVersion = 0x51
	r.deviceId = name
	return r, nil
}

func (r *SDRMcDeviceLocator) DeviceId() string {
	return r.deviceId
}

func (r *SDRMcDeviceLocator) RecordId() uint16 {
	return r.recordId
}

func (r *SDRMcDeviceLocator) RecordType() SDRRecordType {
	return r.rtype
}

func (r *SDRMcDeviceLocator) MarshalBinary() (data []byte, err error) {
	hb := new(bytes.Buffer)
	fb := new(bytes.Buffer)
	db := new(bytes.Buffer)
	binary.Write(hb, binary.LittleEndian, r.SDRRecordHeader)
	binary.Write(fb, binary.LittleEndian, r.sdrMcDeviceLocatorFields)
	db.WriteByte(byte(len(r.DeviceId())))
	db.WriteString(r.DeviceId())

	//merge all
	recLen := uint8(fb.Len() + db.Len())
	hb.WriteByte(byte(recLen))
	hb.Write(fb.Bytes())
	hb.Write(db.Bytes())
	return hb.Bytes(), nil
}

// section 43.1
type sdrFullSensorFields struct { //size 42
	SensorOwnerId        uint8
	SensorOwnerLUN       uint8
	SensorNumber         uint8
	EntityId             uint8
	EntityIns            uint8
	SensorInit           uint8
	SensorCap            uint8
	SensorType           SDRSensorType
	ReadingType          SDRSensorReadingType
	AssertionEventMask   uint16
	DeassertionEventMask uint16
	DiscreteReadingMask  uint16
	Unit                 uint8
	BaseUnit             uint8
	ModifierUnit         uint8
	Linearization        uint8
	MTol                 uint16
	Bacc                 uint32
	AnalogFlag           uint8
	NominalReading       uint8
	NormalMax            uint8
	NormalMin            uint8
	SensorMax            uint8
	SensorMin            uint8
	U_NR                 uint8
	U_C                  uint8
	U_NC                 uint8
	L_NR                 uint8
	L_C                  uint8
	L_NC                 uint8
	PositiveHysteresis   uint8
	NegativeHysteresis   uint8
	reserved             [2]byte
	OEM                  uint8
}

type SDRFullSensor struct {
	SDRRecordHeader
	sdrFullSensorFields
	deviceId string
}

func NewSDRFullSensor(id uint16, name string) (*SDRFullSensor, error) {
	if len(name) > 16 {
		return nil, ErrDeviceIdMustLess16
	}
	r := &SDRFullSensor{}
	r.recordId = id
	r.rtype = SDR_RECORD_TYPE_FULL_SENSOR
	r.SDRVersion = 0x51
	r.deviceId = name
	return r, nil
}

func (r *SDRFullSensor) DeviceId() string {
	return r.deviceId
}

func (r *SDRFullSensor) RecordId() uint16 {
	return r.recordId
}

func (r *SDRFullSensor) RecordType() SDRRecordType {
	return r.rtype
}

func (r *SDRFullSensor) MarshalBinary() (data []byte, err error) {
	hb := new(bytes.Buffer)
	fb := new(bytes.Buffer)
	db := new(bytes.Buffer)
	binary.Write(hb, binary.LittleEndian, r.SDRRecordHeader)
	binary.Write(fb, binary.LittleEndian, r.sdrFullSensorFields)
	db.WriteByte(byte(len(r.DeviceId())))
	db.WriteString(r.DeviceId())

	//merge all
	recLen := uint8(fb.Len() + db.Len())
	hb.WriteByte(byte(recLen))
	hb.Write(fb.Bytes())
	hb.Write(db.Bytes())
	return hb.Bytes(), nil
}

type SDRCompactSensor struct {
}
