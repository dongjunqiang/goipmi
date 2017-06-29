package ipmi

import (
	"bytes"
	"fmt"
)

// RepositoryInfo get the Repository Info of the SDR
func (c *Client) RepositoryInfo() (*SDRRepositoryInfoResponse, error) {
	req := &Request{
		NetworkFunctionStorge,
		CommandGetSDRRepositoryInfo,
		&SDRRepositoryInfoRequest{},
	}
	res := &SDRRepositoryInfoResponse{}
	return res, c.Send(req, res)
}
func (c *Client) GetReserveSDRRepoForReserveId() (*ReserveRepositoryResponse, error) {
	req := &Request{
		NetworkFunctionStorge,
		CommandGetReserveSDRRepo,
		&ReserveSDRRepositoryRequest{},
	}
	res := &ReserveRepositoryResponse{}
	return res, c.send(req, res)
}

//Get SDR Command  33.12
func (c *Client) GetSDR(reservationID uint16, recordID uint16) (record SDRRecord, next uint16) {

	req := &Request{
		NetworkFunctionStorge,
		CommandGetSDR,
		&GetSDRCommandRequest{
			ReservationID:    reservationID,
			RecordID:         recordID,
			OffsetIntoRecord: 0,
			ByteToRead:       5,
		},
	}
	entire1 := new(bytes.Buffer)
	res := &GetSDRCommandResponse{}
	c.Send(req, res)
	readData1 := res.ReadData
	remain1 := readData1[4]
	entire1.Write(res.ReadData)

	req2 := &Request{
		NetworkFunctionStorge,
		CommandGetSDR,
		&GetSDRCommandRequest{
			ReservationID:    reservationID,
			RecordID:         recordID,
			OffsetIntoRecord: 5,
			ByteToRead:       uint8(remain1),
		},
	}
	res2 := &GetSDRCommandResponse{}
	c.Send(req2, res2)
	entire1.Write(res2.ReadData)

	//Unmarshalbinary and assert
	r2, _ := NewSDRFullSensor(0, "")
	r2.UnmarshalBinary(entire1.Bytes())

	return r2, res2.NextRecordID
}
func (c *Client) GetSensorList(reservationID uint16, recordID uint16) {
	var recordId uint16 = 0
	i := 0
	for recordId < 0xffff {
		sdrRecord, nId := c.GetSDR(reservationID, recordId)
		r2 := sdrRecord.(*SDRFullSensor)

		recordId = nId
		i = i+1
	}
	return 
}
