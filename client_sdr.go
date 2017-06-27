package ipmi

// RepositoryInfo get the Repository Info of the SDR
// func (c *Client) GetSDRRepositoryInfo() (*SDRRepositoryInfoResponse, error) {
// 	req := &Request{
// 		NetworkFunctionStorge,
// 		CommandGetSDRRepositoryInfo,
// 		&SDRRepositoryInfoRequest{},
// 	}
// 	res := &SDRRepositoryInfoResponse{}
// 	return res, c.Send(req, res)
// }
// RepositoryInfo get sdr

func (c *Client) GetSDR(reservationID uint16, recordID uint16){
	
	req := &Request{
		NetworkFunctionStorge,
		CommandGetSDR,
		&GetSDRCommandRequest{
			ReservationID: reservationID,
			RecordID:recordID,
		},
	}

	res := &GetSDRCommandResponse{
		ReadData:SDRFullSensor{},
	}
	c.Send(req, res)

// type SDRFullSensor struct {
// 	SDRRecordHeader
// 	sdrFullSensorFields
// 	deviceId string
// }
// type SDRRecordHeader struct {
// 	recordId   uint16
// 	SDRVersion uint8
// 	rtype      SDRRecordType
// }
// type sdrFullSensorFields struct { //size 42
// 	SensorOwnerId        uint8
// 	SensorOwnerLUN       uint8
// 	SensorNumber         uint8
// 	EntityId             uint8
// 	EntityIns            uint8
// 	SensorInit           uint8
// 	SensorCap            uint8
// 	SensorType           SDRSensorType
// 	ReadingType          SDRSensorReadingType
// 	AssertionEventMask   uint16
// 	DeassertionEventMask uint16
// 	DiscreteReadingMask  uint16
// 	Unit                 uint8
// 	BaseUnit             uint8
// 	ModifierUnit         uint8
// 	Linearization        uint8
// 	MTol                 uint16
// 	Bacc                 uint32
// 	AnalogFlag           uint8
// 	NominalReading       uint8
// 	NormalMax            uint8
// 	NormalMin            uint8
// 	SensorMax            uint8
// 	SensorMin            uint8
// 	U_NR                 uint8
// 	U_C                  uint8
// 	U_NC                 uint8
// 	L_NR                 uint8
// 	L_C                  uint8
// 	L_NC                 uint8
// 	PositiveHysteresis   uint8
// 	NegativeHysteresis   uint8
// 	reserved             [2]byte
// 	OEM                  uint8
// }


	//nextRecordID := res.NextRecordID
	//readData := res.ReadData
	//sensorType := readData.SensorType



}







