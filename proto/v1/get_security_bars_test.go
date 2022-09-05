package v1

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/go-test/deep"
)

func TestGetSecurityBarsResponse_Unmarshal(t *testing.T) {
	type WantData struct {
		Count  uint
		KLines []SecurityBarsRespItem
	}
	type args struct {
		category KLINE_TYPE
		filePath string
		data     []byte
	}
	tests := []struct {
		name     string
		args     args
		wantData WantData
		wantErr  bool
	}{
		{
			name: "stock_daily",
			args: args{category: KLINE_TYPE_DAILY, filePath: "test_datas/get_security_bars/stock_daily"},
			wantData: WantData{Count: 3, KLines: []SecurityBarsRespItem{
				{Open: 12.65, Close: 12.61, High: 12.79, Low: 12.58, Vol: 8.6198192e+07, Amount: 1.09266624e+09, Year: 2022, Month: 9, Day: 1, Hour: 15, Minute: 0, Datetime: "2022-09-01 15:00"},
				{Open: 12.62, Close: 12.51, High: 12.69, Low: 12.43, Vol: 7.863628e+07, Amount: 9.83433856e+08, Year: 2022, Month: 9, Day: 2, Hour: 15, Minute: 0, Datetime: "2022-09-02 15:00"},
				{Open: 12.46, Close: 12.49, High: 12.51, Low: 12.37, Vol: 3.4710424e+07, Amount: 4.31302656e+08, Year: 2022, Month: 9, Day: 5, Hour: 15, Minute: 0, Datetime: "2022-09-05 15:00"},
			}},
			wantErr: false,
		},
		{
			name: "stock_5min",
			args: args{category: KLINE_TYPE_5MIN, filePath: "test_datas/get_security_bars/stock_5min"},
			wantData: WantData{Count: 1, KLines: []SecurityBarsRespItem{
				{Open: 12.48, Close: 12.49, High: 12.49, Low: 12.48, Vol: 160728.0000076294, Amount: 1.5469365000004768e+06, Year: 2022, Month: 9, Day: 5, Hour: 13, Minute: 0, Datetime: "2022-09-05 13:00"},
			}},
			wantErr: false,
		},
		{
			name: "stock_15min",
			args: args{category: KLINE_TYPE_15MIN, filePath: "test_datas/get_security_bars/stock_15min"},
			wantData: WantData{Count: 1, KLines: []SecurityBarsRespItem{
				{Open: 12.49, Close: 12.49, High: 12.5, Low: 12.47, Vol: 542712.0000019073, Amount: 1.3324629e+07, Year: 2022, Month: 9, Day: 5, Hour: 13, Minute: 0, Datetime: "2022-09-05 13:00"},
			}},
			wantErr: false,
		},
		{
			name: "stock_30min",
			args: args{category: KLINE_TYPE_30MIN, filePath: "test_datas/get_security_bars/stock_30min"},
			wantData: WantData{Count: 1, KLines: []SecurityBarsRespItem{
				{Open: 12.48, Close: 12.49, High: 12.51, Low: 12.47, Vol: 1.3295480000004768e+06, Amount: 4.2811904e+07, Year: 2022, Month: 9, Day: 5, Hour: 13, Minute: 0, Datetime: "2022-09-05 13:00"},
			}},
			wantErr: false,
		},
		{
			name: "stock_1hour",
			args: args{category: KLINE_TYPE_1HOUR, filePath: "test_datas/get_security_bars/stock_1hour"},
			wantData: WantData{Count: 1, KLines: []SecurityBarsRespItem{
				{Open: 12.4, Close: 12.49, High: 12.51, Low: 12.4, Vol: 1.1227e+07, Amount: 1.3985632e+08, Year: 2022, Month: 9, Day: 5, Hour: 13, Minute: 0, Datetime: "2022-09-05 13:00"},
			}},
			wantErr: false,
		},
		{
			name: "stock_weekly",
			args: args{category: KLINE_TYPE_WEEKLY, filePath: "test_datas/get_security_bars/stock_weekly"},
			wantData: WantData{Count: 1, KLines: []SecurityBarsRespItem{
				{Open: 12.46, Close: 12.49, High: 12.51, Low: 12.37, Vol: 224606.0000076294, Amount: 4.42008e+08, Year: 2022, Month: 9, Day: 5, Hour: 15, Minute: 0, Datetime: "2022-09-05 15:00"},
			}},
			wantErr: false,
		},
		{
			name: "stock_monthly",
			args: args{category: KLINE_TYPE_MONTHLY, filePath: "test_datas/get_security_bars/stock_monthly"},
			wantData: WantData{Count: 1, KLines: []SecurityBarsRespItem{
				{Open: 12.65, Close: 12.49, High: 12.79, Low: 12.37, Vol: 1.4797340000019073e+06, Amount: 2.51810816e+09, Year: 2022, Month: 9, Day: 5, Hour: 15, Minute: 0, Datetime: "2022-09-05 15:00"},
			}},
			wantErr: false,
		},
		{
			name: "stock_3month",
			args: args{category: KLINE_TYPE_3MONTH, filePath: "test_datas/get_security_bars/stock_3month"},
			wantData: WantData{Count: 1, KLines: []SecurityBarsRespItem{
				{Open: 15, Close: 12.49, High: 15.27, Low: 11.94, Vol: 4.9036944e+07, Amount: 6.3628783616e+10, Year: 2022, Month: 9, Day: 5, Hour: 15, Minute: 0, Datetime: "2022-09-05 15:00"},
			}},
			wantErr: false,
		},
		{
			name: "stock_yearly",
			args: args{category: KLINE_TYPE_YEARLY, filePath: "test_datas/get_security_bars/stock_yearly"},
			wantData: WantData{Count: 1, KLines: []SecurityBarsRespItem{
				{Open: 16.48, Close: 12.49, High: 17.56, Low: 11.94, Vol: 1.85604464e+08, Amount: 2.74195185664e+11, Year: 2022, Month: 9, Day: 5, Hour: 15, Minute: 0, Datetime: "2022-09-05 15:00"},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			tt.args.data, err = ioutil.ReadFile(tt.args.filePath)
			if err != nil {
				t.Error(err)
				return
			}
			resp := &SecurityBarsResponse{
				category: tt.args.category,
			}
			if err := resp.Unmarshal(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("GetSecurityBars error = %v, wantErr %v", err, tt.wantErr)
			}
			if resp.Count != tt.wantData.Count {
				t.Errorf("GetSecurityBars count error = %v, wantErr %v", resp.Count, tt.wantData.Count)
			}
			for i, value := range resp.KLines {
				if diff := deep.Equal(value, tt.wantData.KLines[i]); diff != nil {
					t.Error(fmt.Sprintf("deep equal[%v]", i), diff)
				}
			}
		})
	}
}
