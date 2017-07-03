package ipmi

import (
	"net"
	//"reflect"
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
