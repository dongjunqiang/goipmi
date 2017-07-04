package ipmi

import (
	"bytes"
	"errors"
	//"fmt"
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
func (c *Client) GetSDR(reservationID uint16, recordID uint16) (sdr *sDRRecordAndValue, next uint16) {
	req_step1 := &Request{
		NetworkFunctionStorge,
		CommandGetSDR,
		&GetSDRCommandRequest{
			ReservationID:    reservationID,
			RecordID:         recordID,
			OffsetIntoRecord: 0,
			ByteToRead:       5,
		},
	}
	recordKeyBody_Data := new(bytes.Buffer)
	res_step1 := &GetSDRCommandResponse{}
	c.Send(req_step1, res_step1)
	readData_step1 := res_step1.ReadData
	recordType := readData_step1[3]
	lenToRead_step2 := readData_step1[4]
	recordKeyBody_Data.Write(readData_step1)
	req_step2 := &Request{
		NetworkFunctionStorge,
		CommandGetSDR,
		&GetSDRCommandRequest{
			ReservationID:    reservationID,
			RecordID:         recordID,
			OffsetIntoRecord: 5,
			ByteToRead:       uint8(lenToRead_step2),
		},
	}
	res_step2 := &GetSDRCommandResponse{}
	c.Send(req_step2, res_step2)
	recordKeyBody_Data.Write(res_step2.ReadData)
	sdrRecordAndValue := c.CalSdrRecordValue(recordType, recordKeyBody_Data)
	return sdrRecordAndValue, res_step2.NextRecordID
}
func (c *Client) CalSdrRecordValue(recordType uint8, recordKeyBody_Data *bytes.Buffer) *sDRRecordAndValue {
	if recordType == SDR_RECORD_TYPE_FULL_SENSOR {
		//Unmarshalbinary and assert
		fullSensor, _ := NewSDRFullSensor(0, "")
		var sdrRecordAndValue = &sDRRecordAndValue{}
		fullSensor.UnmarshalBinary(recordKeyBody_Data.Bytes())
		sdrRecordAndValue.SDRRecord = fullSensor
		res, avai := c.CalFullSensorValue(fullSensor)
		sdrRecordAndValue.avail = avai
		sdrRecordAndValue.value = res
		return sdrRecordAndValue
	} else if recordType == SDR_RECORD_TYPE_COMPACT_SENSOR {
		//Unmarshalbinary and assert
		var sdrRecordAndValue = &sDRRecordAndValue{}
		compactSensor, _ := NewSDRCompactSensor(0, "")
		compactSensor.UnmarshalBinary(recordKeyBody_Data.Bytes())
		sdrRecordAndValue.SDRRecord = compactSensor
		//threshold type
		if compactSensor.ReadingType == SENSOR_READTYPE_THREADHOLD {
			// has analog value
			if compactSensor.Unit&0xc0 == 0xc0 {
				sensorReading, err := c.GetSensorReading(compactSensor.SensorNumber)
				if err != nil {
					sdrRecordAndValue.avail = false
					sdrRecordAndValue.value = 0.0
				} else {
					sdrRecordAndValue.avail = true
					sdrRecordAndValue.value = float64(sensorReading)
				}
			}
		}
		return sdrRecordAndValue
	}
	return nil
}

func (c *Client) CalFullSensorValue(fullSensor *SDRFullSensor) (float64, bool) {
	var result float64
	var avail bool
	//threshold type
	if fullSensor.ReadingType == SENSOR_READTYPE_THREADHOLD {
		// has analog value
		if fullSensor.Unit&0xc0 != 0xc0 {
			m, b, bexp, rexp := fullSensor.GetMBExp()
			sensorReading, err := c.GetSensorReading(fullSensor.SensorNumber)
			if err != nil {
				result = 0.0
				avail = false
			} else {
				switch (fullSensor.Unit & 0xc0) >> 6 {
				case 0:
					result = (float64(m)*float64(sensorReading) + float64(b)*math.Pow(10, float64(bexp))) * math.Pow(10, float64(rexp))
				case 1:
				case 2:
					result = (float64(int8(m)*int8(sensorReading)) + float64(b)*math.Pow(10, float64(rexp))) * math.Pow(10, float64(bexp))
				}
				avail = true
			}
		}
	}
	return result, avail
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
	if (res.ReadingAvail & 0x20) == 0 {
		readValue := res.SensorReading
		return readValue, nil
	}
	return uint8(0), errors.New("reading unAvailableß")

}
func (c *Client) GetSensorList(reservationID uint16, recordID uint16) []*sDRRecordAndValue {
	var recordId uint16 = 0
	var sdrRecAndVallist = make([]*sDRRecordAndValue, 20, 60)
	for recordId < 0xffff {
		sdrRecordAndValue, nId := c.GetSDR(reservationID, recordId)
		//r2 := sdrRecordAndValue.SDRRecord.(*SDRFullSensor)
		_ = append(sdrRecAndVallist, sdrRecordAndValue)
		recordId = nId
	}
	return sdrRecAndVallist
}
