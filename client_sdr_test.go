package ipmi

import (
	"net"
	//"reflect"
	//"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepositoryInfo(t *testing.T) {
	s := NewSimulator(net.UDPAddr{})
	err := s.Run()
	assert.NoError(t, err)

	client, err := NewClient(s.NewConnection())
	assert.NoError(t, err)

	err = client.Open()
	assert.NoError(t, err)

	s.SetHandler(NetworkFunctionStorge, CommandGetSDRRepositoryInfo, func(*Message) Response {
		return &SDRRepositoryInfoResponse{
			CompletionCode:   CommandCompleted,
			SDRVersion:       0x51,
			RecordCount:      55,
			OperationSupprot: SDR_OP_SUP_ALLOC_INFO,
		}
	})

	resp, err := client.RepositoryInfo()
	assert.NoError(t, err)
	assert.Equal(t, CommandCompleted, resp.CompletionCode)
	assert.Equal(t, uint16(55), resp.RecordCount)
	assert.Equal(t, uint8(SDR_OP_SUP_ALLOC_INFO), resp.OperationSupprot)

	err = client.Close()
	assert.NoError(t, err)
	s.Stop()

}
func TestGetReserveSDRRepoForReserveId1(t *testing.T) {
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
			ReservationId:  25555,
		}
	})

	resp, err := client.GetReserveSDRRepoForReserveId()
	assert.NoError(t, err)
	assert.Equal(t, CommandCompleted, resp.CompletionCode)
	assert.Equal(t, uint16(25555), resp.ReservationId)
	//assert.Equal(t, test.String(), id.ManufacturerID.String())

	err = client.Close()
	assert.NoError(t, err)
	s.Stop()
}

func TestGetSDR(t *testing.T) {
	s := NewSimulator(net.UDPAddr{})
	resp := s.reserveRepository(nil)
	reserve, _ := resp.(*ReserveRepositoryResponse)

	err := s.Run()
	assert.NoError(t, err)

	client, err := NewClient(s.NewConnection())
	assert.NoError(t, err)

	err = client.Open()
	assert.NoError(t, err)

	r := &SDRFullSensor{}
	r.Recordid = 5
	r.Rtype = SDR_RECORD_TYPE_FULL_SENSOR
	r.SDRVersion = 0x51
	r.deviceId = "test deviceId"
	r.Unit = 0x00
	r.SensorNumber = 0x04
	r.BaseUnit = 0x12
	r.SetMBExp(63, 0, 0, 0)
	data1, _ := r.MarshalBinary()

	response := &GetSDRCommandResponse{}
	response.CompletionCode = CommandCompleted
	response.NextRecordID = 65535

	s.SetHandler(NetworkFunctionStorge, CommandGetSDR, func(m *Message) Response {
		request := &GetSDRCommandRequest{}
		if err := m.Request(request); err != nil {
			return err
		}
		response.ReadData = data1[request.OffsetIntoRecord : request.OffsetIntoRecord+request.ByteToRead]
		return response
	})

	sdrRecordAndValue, nextRecordId := client.GetSDR(reserve.ReservationId, 0)
	r2 := sdrRecordAndValue.SDRRecord.(*SDRFullSensor)
	assert.Equal(t, SDRRecordType(SDR_RECORD_TYPE_FULL_SENSOR), r2.Rtype)
	assert.Equal(t, "test deviceId", r2.DeviceId())
	assert.Equal(t, uint16(65535), nextRecordId)

	err = client.Close()
	assert.NoError(t, err)
	s.Stop()

}

func TestCalFullSensorValue(t *testing.T) {
	fs1 := &SDRFullSensor{}
	fs1.SetMBExp(8, 0, 0, 0)
	fs1.ReadingType = SENSOR_READTYPE_THREADHOLD
	fs1.Unit = 0x0
	res, avail := CalFullSensorValue(fs1, 0x11)
	assert.Equal(t, float64(136.0), res)
	assert.Equal(t, true, avail)

	fs1.SetMBExp(1, 0, 0, 0)
	fs1.ReadingType = SENSOR_READTYPE_THREADHOLD
	fs1.Unit = 0x80
	res, avail = CalFullSensorValue(fs1, 0xcf)
	assert.Equal(t, float64(-49.0), res)
	assert.Equal(t, true, avail)

	fs1.SetMBExp(2, 0, 0, -2)
	fs1.ReadingType = SENSOR_READTYPE_THREADHOLD
	fs1.Unit = 0x00
	res, avail = CalFullSensorValue(fs1, 0xa8)
	assert.Equal(t, float64(3.36), res)

	fs1.SetMBExp(2, 0, 0, -2)
	fs1.ReadingType = SENSOR_READTYPE_THREADHOLD
	fs1.Unit = 0x00
	res, avail = CalFullSensorValue(fs1, 0xa8)
	assert.Equal(t, float64(3.36), res)

	assert.Equal(t, true, avail)
}
func TestCalCompactSensorValue(t *testing.T) {
	cs1 := &SDRCompactSensor{}
	cs1.ReadingType = SENSOR_READTYPE_SENSORSPECIF
	cs1.Unit = 0xc0
	res, avail := CalCompactSensorValue(cs1, 0x11)
	assert.Equal(t, float64(17), res)
	assert.Equal(t, true, avail)
}
func TestGetSensorReading(t *testing.T) {

	s := NewSimulator(net.UDPAddr{})

	err := s.Run()
	assert.NoError(t, err)

	client, err := NewClient(s.NewConnection())
	assert.NoError(t, err)

	s.SetHandler(NetworkFunctionSensorEvent, CommandGetSensorReading, func(m *Message) Response {
		return &GetSensorReadingResponse{
			CompletionCode: CommandCompleted,
			SensorReading:  56,
		}
	})

	err = client.Open()
	assert.NoError(t, err)

	if SensorReading, err := client.GetSensorReading(0x04); err == nil {
		assert.Equal(t, uint8(56), SensorReading)
	}
}
func TestGetSensorList(t *testing.T) {

	s := NewSimulator(net.UDPAddr{})
	resp := s.reserveRepository(nil)
	reserve, _ := resp.(*ReserveRepositoryResponse)

	err := s.Run()
	assert.NoError(t, err)

	client, err := NewClient(s.NewConnection())
	assert.NoError(t, err)

	err = client.Open()
	assert.NoError(t, err)
	_ = client.GetSensorList(reserve.ReservationId, 0)

	//todo 获取到所有的sensorRecord，信息
	//	for _, sdrRecordAndValue := range sdrRecAndVallist {
	//		fmt.Println("type======", reflect.TypeOf(sdrRecordAndValue))
	//		//fmt.Println("sdrRecordAndValue======", sdrRecordAndValue.RecordType())
	//	}
}
