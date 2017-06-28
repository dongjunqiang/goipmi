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

const (
	CommandGetSDRRepositoryInfo = Command(0x20)
	CommandGetReserveSDRRepo    = Command(0x22)
	CommandGetSDR               = Command(0x23)
)

type ReserveSDRRepositoryRequest struct{}

// section 33.11
type ReserveRepositoryResponse struct {
	CompletionCode
	ReservationId uint16
}
type GetSDRCommandRequest struct {
	ReservationID    uint16
	RecordID         uint16
	OffsetIntoRecord uint8
	ByteToRead       uint8 //FFH means read entire record
}

type GetSDRCommandResponse struct {
	CompletionCode
	NextRecordID uint16
	ReadData     []byte
}
