package v1

import (
	"fmt"

	"github.com/cyclegen-community/tdx-go/proto"
	"github.com/cyclegen-community/tdx-go/utils"
)

// 返回结构体
type TransactionData struct {
	Time      string  `json:"time"`
	Price     float64 `json:"price"`
	Vol       int     `json:"vol"`
	Num       int     `json:"num"`
	Buyorsell int     `json:"buyorsell"`
}

type GetTransactionDataResponse struct {
	Count               int               `json:"count"`
	TransactionDataList []TransactionData `json:"minute_time_data"`
}

// 原始返回结构体
type GetTransactionResponseRaw struct {
	Count                  uint16 `struc:"uint16,little,sizeof=TransactionDataRawList"`
	TransactionDataRawList []TransactionDataRaw
}

type TransactionDataRaw struct {
	Time      int       `struc:"uint16,little"`
	PriceRaw  PriceType `struc:"CustomType,little"`
	Vol       PriceType `struc:"CustomType,little"`
	Num       PriceType `struc:"CustomType,little"`
	Buyorsell PriceType `struc:"CustomType,little"`
	DontWant  PriceType `struc:"CustomType,little"` // 妈的，必须要一个字段名，如果字段名是：_，就会忽略解析
}

// 请求结构体
type GetTransactionDataRequest struct {
	Unknown1 []byte `struc:"[12]byte,little"`
	Market   Market `struc:"uint16,little"`
	Code     string `struc:"[6]byte,little"`
	Start    int    `struc:"uint16,little"`
	Count    int    `struc:"uint16,little"`
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

func (r *GetTransactionDataRequest) Marshal() ([]byte, error) {
	return proto.DefaultMarshal(r)
}

func (r *GetTransactionDataResponse) Unmarshal(data []byte) error {
	respR := GetTransactionResponseRaw{}
	if err := proto.DefaultUnmarshal(data, &respR); err != nil {
		return err
	}
	lastPrice := 0
	for i := 0; i < int(respR.Count); i++ {
		itemR := respR.TransactionDataRawList[i]
		item := TransactionData{}

		lastPrice = lastPrice + itemR.PriceRaw.getValue()
		item.Time = fmt.Sprintf("%02d:%02d", int(itemR.Time/60), int(itemR.Time%60))
		item.Price = float64(lastPrice) / 100
		item.Vol = itemR.Vol.getValue()
		item.Num = itemR.Num.getValue()
		item.Buyorsell = itemR.Buyorsell.getValue()

		r.TransactionDataList = append(r.TransactionDataList, item)
	}

	r.Count = int(respR.Count)
	fmt.Printf("%+v\n", r)
	return nil
}
