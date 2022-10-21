package v1

import (
	"fmt"

	"github.com/cyclegen-community/tdx-go/proto"
	"github.com/cyclegen-community/tdx-go/utils"
)

type GetTransactionDataRequest struct {
	Unknown1 []byte `struc:"[12]byte"`
	Market   Market `struc:"uint16,little"`
	Code     string `struc:"[6]byte,little"`
	Start    int    `struc:"uint16,little"`
	Count    int    `struc:"uint16,little"`
}

func (r *GetTransactionDataRequest) Marshal() ([]byte, error) {
	return proto.DefaultMarshal(r)
}

type GetTransactionDataResponse struct {
	Count               int               `json:"count"`
	TransactionDataList []TransactionData `json:"minute_time_data"`
}

func (r *GetTransactionDataResponse) Unmarshal(data []byte) error {
	respR := GetTransactionResponseRaw{}
	if err := proto.DefaultUnmarshal(data, &respR); err != nil {
		return err
	}
	for i := 0; i < int(respR.Count); i++ {
		itemR := respR.TransactionDataRawList[i]
		item := TransactionData{}
		item.Price = float64(itemR.PriceRaw.getValue())
		item.Num = itemR.Num.getValue()
		item.Vol = itemR.Vol.getValue()
		item.Buyorsell = itemR.Buyorsell.getValue()
		r.TransactionDataList = append(r.TransactionDataList, item)
	}

	r.Count = int(respR.Count)
	fmt.Printf("%+v\n", r)
	return nil
}

type TransactionData struct {
	Time      string  `json:"time"`
	Price     float64 `json:"price"`
	Vol       int     `json:"vol"`
	Num       int     `json:"num"`
	Buyorsell int     `json:"buyorsell"`
}

type GetTransactionResponseRaw struct {
	Count                  uint16 `struc:"uint32,little,sizeof=TransactionDataRawList"`
	TransactionDataRawList []TransactionDataRaw
}

type TransactionDataRaw struct {
	Time      uint32    `struc:"uint32,little"`
	PriceRaw  PriceType `struc:"CustomType,little"`
	Reversed1 PriceType `struc:"CustomType,little"`
	Vol       PriceType `struc:"CustomType,little"`
	Num       PriceType `struc:"CustomType,little"`
	Buyorsell PriceType `struc:"CustomType,little"`
}

func NewGetTransactionRequest(market Market, code string, start int, count int) (*GetTransactionDataRequest, error) {
	request := &GetTransactionDataRequest{
		Unknown1: utils.HexString2Bytes("0c 17 08 01 01 01 0e 00 0e 00 c5 0f"),
		Market:   market,
		Code:     code,
		Start:    start,
		Count:    count,
	}
	return request, nil
}

func NewGetTransactionData(market Market, code string, start int, count int) (*GetTransactionDataRequest, *GetTransactionDataResponse, error) {
	var response GetTransactionDataResponse
	var request, err = NewGetTransactionRequest(market, code, start, count)
	return request, &response, err
}
