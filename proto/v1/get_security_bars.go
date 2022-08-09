package v1

import (
	"fmt"

	"github.com/cyclegen-community/tdx-go/proto"
)

// 请求包结构
type GetSecurityBarsRequestParams struct {
	Market   int    `struc:"uint16,little";json:"market"`
	Code     string `struc:"[6]byte,little";json:"code"`
	Category int    `struc:"uint16,little";json:"category"`
	Start    int    `struc:"uint16,little";json:"start"`
	Count    int    `struc:"uint16,little";json:"count"`
}

func (req *GetSecurityBarsRequestParams) Marshal() ([]byte, error) {
	return proto.DefaultMarshal(req)
}

type GetSecurityBarsRequest struct {
	H1       int        `struc:"uint16,little";json:"h1"`
	I2       int        `struc:"uint32,little";json:"i2"`
	H3       int        `struc:"uint16,little";json:"h3"`
	H4       int        `struc:"uint16,little";json:"h4"`
	H5       int        `struc:"uint16,little";json:"h5"`
	Market   Market     `struc:"uint16,little";json:"market"`
	Code     string     `struc:"[6]byte,little";json:"code"`
	Category KLINE_TYPE `struc:"uint16,little";json:"category"`
	H9       int        `struc:"uint16,little";json:"h9"`
	Start    int        `struc:"uint16,little";json:"start"`
	Count    int        `struc:"uint16,little";json:"count"`
	I12      int        `struc:"uint32,little";json:"i12"`
	I13      int        `struc:"uint32,little";json:"i13"`
	H14      int        `struc:"uint16,little";json:"h14"`
}

// 请求包序列化输出
func (r *GetSecurityBarsRequest) Marshal() ([]byte, error) {
	return proto.DefaultMarshal(r)
}

type Kline struct {
	Open     float64 `json:"open"`
	Close    float64 `json:"close"`
	High     float64 `json:"high"`
	Low      float64 `json:"low"`
	Vol      float64 `json:"vol"`
	Amount   float64 `json:"amount"`
	Year     int     `json:"year"`
	Month    int     `json:"month"`
	Day      int     `json:"day"`
	Hour     int     `json:"hour"`
	Minute   int     `json:"minute"`
	Datetime string  `json:"datetime"`
}

// 响应包结构
type GetSecurityBarsResponse struct {
	klines []Kline

	// 存放这个变量，解析返回值需要用到
	category KLINE_TYPE
}

// 内部套用原始结构解析，外部为经过解析之后的响应信息
func (resp *GetSecurityBarsResponse) Unmarshal(data []byte) error {
	var err error
	var values []interface{}
	bp := new(BinaryPack)
	data, values, err = bp.UnPack([]string{"H"}, data)
	if err != nil {
		return err
	}

	preDiffBase := float64(0)
	count := values[0].(int)
	for i := 0; i < count; i++ {
		dataTemp, year, month, day, hour, minute := getDatetime(bp, int(resp.category), data)
		dataTemp, priceOpenDiff := bp.UnPackPrice(dataTemp)
		dataTemp, priceCloseDiff := bp.UnPackPrice(dataTemp)
		dataTemp, priceHighDiff := bp.UnPackPrice(dataTemp)
		dataTemp, priceLowDiff := bp.UnPackPrice(dataTemp)

		dataTemp, vol, err := bp.UnPackAmount(dataTemp)
		if err != nil {
			return err
		}

		dataTemp, dbVol, err := bp.UnPackAmount(dataTemp)
		if err != nil {
			return err
		}

		open := calPrice1000(priceOpenDiff, preDiffBase)
		priceOpenDiff += preDiffBase
		close := calPrice1000(priceOpenDiff, priceCloseDiff)
		high := calPrice1000(priceOpenDiff, priceHighDiff)
		low := calPrice1000(priceOpenDiff, priceLowDiff)

		kline := Kline{
			Open:     open,
			Close:    close,
			High:     high,
			Low:      low,
			Vol:      vol,
			Amount:   dbVol,
			Year:     year,
			Month:    month,
			Day:      day,
			Hour:     hour,
			Minute:   minute,
			Datetime: fmt.Sprintf("%d-%02d-%02d %02d:%02d", year, month, day, hour, minute),
		}
		resp.klines = append(resp.klines, kline)

		preDiffBase = priceOpenDiff + priceCloseDiff
		data = dataTemp
	}

	fmt.Printf("%+v\n", resp)

	return nil
}

// func getPrice(b []byte, pos *int) int {
// 	/*
// 		    0x7f与常量做与运算实质是保留常量（转换为二进制形式）的后7位数，既取值区间为[0,127]
// 		    0x3f与常量做与运算实质是保留常量（转换为二进制形式）的后6位数，既取值区间为[0,63]

// 			0x80 1000 0000
// 			0x7f 0111 1111
// 			0x40  100 0000
// 			0x3f  011 1111
// 	*/
// 	posByte := 6
// 	bData := b[*pos]
// 	data := int(bData & 0x3f)
// 	bSign := false
// 	if (bData & 0x40) > 0 {
// 		bSign = true
// 	}

// 	if (bData & 0x80) > 0 {
// 		for {
// 			*pos += 1
// 			bData = b[*pos]
// 			data += (int(bData&0x7f) << posByte)

// 			posByte += 7

// 			if (bData & 0x80) <= 0 {
// 				break
// 			}
// 		}
// 	}
// 	*pos++

// 	if bSign {
// 		data = -data
// 	}
// 	return data
// }

func getDatetime(bp *BinaryPack, category int, msg []byte) (restMsg []byte, year int, month int, day int, hour int, minute int) {
	hour = 15
	if category < 4 || category == 7 || category == 8 {
		msg, values, _ := bp.UnPack([]string{"H", "H"}, msg)
		zipday := values[0].(int)
		tminutes := values[1].(int)

		year = int((zipday >> 11) + 2004)
		month = int((zipday % 2048) / 100)
		day = int((zipday % 2048) % 100)
		hour = int(tminutes / 60)
		minute = int(tminutes % 60)
		return msg, year, month, day, hour, minute

	} else {
		msg, values, _ := bp.UnPack([]string{"I"}, msg)
		zipday := values[0].(int)

		year = int(zipday / 10000)
		month = int((zipday % 10000) / 100)
		day = int(zipday % 100)
		return msg, year, month, day, 15, 0
	}
}

// todo: 检测market是否为合法值
func NewGetSecurityBarsRequest(market Market, code string, category KLINE_TYPE, start int, count int) (*GetSecurityBarsRequest, error) {
	request := &GetSecurityBarsRequest{
		H1:       0x10c,
		I2:       0x01016408,
		H3:       0x1c,
		H4:       0x1c,
		H5:       0x052d,
		Market:   market,
		Code:     code,
		Category: category,
		H9:       1,
		Start:    start,
		Count:    count,
		I12:      0,
		I13:      0,
		H14:      0,
	}

	return request, nil
}

type KLINE_TYPE int

const (
	/*
		K线种类
		# K 线种类
		# 0 -   5 分钟K 线
		# 1 -   15 分钟K 线
		# 2 -   30 分钟K 线
		# 3 -   1 小时K 线
		# 4 -   日K 线
		# 5 -   周K 线
		# 6 -   月K 线
		# 7 -   1 分钟
		# 8 -   1 分钟K 线
		# 9 -   日K 线
		# 10 -  季K 线
		# 11 -  年K 线
	*/
	KLINE_TYPE_5MIN      KLINE_TYPE = iota
	KLINE_TYPE_15MIN                = 1
	KLINE_TYPE_30MIN                = 2
	KLINE_TYPE_1HOUR                = 3
	KLINE_TYPE_DAILY                = 4
	KLINE_TYPE_WEEKLY               = 5
	KLINE_TYPE_MONTHLY              = 6
	KLINE_TYPE_EXHQ_1MIN            = 7
	KLINE_TYPE_1MIN                 = 8
	KLINE_TYPE_RI_K                 = 9
	KLINE_TYPE_3MONTH               = 10
	KLINE_TYPE_YEARLY               = 11
)

func NewGetSecurityBars(market Market, code string, category KLINE_TYPE, start int, count int) (*GetSecurityBarsRequest, *GetSecurityBarsResponse, error) {
	var response GetSecurityBarsResponse
	var request, err = NewGetSecurityBarsRequest(market, code, category, start, count)
	response.category = category
	return request, &response, err
}

func calPrice1000(base float64, diff float64) float64 {
	return float64(base+diff) / 1000
}
