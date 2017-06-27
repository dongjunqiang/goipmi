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
	"net"
	"testing"
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
	})
	rep.addRecord(&sDRRecordAndValue{
		SDRRecord: r2,
		value:     23.3,
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
	})
	rep.addRecord(&sDRRecordAndValue{
		SDRRecord: r2,
		value:     33.6,
	})

}

func TestSimulatorSDR_L_SetMBExp(t *testing.T) {
	r2, _ := NewSDRFullSensor(2, "System 3.3V")
	r2.SetMBExp(2, 0, 0, 0)
	assert.Equal(t, uint16(0x0002), r2.MTol)
	assert.Equal(t, uint16(0x0000), r2.Bacc)
	r2.SetMBExp(-2, 5, 0, -2)
	assert.Equal(t, uint16(0xc0fe), r2.MTol)
	assert.Equal(t, uint16(0x0005), r2.Bacc)
	assert.Equal(t, uint8(0xe0), r2.RBexp)
	r2.SetMBExp(-2, -5, -2, -2)
	assert.Equal(t, uint16(0xc0fe), r2.MTol)
	assert.Equal(t, uint16(0xc0fb), r2.Bacc)
	assert.Equal(t, uint8(0xee), r2.RBexp)
}

func TestSimulatorSDR_L_GetMBExp(t *testing.T) {
	r1, _ := NewSDRFullSensor(2, "System 3.3V")
	r1.SetMBExp(2, 0, 0, 0)
	m1, b1, be1, re1 := r1.GetMBExp()
	assert.Equal(t, int16(2), m1)
	assert.Equal(t, int16(0), b1)
	assert.Equal(t, int8(0), be1)
	assert.Equal(t, int8(0), re1)

	r1.SetMBExp(-2, 5, 0, -2)
	m2, b2, be2, re2 := r1.GetMBExp()
	assert.Equal(t, int16(-2), m2)
	assert.Equal(t, int16(5), b2)
	assert.Equal(t, int8(0), be2)
	assert.Equal(t, int8(-2), re2)

	r1.SetMBExp(-2, -5, -2, -2)
	m3, b3, be3, re3 := r1.GetMBExp()
	assert.Equal(t, int16(-2), m3)
	assert.Equal(t, int16(-5), b3)
	assert.Equal(t, int8(-2), be3)
	assert.Equal(t, int8(-2), re3)
}

func TestSimulatorSDR_L_CalValue1(t *testing.T) {
	r, _ := NewSDRFullSensor(2, "System 3.3V")
	//r.Unit = 0x80 //2's complement signed
	r.Unit = 0x00          //unsigned
	r.BaseUnit = 0x04      //Voltage
	r.Linearization = 0x00 //no linearization
	r.SetMBExp(2, 0, 0, -2)
	v := r.CalValue(3.36)
	assert.Equal(t, uint8(0xa8), v)
}

func TestSimulatorSDR_L_CalValue2(t *testing.T) {
	r, _ := NewSDRFullSensor(2, "Exhaust Temp")
	r.Unit = 0x80          //2's complement signed
	r.BaseUnit = 0x01      //Temprature
	r.Linearization = 0x00 //no linearization
	r.SetMBExp(1, 0, 0, 0)
	v := r.CalValue(23.0)
	assert.Equal(t, uint8(0x17), v)
}

func TestSimulatorSDR_L_CalValue3(t *testing.T) {
	r, _ := NewSDRFullSensor(2, "CPU1 DTS")
	r.Unit = 0x80          //2's complement signed
	r.BaseUnit = 0x01      //Temprature
	r.Linearization = 0x00 //no linearization
	r.SetMBExp(1, 0, 0, 0)
	v := r.CalValue(-49.0)
	assert.Equal(t, uint8(0xcf), v)
}

func TestSimulatorSDR_L_CalValue4(t *testing.T) {
	r, _ := NewSDRFullSensor(2, "Fan 1")
	r.Unit = 0x00          //2's complement signed
	r.BaseUnit = 0x12      //RPM
	r.Linearization = 0x00 //no linearization
	r.SetMBExp(63, 0, 0, 0)
	v := r.CalValue(2583.0)
	assert.Equal(t, uint8(0x29), v)
}
