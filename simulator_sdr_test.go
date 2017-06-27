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
	rep.addRecord(r1)
	rep.addRecord(r2)
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
