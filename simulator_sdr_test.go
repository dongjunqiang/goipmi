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
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimulatorSDR_L_ReserveRepo(t *testing.T) {
	s := NewSimulator(net.UDPAddr{Port: 0})
	resp := s.reserveRepository(nil)

	rrr, ok := resp.(*ReserveRepositoryResponse)
	assert.True(t, ok)
	assert.Equal(t, CommandCompleted, rrr.CompletionCode)
	fmt.Println("ReservationId: ", rrr.ReservationId)
}

func TestSimulatorSDR_L_GetRecord(t *testing.T) {
	rep := NewRepo()
	r1, _ := NewSDRMcDeviceLocator(1, "aaa")
	r2, _ := NewSDRFullSensor(2, "bbb")
	rep.addRecord(&sDRRecordAndValue{
		SDRRecord: r1,
		value:     -49.0,
		avail:     true,
	})
	rep.addRecord(&sDRRecordAndValue{
		SDRRecord: r2,
		value:     23.3,
		avail:     true,
	})
	d0, next0 := rep.getRecordById(0)
	assert.NotNil(t, d0)
	assert.Equal(t, uint16(2), next0)

	d1, next1 := rep.getRecordById(1)
	assert.NotNil(t, d1)
	assert.Equal(t, uint16(2), next1)

	d2, next2 := rep.getRecordById(2)
	assert.NotNil(t, d2)
	assert.Equal(t, uint16(0xffff), next2)
}

func TestSimulatorSDR_L_GetSDR(t *testing.T) {
	s := NewSimulator(net.UDPAddr{Port: 0})
	resp := s.reserveRepository(nil)
	reserve, ok := resp.(*ReserveRepositoryResponse)
	assert.True(t, ok)
	assert.Equal(t, CommandCompleted, reserve.CompletionCode)

	entire := new(bytes.Buffer)
	//get first 5 bytes
	getsdr_req := &GetSDRCommandRequest{}
	getsdr_req.ReservationID = reserve.ReservationId
	getsdr_req.RecordID = 0
	getsdr_req.OffsetIntoRecord = 0
	getsdr_req.ByteToRead = 5
	m := &Message{}
	m.Data = messageDataToBytes(getsdr_req)

	rec_resp := s.getSDR(m)
	rec, ok := rec_resp.(*GetSDRCommandResponse)
	assert.True(t, ok)
	assert.Equal(t, CommandCompleted, rec.CompletionCode)
	assert.Equal(t, uint16(2), rec.NextRecordID)

	remain := int(rec.ReadData[4])
	fmt.Printf("remaining %d bytes\n", remain)
	//copy
	entire.Write(rec.ReadData)

	//get remain bytes
	getsdr_req.OffsetIntoRecord = 5
	getsdr_req.ByteToRead = uint8(remain)
	m.Data = messageDataToBytes(getsdr_req)
	rec_resp = s.getSDR(m)
	rec, ok = rec_resp.(*GetSDRCommandResponse)
	assert.True(t, ok)
	assert.Equal(t, CommandCompleted, rec.CompletionCode)
	assert.Equal(t, uint16(2), rec.NextRecordID)
	assert.Equal(t, remain, len(rec.ReadData))
	entire.Write(rec.ReadData)

	//Unmarshalbinary and assert
	r2, _ := NewSDRFullSensor(0, "")
	r2.UnmarshalBinary(entire.Bytes())
	assert.Equal(t, uint16(1), r2.RecordId())
	assert.Equal(t, SDRRecordType(SDR_RECORD_TYPE_FULL_SENSOR), r2.RecordType())
	assert.Equal(t, "Fan 1", r2.DeviceId())
	assert.Equal(t, uint8(0x00), r2.Unit)
	assert.Equal(t, uint8(0x12), r2.BaseUnit)
	M, _, _, _ := r2.GetMBExp()
	assert.Equal(t, int16(63), M)

}

//todo
func TestSimulatorSDR_L_GetSensorReading(t *testing.T) {
	rep := NewRepo()
	r1, _ := NewSDRMcDeviceLocator(1, "")
	r2, _ := NewSDRFullSensor(2, "System 3.3V")
	r2.SensorNumber = 0xb5
	r2.Unit = 0x00          //unsigned
	r2.BaseUnit = 0x04      //Voltage
	r2.Linearization = 0x00 //no linearization

	rep.addRecord(&sDRRecordAndValue{
		SDRRecord: r1,
		value:     -49.0,
		avail:     true,
	})
	rep.addRecord(&sDRRecordAndValue{
		SDRRecord: r2,
		value:     33.6,
		avail:     true,
	})

	getSensorNum_req := &GetSensorReadingRequest{}
	getSensorNum_req.SensorNumber = 5

	//	m := &Message{}
	//	m.Data = messageDataToBytes(getSensorNum_req)

	//	senN_resp := s.getSensorReading(m)
	//	rec, ok := senN_resp.(*GetSensorReadingResponse)
	//	assert.True(t, ok)
	//	assert.Equal(t, CommandCompleted, rec.CompletionCode)
	//	assert.Equal(t, uint16(2), rec.NextRecordID)
}
