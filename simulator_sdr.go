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
	//"fmt"
	"math"
	"math/rand"
	"time"
)

const sdrOpSupport = SDR_OP_SUP_RESERVE_REPO

type repo struct {
	sdrRepo []SDRRecord
}

func NewRepo() *repo {
	rep := &repo{}
	rep.sdrRepo = make([]SDRRecord, 0)
	return rep
}

func (rep *repo) addRecord(rec SDRRecord) {
	rep.sdrRepo = append(rep.sdrRepo, rec)
}

// return the marshaled bytes and next record id
func (rep *repo) getRecordById(id uint16) ([]byte, uint16) {

	if len(rep.sdrRepo) < 1 {
		return nil, 0xffff
	}

	//if not match, try to find the nearest one
	var distance = uint16(0xffff)
	var index = 0
	var record SDRRecord
	var next uint16
	var data []byte
	for i, rec := range rep.sdrRepo {
		_distance := uint16(math.Abs(float64(rec.RecordId()) - float64(id)))
		if _distance < distance {
			distance = _distance
			index = i
		}
	}

	if index >= len(rep.sdrRepo)-1 {
		next = 0xffff
	} else {
		next = rep.sdrRepo[index+1].RecordId()
	}

	record = rep.sdrRepo[index]

	switch record.RecordType() {
	case SDR_RECORD_TYPE_FULL_SENSOR:
		r, _ := record.(*SDRFullSensor)
		data, _ = r.MarshalBinary()
		break
	case SDR_RECORD_TYPE_MC_DEVICE_LOCATOR:
		r, _ := record.(*SDRMcDeviceLocator)
		data, _ = r.MarshalBinary()
		break
	default:
		panic("Unsupport record type")
	}

	return data, next
}

func (s *Simulator) repositoryInfo(*Message) Response {
	return &SDRRepositoryInfoResponse{
		CompletionCode:              CommandCompleted,
		SDRVersion:                  0x51,
		RecordCount:                 10,
		FreeSpaceInBytes:            0,
		TimestampMostRecentAddition: 0,
		TimestampMostRecentErase:    0,
		OperationSupprot:            sdrOpSupport,
	}
}

func (s *Simulator) reserveRepository(*Message) Response {

	var rId uint16
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rId = uint16(r.Intn(65536))

	return &ReserveRepositoryResponse{
		CompletionCode: CommandCompleted,
		ReservationId:  rId,
	}
}

func (s *Simulator) getNextSDR(*Message) Response {
	return nil
}
