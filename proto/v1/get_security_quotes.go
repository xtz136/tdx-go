package v1

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/cyclegen-community/tdx-go/proto"
)

func calPrice(base float64, diff float64) float64 {
	return float64(base+diff) / 100
}

func formatTime(timeStamp string) (string, error) {
	if len(timeStamp) < 6 {
		return "", nil
	}

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
type GetSecurityQuotesRequestParams struct {
	Market Market `struc:"uint8,little";json:"market"`
	Code   string `struct:"[6]byte";json:"code"`
}

func (req *GetSecurityQuotesRequestParams) Marshal() ([]byte, error) {
	return proto.DefaultMarshal(req)
}

type GetSecurityQuotesRequest struct {
	H1         int `struc:"uint16,little";json:"h1"`
	I2         int `struc:"uint32,little";json:"i2"`
	H3         int `struc:"uint16,little";json:"h3"`
	H4         int `struc:"uint16,little";json:"h4"`
	I5         int `struc:"uint32,little";json:"i5"`
	I6         int `struc:"uint32,little";json:"i6"`
	H7         int `struc:"uint16,little";json:"h7"`
	H8         int `struc:"uint16,little";json:"h8"`
	securities []GetSecurityQuotesRequestParams
}

// 请求包序列化输出
func (r *GetSecurityQuotesRequest) Marshal() ([]byte, error) {
	RequestByte, err := proto.DefaultMarshal(r)
	if err != nil {
		return nil, err
	}
	for _, p := range r.securities {
		paramByte, err := proto.DefaultMarshal(&p)
		if err != nil {
			return nil, err
		}
		RequestByte = append(RequestByte, paramByte...)
	}
	return RequestByte, nil
}

func NewGetSecurityQuotesRequest(securities []GetSecurityQuotesRequestParams) (*GetSecurityQuotesRequest, error) {
	stockLen := len(securities)
	if stockLen <= 0 {
		return nil, errors.New("securities must not be empty")
	}
	pkgDataLen := stockLen*7 + 12

	request := &GetSecurityQuotesRequest{
		H1:         0x10c,
		I2:         0x02006320,
		H3:         pkgDataLen,
		H4:         pkgDataLen,
		I5:         0x5053e,
		I6:         0,
		H7:         0,
		H8:         stockLen,
		securities: securities,
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

type GetSecurityQuotesResponse struct {
	Count          int
	SecurityQuotes []SecurityQuote
}

// 响应体原始结构

// 内部套用原始结构解析，外部为经过解析之后的响应信息
func (resp *GetSecurityQuotesResponse) Unmarshal(data []byte) error {
	// skip b1 cb
	data = data[2:]

	var err error
	var values []interface{}

	bp := new(BinaryPack)
	data, values, err = bp.UnPack([]string{"H"}, data)
	if err != nil {
		return err
	}
	resp.Count = values[0].(int)

	for i := 0; i < resp.Count; i++ {
		data, values, err = bp.UnPack([]string{"B", "6s", "H"}, data)
		if err != nil {
			return err
		}
		market := Market(values[0].(int))
		code := values[1].(string)
		active1 := values[2].(int)

		data1, price := bp.UnPackPrice(data)
		data1, lastCloseDiff := bp.UnPackPrice(data1)
		data1, openDiff := bp.UnPackPrice(data1)
		data1, highDiff := bp.UnPackPrice(data1)
		data1, lowDiff := bp.UnPackPrice(data1)
		data1, reversedBytes0 := bp.UnPackPrice(data1)
		data1, reversedBytes1 := bp.UnPackPrice(data1)
		data1, vol := bp.UnPackPrice(data1)
		data1, curVol := bp.UnPackPrice(data1)

		data1, amount, err := bp.UnPackAmount(data1)
		if err != nil {
			return err
		}

		data1, sVol := bp.UnPackPrice(data1)
		data1, bVol := bp.UnPackPrice(data1)
		data1, reversedBytes2 := bp.UnPackPrice(data1)
		data1, reversedBytes3 := bp.UnPackPrice(data1)

		data1, bid1 := bp.UnPackPrice(data1)
		data1, ask1 := bp.UnPackPrice(data1)
		data1, bidVol1 := bp.UnPackPrice(data1)
		data1, askVol1 := bp.UnPackPrice(data1)

		data1, bid2 := bp.UnPackPrice(data1)
		data1, ask2 := bp.UnPackPrice(data1)
		data1, bidVol2 := bp.UnPackPrice(data1)
		data1, askVol2 := bp.UnPackPrice(data1)

		data1, bid3 := bp.UnPackPrice(data1)
		data1, ask3 := bp.UnPackPrice(data1)
		data1, bidVol3 := bp.UnPackPrice(data1)
		data1, askVol3 := bp.UnPackPrice(data1)

		data1, bid4 := bp.UnPackPrice(data1)
		data1, ask4 := bp.UnPackPrice(data1)
		data1, bidVol4 := bp.UnPackPrice(data1)
		data1, askVol4 := bp.UnPackPrice(data1)

		data1, bid5 := bp.UnPackPrice(data1)
		data1, ask5 := bp.UnPackPrice(data1)
		data1, bidVol5 := bp.UnPackPrice(data1)
		data1, askVol5 := bp.UnPackPrice(data1)

		data1, values, err = bp.UnPack([]string{"H"}, data1)
		if err != nil {
			return err
		}
		reversedBytes4 := values[0].(int)

		data1, reversedBytes5 := bp.UnPackPrice(data1)
		data1, reversedBytes6 := bp.UnPackPrice(data1)
		data1, reversedBytes7 := bp.UnPackPrice(data1)
		data1, reversedBytes8 := bp.UnPackPrice(data1)

		data1, values, err = bp.UnPack([]string{"h", "H"}, data1)
		if err != nil {
			return err
		}
		reversedBytes9 := values[0].(int)
		active2 := values[1].(int)
		data = data1

		servertime, err := formatTime(fmt.Sprintf("%d", int(reversedBytes0)-int(reversedBytes1)))
		if err != nil {
			return err
		}

		quote := SecurityQuote{
			Market:         market,
			Code:           code,
			Active1:        active1,
			Price:          calPrice(price, 0),
			LastClose:      calPrice(price, lastCloseDiff),
			Open:           calPrice(price, openDiff),
			High:           calPrice(price, highDiff),
			Low:            calPrice(price, lowDiff),
			Servertime:     servertime,
			ReversedBytes0: int(reversedBytes0),
			ReversedBytes1: int(reversedBytes1),
			Vol:            int(vol),
			CurVol:         int(curVol),
			Amount:         amount,
			SVol:           int(sVol),
			BVol:           int(bVol),
			ReversedBytes2: int(reversedBytes2),
			ReversedBytes3: int(reversedBytes3),
			Bid1:           calPrice(price, bid1),
			Ask1:           calPrice(price, ask1),
			BidVol1:        int(bidVol1),
			AskVol1:        int(askVol1),
			Bid2:           calPrice(price, bid2),
			Ask2:           calPrice(price, ask2),
			BidVol2:        int(bidVol2),
			AskVol2:        int(askVol2),
			Bid3:           calPrice(price, bid3),
			Ask3:           calPrice(price, ask3),
			BidVol3:        int(bidVol3),
			AskVol3:        int(askVol3),
			Bid4:           calPrice(price, bid4),
			Ask4:           calPrice(price, ask4),
			BidVol4:        int(bidVol4),
			AskVol4:        int(askVol4),
			Bid5:           calPrice(price, bid5),
			Ask5:           calPrice(price, ask5),
			BidVol5:        int(bidVol5),
			AskVol5:        int(askVol5),
			ReversedBytes4: []int{reversedBytes4},
			ReversedBytes5: int(reversedBytes5),
			ReversedBytes6: int(reversedBytes6),
			ReversedBytes7: int(reversedBytes7),
			ReversedBytes8: int(reversedBytes8),
			ReversedBytes9: float64(reversedBytes9) / 100.0,
			Active2:        active2,
		}

		resp.SecurityQuotes = append(resp.SecurityQuotes, quote)
	}

	// fmt.Printf("%+v\n", resp)

	return nil
}

func NewGetSecurityQuotes(securities []GetSecurityQuotesRequestParams) (*GetSecurityQuotesRequest, *GetSecurityQuotesResponse, error) {
	var response GetSecurityQuotesResponse
	var request, err = NewGetSecurityQuotesRequest(securities)
	return request, &response, err
}
