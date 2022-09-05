package v1

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/cyclegen-community/tdx-go/proto"
)

func calPrice(base int, diff int) float64 {
	return float64(base+diff) / 100
}

func formatTime(timeStamp string) (string, error) {
	if len(timeStamp) < 6 {
		return "", nil
	}

	// 0-8点的时候，前面少了一个0
	if len(timeStamp) < 8 {
		timeStamp = fmt.Sprintf("%08s", timeStamp)
	}

	time := timeStamp[:2] + ":"

	t2, err := strconv.Atoi(timeStamp[4:8])
	if err != nil {
		return "", err
	}
	time += timeStamp[2:4] + ":"
	time += fmt.Sprintf("%6.3f", float64(t2)*60/10000.0)

	return time, nil
}

// 请求包结构
type SecurityQuotesRequestParams struct {
	Market Market `struc:"uint8,little"`
	Code   string `struct:"[6]byte"`
}

func (req *SecurityQuotesRequestParams) Marshal() ([]byte, error) {
	return proto.DefaultMarshal(req)
}

type SecurityQuotesRequest struct {
	H1     int `struc:"uint16,little"`
	I2     int `struc:"uint32,little"`
	H3     int `struc:"uint16,little"`
	H4     int `struc:"uint16,little"`
	I5     int `struc:"uint32,little"`
	I6     int `struc:"uint32,little"`
	H7     int `struc:"uint16,little"`
	H8     int `struc:"uint16,little"`
	params []SecurityQuotesRequestParams
}

// 请求包序列化输出
func (r *SecurityQuotesRequest) Marshal() ([]byte, error) {
	RequestByte, err := proto.DefaultMarshal(r)
	if err != nil {
		return nil, err
	}
	for _, p := range r.params {
		paramByte, err := proto.DefaultMarshal(&p)
		if err != nil {
			return nil, err
		}
		RequestByte = append(RequestByte, paramByte...)
	}
	return RequestByte, nil
}

func NewGetSecurityQuotesRequest(params []SecurityQuotesRequestParams) (*SecurityQuotesRequest, error) {
	stockLen := len(params)
	if stockLen <= 0 {
		return nil, errors.New("securities must not be empty")
	}
	pkgDataLen := stockLen*7 + 12

	request := &SecurityQuotesRequest{
		H1:     0x10c,
		I2:     0x02006320,
		H3:     pkgDataLen,
		H4:     pkgDataLen,
		I5:     0x5053e,
		I6:     0,
		H7:     0,
		H8:     stockLen,
		params: params,
	}

	return request, nil
}

// 响应包结构
type SecurityQuote struct {
	Market         Market  `json:"market"`
	Code           string  `json:"code"`
	Active1        int     `json:"active1"`
	Price          float64 `json:"price"`
	LastClose      float64 `json:"last_close"`
	Open           float64 `json:"open"`
	High           float64 `json:"high"`
	Low            float64 `json:"low"`
	Servertime     string  `json:"servertime"`
	ReversedBytes0 int     `json:"reversed_bytes0"`
	ReversedBytes1 int     `json:"reversed_bytes1"`
	Vol            int     `json:"vol"`
	CurVol         int     `json:"cur_vol"`
	Amount         float64 `json:"amount"`
	SVol           int     `json:"s_vol"`
	BVol           int     `json:"b_vol"`
	ReversedBytes2 int     `json:"reversed_bytes2"`
	ReversedBytes3 int     `json:"reversed_bytes3"`
	Bid1           float64 `json:"bid1"`
	Ask1           float64 `json:"ask1"`
	BidVol1        int     `json:"bid_vol1"`
	AskVol1        int     `json:"ask_vol1"`
	Bid2           float64 `json:"bid2"`
	Ask2           float64 `json:"ask2"`
	BidVol2        int     `json:"bid_vol2"`
	AskVol2        int     `json:"ask_vol2"`
	Bid3           float64 `json:"bid3"`
	Ask3           float64 `json:"ask3"`
	BidVol3        int     `json:"bid_vol3"`
	AskVol3        int     `json:"ask_vol3"`
	Bid4           float64 `json:"bid4"`
	Ask4           float64 `json:"ask4"`
	BidVol4        int     `json:"bid_vol4"`
	AskVol4        int     `json:"ask_vol4"`
	Bid5           float64 `json:"bid5"`
	Ask5           float64 `json:"ask5"`
	BidVol5        int     `json:"bid_vol5"`
	AskVol5        int     `json:"ask_vol5"`
	ReversedBytes4 []int   `json:"reversed_bytes4"`
	ReversedBytes5 int     `json:"reversed_bytes5"`
	ReversedBytes6 int     `json:"reversed_bytes6"`
	ReversedBytes7 int     `json:"reversed_bytes7"`
	ReversedBytes8 int     `json:"reversed_bytes8"`
	ReversedBytes9 float64 `json:"reversed_bytes9"`
	Active2        int     `json:"active2"`
}

type SecurityQuotesResponse struct {
	Count          int
	SecurityQuotes []SecurityQuote
}

// 响应体原始结构
type SecurityQuoteRaw struct {
	Market         Market    `struc:"uint8,little"`
	Code           string    `struc:"[6]byte,little"`
	Active1        int       `struc:"uint16,little"`
	Price          PriceType `struc:"CustomType,little"`
	LastCloseDiff  PriceType `struc:"CustomType,little"`
	OpenDiff       PriceType `struc:"CustomType,little"`
	HighDiff       PriceType `struc:"CustomType,little"`
	LowDiff        PriceType `struc:"CustomType,little"`
	ReversedBytes0 PriceType `struc:"CustomType,little"`
	ReversedBytes1 PriceType `struc:"CustomType,little"`
	Vol            PriceType `struc:"CustomType,little"`
	CurVol         PriceType `struc:"CustomType,little"`
	Amount         int       `struc:"uint32,little"`
	SVol           PriceType `struc:"CustomType,little"`
	BVol           PriceType `struc:"CustomType,little"`
	ReversedBytes2 PriceType `struc:"CustomType,little"`
	ReversedBytes3 PriceType `struc:"CustomType,little"`
	Bid1           PriceType `struc:"CustomType,little"`
	Ask1           PriceType `struc:"CustomType,little"`
	BidVol1        PriceType `struc:"CustomType,little"`
	AskVol1        PriceType `struc:"CustomType,little"`
	Bid2           PriceType `struc:"CustomType,little"`
	Ask2           PriceType `struc:"CustomType,little"`
	BidVol2        PriceType `struc:"CustomType,little"`
	AskVol2        PriceType `struc:"CustomType,little"`
	Bid3           PriceType `struc:"CustomType,little"`
	Ask3           PriceType `struc:"CustomType,little"`
	BidVol3        PriceType `struc:"CustomType,little"`
	AskVol3        PriceType `struc:"CustomType,little"`
	Bid4           PriceType `struc:"CustomType,little"`
	Ask4           PriceType `struc:"CustomType,little"`
	BidVol4        PriceType `struc:"CustomType,little"`
	AskVol4        PriceType `struc:"CustomType,little"`
	Bid5           PriceType `struc:"CustomType,little"`
	Ask5           PriceType `struc:"CustomType,little"`
	BidVol5        PriceType `struc:"CustomType,little"`
	AskVol5        PriceType `struc:"CustomType,little"`
	ReversedBytes4 int       `struc:"uint16,little"`
	ReversedBytes5 PriceType `struc:"CustomType,little"`
	ReversedBytes6 PriceType `struc:"CustomType,little"`
	ReversedBytes7 PriceType `struc:"CustomType,little"`
	ReversedBytes8 PriceType `struc:"CustomType,little"`
	ReversedBytes9 int       `struc:"uint16,little"`
	Active2        int       `struc:"uint16,little"`
}

type SecurityQuoteResponseRaw struct {
	Count             int `struc:"uint16,little,sizeof=SecurityQuoteRaws"`
	SecurityQuoteRaws []SecurityQuoteRaw
}

// 内部套用原始结构解析，外部为经过解析之后的响应信息
func (resp *SecurityQuotesResponse) Unmarshal(data []byte) error {
	// skip b1 cb
	data = data[2:]

	var respR SecurityQuoteResponseRaw
	if err := proto.DefaultUnmarshal(data, &respR); err != nil {
		return err
	}

	resp.Count = respR.Count

	for i := 0; i < respR.Count; i++ {
		itemR := respR.SecurityQuoteRaws[i]
		servertime, err := formatTime(fmt.Sprintf("%d", int(itemR.ReversedBytes0.getValue())-int(itemR.ReversedBytes1.getValue())))
		if err != nil {
			return err
		}

		price := itemR.Price.getValue()
		quote := SecurityQuote{
			Market:         itemR.Market,
			Code:           itemR.Code,
			Active1:        itemR.Active1,
			Price:          calPrice(price, 0),
			LastClose:      calPrice(price, itemR.LastCloseDiff.getValue()),
			Open:           calPrice(price, itemR.OpenDiff.getValue()),
			High:           calPrice(price, itemR.HighDiff.getValue()),
			Low:            calPrice(price, itemR.LowDiff.getValue()),
			Servertime:     servertime,
			ReversedBytes0: itemR.ReversedBytes0.getValue(),
			ReversedBytes1: itemR.ReversedBytes1.getValue(),
			Vol:            itemR.Vol.getValue(),
			CurVol:         itemR.CurVol.getValue(),
			Amount:         ParseVolume(itemR.Amount),
			SVol:           itemR.SVol.getValue(),
			BVol:           itemR.BVol.getValue(),
			ReversedBytes2: itemR.ReversedBytes2.getValue(),
			ReversedBytes3: itemR.ReversedBytes3.getValue(),
			Bid1:           calPrice(price, itemR.Bid1.getValue()),
			Ask1:           calPrice(price, itemR.Ask1.getValue()),
			BidVol1:        itemR.BidVol1.getValue(),
			AskVol1:        itemR.AskVol1.getValue(),
			Bid2:           calPrice(price, itemR.Bid2.getValue()),
			Ask2:           calPrice(price, itemR.Ask2.getValue()),
			BidVol2:        itemR.BidVol2.getValue(),
			AskVol2:        itemR.AskVol2.getValue(),
			Bid3:           calPrice(price, itemR.Bid3.getValue()),
			Ask3:           calPrice(price, itemR.Ask3.getValue()),
			BidVol3:        itemR.BidVol3.getValue(),
			AskVol3:        itemR.AskVol3.getValue(),
			Bid4:           calPrice(price, itemR.Bid4.getValue()),
			Ask4:           calPrice(price, itemR.Ask4.getValue()),
			BidVol4:        itemR.BidVol4.getValue(),
			AskVol4:        itemR.AskVol4.getValue(),
			Bid5:           calPrice(price, itemR.Bid5.getValue()),
			Ask5:           calPrice(price, itemR.Ask5.getValue()),
			BidVol5:        itemR.BidVol5.getValue(),
			AskVol5:        itemR.AskVol5.getValue(),
			ReversedBytes4: []int{itemR.ReversedBytes4},
			ReversedBytes5: itemR.ReversedBytes5.getValue(),
			ReversedBytes6: itemR.ReversedBytes6.getValue(),
			ReversedBytes7: itemR.ReversedBytes7.getValue(),
			ReversedBytes8: itemR.ReversedBytes8.getValue(),
			ReversedBytes9: float64(itemR.ReversedBytes9) / 100.0,
			Active2:        itemR.Active2,
		}

		resp.SecurityQuotes = append(resp.SecurityQuotes, quote)
	}

	// fmt.Printf("%+v\n", resp)

	return nil
}

func NewGetSecurityQuotes(params []SecurityQuotesRequestParams) (*SecurityQuotesRequest, *SecurityQuotesResponse, error) {
	var response SecurityQuotesResponse
	var request, err = NewGetSecurityQuotesRequest(params)
	return request, &response, err
}
