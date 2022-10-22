package v1

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/go-test/deep"
)

func TestGetTransactionDataResponse_Unmarshal(t *testing.T) {
	type fields struct {
		Count               int
		TransactionDataList []TransactionData
	}
	type args struct {
		filePath string
		data     []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{Count: 30, TransactionDataList: []TransactionData{
				{Time: "14:55", Price: 11.1, Vol: 81, Num: 13, Buyorsell: 0}, {Time: "14:55", Price: 11.09, Vol: 171, Num: 17, Buyorsell: 1}, {Time: "14:55", Price: 11.09, Vol: 101, Num: 14, Buyorsell: 1}, {Time: "14:55", Price: 11.1, Vol: 670, Num: 28, Buyorsell: 0}, {Time: "14:55", Price: 11.09, Vol: 624, Num: 45, Buyorsell: 1}, {Time: "14:55", Price: 11.09, Vol: 132, Num: 9, Buyorsell: 1}, {Time: "14:55", Price: 11.1, Vol: 222, Num: 23, Buyorsell: 0}, {Time: "14:55", Price: 11.09, Vol: 189, Num: 21, Buyorsell: 1}, {Time: "14:56", Price: 11.09, Vol: 180, Num: 25, Buyorsell: 1}, {Time: "14:56", Price: 11.09, Vol: 115, Num: 18, Buyorsell: 1}, {Time: "14:56", Price: 11.09, Vol: 243, Num: 32, Buyorsell: 1}, {Time: "14:56", Price: 11.1, Vol: 366, Num: 63, Buyorsell: 0}, {Time: "14:56", Price: 11.09, Vol: 209, Num: 29, Buyorsell: 1}, {Time: "14:56", Price: 11.09, Vol: 327, Num: 54, Buyorsell: 1}, {Time: "14:56", Price: 11.09, Vol: 372, Num: 63, Buyorsell: 1}, {Time: "14:56", Price: 11.09, Vol: 123, Num: 47, Buyorsell: 1}, {Time: "14:56", Price: 11.09, Vol: 114, Num: 13, Buyorsell: 1}, {Time: "14:56", Price: 11.09, Vol: 155, Num: 13, Buyorsell: 1}, {Time: "14:56", Price: 11.1, Vol: 178, Num: 20, Buyorsell: 0}, {Time: "14:56", Price: 11.09, Vol: 886, Num: 59, Buyorsell: 1}, {Time: "14:56", Price: 11.09, Vol: 564, Num: 49, Buyorsell: 1}, {Time: "14:56", Price: 11.1, Vol: 187, Num: 23, Buyorsell: 0}, {Time: "14:56", Price: 11.1, Vol: 455, Num: 30, Buyorsell: 0}, {Time: "14:56", Price: 11.1, Vol: 310, Num: 32, Buyorsell: 0}, {Time: "14:56", Price: 11.1, Vol: 635, Num: 73, Buyorsell: 0}, {Time: "14:56", Price: 11.09, Vol: 198, Num: 15, Buyorsell: 1}, {Time: "14:56", Price: 11.1, Vol: 486, Num: 28, Buyorsell: 0}, {Time: "14:56", Price: 11.1, Vol: 164, Num: 18, Buyorsell: 0}, {Time: "14:57", Price: 11.1, Vol: 87, Num: 12, Buyorsell: 0}, {Time: "15:00", Price: 11.09, Vol: 8213, Num: 454, Buyorsell: 2},
			}},
			args:    args{filePath: "test_datas/get_transaction_data/default"},
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
			resp := &GetTransactionDataResponse{}
			if err := resp.Unmarshal(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("GetTransactionDataResponse error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if resp.Count != tt.fields.Count {
				t.Errorf("GetSecurityQuotes count error = %v, wantErr %v", resp.Count, tt.fields.Count)
				return
			}
			for i, value := range resp.TransactionDataList {
				if diff := deep.Equal(value, tt.fields.TransactionDataList[i]); diff != nil {
					t.Error(fmt.Sprintf("deep equal[%v]", i), diff)
					return
				}
			}
		})
	}
}
