package ipmi

import (
	"bytes"
	"errors"
	"fmt"
	"math"
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

	fmt.Println("GetSensorList r2.IDString=", r2.DeviceId())
	m, b, bexp, rexp := r2.GetMBExp()

	sensorReading, err := c.GetSensorReading(r2.SensorNumber)
	if err != nil {

	} else {
		var result float64
		switch (r2.Unit & 0xc0) >> 6 {
		case 0:
			result = (float64(m)*float64(sensorReading) + float64(b)*math.Pow(10, float64(bexp))) * math.Pow(10, float64(rexp))
		case 1:
		case 2:
			fmt.Println("sensorReading==", int8(sensorReading))
			result = (float64(int8(m)*int8(sensorReading)) + float64(b)*math.Pow(10, float64(rexp))) * math.Pow(10, float64(bexp))
		}
		fmt.Println("result=", result)
	}

	return r2, res2.NextRecordID
}

//Get Sensor Reading  35.14
func (c *Client) GetSensorReading(sensorNum uint8) (sensorReading uint8, err error) {
	req := &Request{
		NetworkFunctionSensorEvent,
		CommandGetSensorReading,
		&GetSensorReadingRequest{
			SensorNumber: sensorNum,
		},
	}
	res := &GetSensorReadingResponse{}
	c.Send(req, res)
	if res == nil {
		return uint8(0), errors.New("can not found sensor number")
	}
	readValue := res.SensorReading

	return readValue, nil
}
func (c *Client) GetSensorList(reservationID uint16, recordID uint16) {
	var recordId uint16 = 0
	for recordId < 0xffff {
		_, nId := c.GetSDR(reservationID, recordId)
		//r2 := sdrRecord.(*SDRFullSensor)
		//fmt.Println("GetSensorList r2.IDString=",r2.DeviceId())
		recordId = nId
	}
}
