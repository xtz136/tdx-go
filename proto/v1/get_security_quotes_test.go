package v1

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/go-test/deep"
)

func TestGetSecurityQuotesResponse_Unmarshal(t *testing.T) {
	type fields struct {
		Count          int
		SecurityQuotes []SecurityQuote
	}
	type args struct {
		filePath string
		data     []byte
	}
	tests := []struct {
		name     string
		wantData fields
		args     args
		wantErr  bool
	}{
		{
			name: "stock",
			args: args{filePath: "test_datas/get_security_quotes/stock"},
			wantData: fields{Count: 2, SecurityQuotes: []SecurityQuote{
				{Market: 0, Code: "000001", Active1: 4389, Price: 12.51, LastClose: 12.61, Open: 12.62, High: 12.69, Low: 12.43, Servertime: "15:00: 0.006", ReversedBytes0: 14998750, ReversedBytes1: -1251, Vol: 786362, CurVol: 10006, Amount: 9.83433856e+08, SVol: 446986, BVol: 339377, ReversedBytes2: 0, ReversedBytes3: 46626, Bid1: 12.5, Ask1: 12.51, BidVol1: 4093, AskVol1: 2418, Bid2: 12.49, Ask2: 12.52, BidVol2: 2050, AskVol2: 1978, Bid3: 12.48, Ask3: 12.53, BidVol3: 1897, AskVol3: 2298, Bid4: 12.47, Ask4: 12.54, BidVol4: 2073, AskVol4: 3487, Bid5: 12.46, Ask5: 12.55, BidVol5: 2471, AskVol5: 2956, ReversedBytes4: []int{5206}, ReversedBytes5: 1, ReversedBytes6: -24, ReversedBytes7: -23, ReversedBytes8: 15, ReversedBytes9: 0.08, Active2: 4389},
				{Market: 0, Code: "000002", Active1: 4673, Price: 16.8, LastClose: 16.84, Open: 16.93, High: 16.95, Low: 16.57, Servertime: "15:00: 0.006", ReversedBytes0: 14998321, ReversedBytes1: -1680, Vol: 696804, CurVol: 6635, Amount: 1.165902336e+09, SVol: 314991, BVol: 381813, ReversedBytes2: -1, ReversedBytes3: 69442, Bid1: 16.8, Ask1: 16.81, BidVol1: 790, AskVol1: 598, Bid2: 16.79, Ask2: 16.82, BidVol2: 122, AskVol2: 217, Bid3: 16.78, Ask3: 16.83, BidVol3: 650, AskVol3: 586, Bid4: 16.77, Ask4: 16.84, BidVol4: 365, AskVol4: 263, Bid5: 16.76, Ask5: 16.85, BidVol5: 325, AskVol5: 442, ReversedBytes4: []int{726}, ReversedBytes5: 0, ReversedBytes6: 17, ReversedBytes7: 35, ReversedBytes8: -20, ReversedBytes9: -0.17, Active2: 4673},
			}},
			wantErr: false,
		},
		{
			name: "cb",
			args: args{filePath: "test_datas/get_security_quotes/cb"},
			wantData: fields{Count: 1, SecurityQuotes: []SecurityQuote{
				{Market: 0, Code: "128072", Active1: 754, Price: 11659, LastClose: 11620.9, Open: 11598, High: 11670, Low: 11590.3, Servertime: "11:05:26.994", ReversedBytes0: 1104284084, ReversedBytes1: -1165900, Vol: 16476, CurVol: 25, Amount: 1.9180354e+07, SVol: 7551, BVol: 8925, ReversedBytes2: 620, ReversedBytes3: 0, Bid1: 11652.7, Ask1: 11659.6, BidVol1: 7, AskVol1: 1, Bid2: 11647.2, Ask2: 11660, BidVol2: 9, AskVol2: 16, Bid3: 11647, Ask3: 11664.4, BidVol3: 10, AskVol3: 2, Bid4: 11645, Ask4: 11665, BidVol4: 2, AskVol4: 10, Bid5: 11644, Ask5: 11666, BidVol5: 1, AskVol5: 10, ReversedBytes4: []int{16396}, ReversedBytes5: 1, ReversedBytes6: 630, ReversedBytes7: -28190, ReversedBytes8: 5140, ReversedBytes9: 0, Active2: 754},
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
			resp := &SecurityQuotesResponse{}
			if err := resp.Unmarshal(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("GetSecurityQuotes error = %v, wantErr %v", err, tt.wantErr)
			}
			if resp.Count != tt.wantData.Count {
				t.Errorf("GetSecurityQuotes count error = %v, wantErr %v", resp.Count, tt.wantData.Count)
			}
			for i, value := range resp.SecurityQuotes {
				if diff := deep.Equal(value, tt.wantData.SecurityQuotes[i]); diff != nil {
					t.Error(fmt.Sprintf("deep SecurityQuotes[%v]", i), diff)
				}
			}
		})
	}
}
