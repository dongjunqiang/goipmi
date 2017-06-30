/*
Copyright (c) 2014 VMware, Inc. All Rights Reserved.

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

//func TestClient(t *testing.T) {
//	s := NewSimulator(net.UDPAddr{Port: 0})
//	err := s.Run()
//	assert.NoError(t, err)

//	client, err := NewClient(s.NewConnection())
//	assert.NoError(t, err)

//	err = client.Open()
//	assert.NoError(t, err)

//	err = client.SetBootDevice(BootDevicePxe)
//	assert.NoError(t, err)

//	err = client.Close()
//	assert.NoError(t, err)
//	s.Stop()
//}

//func TestDeviceID(t *testing.T) {
//	s := NewSimulator(net.UDPAddr{})
//	err := s.Run()
//	assert.NoError(t, err)

//	client, err := NewClient(s.NewConnection())
//	assert.NoError(t, err)

//	err = client.Open()
//	assert.NoError(t, err)

//	tests := []OemID{
//		OemDell, OemHP,
//	}

//	for _, test := range tests {
//		s.SetHandler(NetworkFunctionApp, CommandGetDeviceID, func(*Message) Response {
//			return &DeviceIDResponse{
//				CompletionCode: CommandCompleted,
//				ManufacturerID: test,
//			}
//		})

//		id, err := client.DeviceID()
//		assert.NoError(t, err)
//		assert.Equal(t, test, id.ManufacturerID)
//		assert.Equal(t, test.String(), id.ManufacturerID.String())
//	}

//	err = client.Close()
//	assert.NoError(t, err)
//	s.Stop()
//}

//

func TestGetReserveSDRRepoForReserveId(t *testing.T) {
	s := NewSimulator(net.UDPAddr{})
	err := s.Run()
	assert.NoError(t, err)

	client, err := NewClient(s.NewConnection())
	assert.NoError(t, err)

	err = client.Open()
	assert.NoError(t, err)

	s.SetHandler(NetworkFunctionStorge, CommandGetReserveSDRRepo, func(*Message) Response {
		return &ReserveRepositoryResponse{
			CompletionCode: CommandCompleted,
			ReservationId:  28661,
		}
	})
	fmt.Println("TestGetReserveSDRRepoForReserveId")
	resp, err := client.GetReserveSDRRepoForReserveId()
	assert.NoError(t, err)
	assert.Equal(t, CommandCompleted, resp.CompletionCode)
	assert.Equal(t, uint16(28661), resp.ReservationId)
	//assert.Equal(t, uint(8), resp.OperationSupprot)
	//assert.Equal(t, test.String(), id.ManufacturerID.String())

	err = client.Close()
	assert.NoError(t, err)
	s.Stop()
}
