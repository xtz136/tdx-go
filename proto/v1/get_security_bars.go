package v1

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/cyclegen-community/tdx-go/proto"
	"github.com/lunixbochs/struc"
)

type dateTimeType struct {
	year   int
	month  int
	day    int
	hour   int
	minute int
}

func (c *dateTimeType) Pack(p []byte, opt *struc.Options) (int, error) {
	return 0, nil
}

func (c *dateTimeType) Unpack(r io.Reader, length int, opt *struc.Options) error {
	category := 3

	if category < 4 || category == 7 || category == 8 {
		data := make([]byte, 4)
		r.Read(data)
		zipday := binary.LittleEndian.Uint16(data[0:2])
		tminutes := binary.LittleEndian.Uint16(data[2:4])

		c.year = int((zipday >> 11) + 2004)
		c.month = int((zipday % 2048) / 100)
		c.day = int((zipday % 2048) % 100)
		c.hour = int(tminutes / 60)
		c.minute = int(tminutes % 60)
		return nil
	} else {
		data := make([]byte, 4)
		r.Read(data)
		// buf := bytes.NewBuffer(data)
		// var x int32
		// binary.Read(buf, binary.LittleEndian, &x)

		zipday := binary.LittleEndian.Uint32(data)

		c.year = int(zipday / 10000)
		c.month = int((zipday % 10000) / 100)
		c.day = int(zipday % 100)
		c.hour = 15
		c.minute = 0
		return nil
	}
}
func (c *dateTimeType) Size(opt *struc.Options) int {
	return -1
}
func (c *dateTimeType) String() string {
	return c.getValue()
}
func (c *dateTimeType) getParts() (int, int, int, int, int) {
	return c.year, c.month, c.day, c.hour, c.minute
}
func (c *dateTimeType) getValue() string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d", c.year, c.month, c.day, c.hour, c.minute)
}

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

type getSecurityBarsResponseItemRaw struct {
	Datetime  dateTimeType `struc:"CustomType,little";json:"datetime"`
	OpenDiff  PriceType    `struc:"CustomType,little";json:"open"`
	CloseDiff PriceType    `struc:"CustomType,little";json:"close"`
	HighDiff  PriceType    `struc:"CustomType,little";json:"close"`
	LowDiff   PriceType    `struc:"CustomType,little";json:"close"`
	Vol       VolumeType   `struc:"CustomType,little";json:"close"`
	DBVol     VolumeType   `struc:"CustomType,little";json:"close"`
}

type GetSecurityBarsResponseRaw struct {
	Count uint `struc:"uint16,little,sizeof=Lines";json:"count"`
	Lines []getSecurityBarsResponseItemRaw
}

type GetSecurityBarsRequestItem struct {
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
	Count  uint                         `json:"count"`
	KLines []GetSecurityBarsRequestItem `json:"datas"`

	// 存放这个变量，解析返回值需要用到
	category KLINE_TYPE
}

// 内部套用原始结构解析，外部为经过解析之后的响应信息
func (resp *GetSecurityBarsResponse) Unmarshal(data []byte) error {
	var respRaw GetSecurityBarsResponseRaw
	if err := proto.DefaultUnmarshal(data, &respRaw); err != nil {
		return err
	}
	resp.Count = respRaw.Count
	preDiffBase := 0
	for _, item := range respRaw.Lines {
		openDiff := item.OpenDiff.getValue()
		closeDiff := item.CloseDiff.getValue()
		highDiff := item.HighDiff.getValue()
		lowDiff := item.LowDiff.getValue()

		open := calPrice1000(openDiff, preDiffBase)
		openDiff += preDiffBase
		close := calPrice1000(openDiff, closeDiff)
		high := calPrice1000(openDiff, highDiff)
		low := calPrice1000(openDiff, lowDiff)

		year, month, day, hour, minue := item.Datetime.getParts()

		resp.KLines = append(resp.KLines, GetSecurityBarsRequestItem{
			Open:     open,
			Close:    close,
			High:     high,
			Low:      low,
			Vol:      item.Vol.getValue(),
			Amount:   item.DBVol.getValue(),
			Year:     year,
			Month:    month,
			Day:      day,
			Hour:     hour,
			Minute:   minue,
			Datetime: item.Datetime.getValue(),
		})

		preDiffBase = openDiff + closeDiff
	}
	return nil
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

func calPrice1000(base int, diff int) float64 {
	return float64(base+diff) / 1000
}
