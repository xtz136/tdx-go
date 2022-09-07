package v1

import (
	"github.com/cyclegen-community/tdx-go/proto"
	"github.com/cyclegen-community/tdx-go/utils"
)

type GetMinuteTimeDataRequest struct {
	Unknown1 []byte `struc:"[12]byte"`
	Market   Market `struc:"uint16,little"`
	Code     string `struc:"[6]byte,little"`
	Unknown2 uint32 `struc:"uint32,little"`
}

func (r *GetMinuteTimeDataRequest) Marshal() ([]byte, error) {
	return proto.DefaultMarshal(r)
}

type GetMinuteTimeDataResponse struct {
	Count              int              `json:"count"`
	MinuteTimeDataList []MinuteTimeData `json:"minute_time_data"`
}

func (r *GetMinuteTimeDataResponse) Unmarshal(data []byte) error {
	respR := GetMinuteTimeResponseRaw{}
	if err := proto.DefaultUnmarshal(data, &respR); err != nil {
		return err
	}
	lastPrice := 0
	for i := 0; i < int(respR.Count); i++ {
		itemR := respR.MinuteTimeDataRawList[i]
		lastPrice += itemR.PriceRaw.getValue()
		item := MinuteTimeData{Price: float64(lastPrice) / 100, Vol: itemR.Vol.getValue()}
		r.MinuteTimeDataList = append(r.MinuteTimeDataList, item)
	}

	r.Count = int(respR.Count)
	// fmt.Printf("%+v\n", r)
	return nil
}

type MinuteTimeData struct {
	Price float64 `json:"price"`
	Vol   int     `json:"vol"`
}

type GetMinuteTimeResponseRaw struct {
	Count                 uint16 `struc:"uint32,little,sizeof=MinuteTimeDataRawList"`
	MinuteTimeDataRawList []MinuteTimeDataRaw
}

type MinuteTimeDataRaw struct {
	PriceRaw  PriceType `struc:"CustomType,little"`
	Reversed1 PriceType `struc:"CustomType,little"`
	Vol       PriceType `struc:"CustomType,little"`
}

func NewGetMinuteTimeRequest(market Market, code string) (*GetMinuteTimeDataRequest, error) {
	request := &GetMinuteTimeDataRequest{
		Unknown1: utils.HexString2Bytes("0c 1b 08 00 01 01 0e 00 0e 00 1d 05"),
		Market:   market,
		Code:     code,
		Unknown2: 0,
	}
	return request, nil
}

func NewGetMinuteTimeData(market Market, code string) (*GetMinuteTimeDataRequest, *GetMinuteTimeDataResponse, error) {
	var response GetMinuteTimeDataResponse
	var request, err = NewGetMinuteTimeRequest(market, code)
	return request, &response, err
}
