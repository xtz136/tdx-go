package v1

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/go-test/deep"
)

func TestGetMinuteTimeDataResponse_Unmarshal(t *testing.T) {
	type WantData struct {
		Count              int
		MinuteTimeDataList []MinuteTimeData
	}
	type args struct {
		filePath string
		data     []byte
	}
	tests := []struct {
		name     string
		wantData WantData
		args     args
		wantErr  bool
	}{
		{
			name: "basic",
			wantData: WantData{Count: 120, MinuteTimeDataList: []MinuteTimeData{
				{Price: 0, Vol: 48}, {Price: 0.48, Vol: 48}, {Price: 0.97, Vol: 9}, {Price: 13.28, Vol: 11}, {Price: 13.39, Vol: 11298770}, {Price: 1.08, Vol: 27}, {Price: -48.39, Vol: -14}, {Price: 3455.68, Vol: 0}, {Price: 4034.82, Vol: 0}, {Price: 4065.51, Vol: -2}, {Price: 4065.52, Vol: 398}, {Price: 4065.49, Vol: 13932}, {Price: 4080.79, Vol: 11}, {Price: 4093.16, Vol: 44155}, {Price: 4093.14, Vol: 32130}, {Price: 4093.17, Vol: 15253}, {Price: 4093.16, Vol: 19040}, {Price: 4093.15, Vol: 8875}, {Price: 4093.16, Vol: 15099}, {Price: 4093.15, Vol: 6243}, {Price: 4093.16, Vol: 6244}, {Price: 4093.17, Vol: 4068}, {Price: 4093.17, Vol: 4004}, {Price: 4093.16, Vol: 5850}, {Price: 4093.15, Vol: 5347}, {Price: 4093.16, Vol: 4637}, {Price: 4093.16, Vol: 7608}, {Price: 4093.15, Vol: 7256}, {Price: 4093.15, Vol: 2828}, {Price: 4093.15, Vol: 2286}, {Price: 4093.16, Vol: 2797}, {Price: 4093.16, Vol: 3718}, {Price: 4093.15, Vol: 13583}, {Price: 4093.15, Vol: 12771}, {Price: 4093.13, Vol: 6609}, {Price: 4093.13, Vol: 5632}, {Price: 4093.14, Vol: 3076}, {Price: 4093.14, Vol: 6013}, {Price: 4093.14, Vol: 8749}, {Price: 4093.14, Vol: 3273}, {Price: 4093.14, Vol: 3382}, {Price: 4093.15, Vol: 2185}, {Price: 4093.14, Vol: 3631}, {Price: 4093.15, Vol: 3365}, {Price: 4093.15, Vol: 1534}, {Price: 4093.16, Vol: 2946}, {Price: 4093.15, Vol: 1758}, {Price: 4093.17, Vol: 6103}, {Price: 4093.16, Vol: 3349}, {Price: 4093.16, Vol: 4522}, {Price: 4093.16, Vol: 4394}, {Price: 4093.14, Vol: 8551}, {Price: 4093.15, Vol: 4159}, {Price: 4093.15, Vol: 1057}, {Price: 4093.14, Vol: 5449}, {Price: 4093.13, Vol: 7006}, {Price: 4093.14, Vol: 1799}, {Price: 4093.14, Vol: 1834}, {Price: 4093.14, Vol: 2040}, {Price: 4093.15, Vol: 4517}, {Price: 4093.14, Vol: 1897}, {Price: 4093.15, Vol: 3939}, {Price: 4093.14, Vol: 4506}, {Price: 4093.16, Vol: 3663}, {Price: 4093.16, Vol: 2359}, {Price: 4093.16, Vol: 1180}, {Price: 4093.15, Vol: 940}, {Price: 4093.16, Vol: 1018}, {Price: 4093.14, Vol: 4715}, {Price: 4093.14, Vol: 4572}, {Price: 4093.14, Vol: 2434}, {Price: 4093.15, Vol: 9817}, {Price: 4093.13, Vol: 6578}, {Price: 4093.13, Vol: 3171}, {Price: 4093.14, Vol: 3787}, {Price: 4093.14, Vol: 4885}, {Price: 4093.14, Vol: 2364}, {Price: 4093.15, Vol: 1298}, {Price: 4093.15, Vol: 2682}, {Price: 4093.15, Vol: 2010}, {Price: 4093.15, Vol: 2173}, {Price: 4093.14, Vol: 1156}, {Price: 4093.14, Vol: 1430}, {Price: 4093.14, Vol: 1143}, {Price: 4093.14, Vol: 5796}, {Price: 4093.13, Vol: 1922}, {Price: 4093.13, Vol: 4401}, {Price: 4093.13, Vol: 18171}, {Price: 4093.12, Vol: 9305}, {Price: 4093.13, Vol: 4023}, {Price: 4093.12, Vol: 2178}, {Price: 4093.13, Vol: 3430}, {Price: 4093.13, Vol: 2027}, {Price: 4093.13, Vol: 1976}, {Price: 4093.12, Vol: 4864}, {Price: 4093.13, Vol: 3312}, {Price: 4093.12, Vol: 13289}, {Price: 4093.11, Vol: 6677}, {Price: 4093.11, Vol: 3602}, {Price: 4093.12, Vol: 3534}, {Price: 4093.12, Vol: 1871}, {Price: 4093.12, Vol: 3045}, {Price: 4093.11, Vol: 3384}, {Price: 4093.11, Vol: 10037}, {Price: 4093.11, Vol: 3923}, {Price: 4093.1, Vol: 6496}, {Price: 4093.1, Vol: 27039}, {Price: 4093.1, Vol: 4700}, {Price: 4093.11, Vol: 10653}, {Price: 4093.1, Vol: 3694}, {Price: 4093.1, Vol: 1955}, {Price: 4093.1, Vol: 3000}, {Price: 4093.1, Vol: 9479}, {Price: 4093.09, Vol: 2239}, {Price: 4093.09, Vol: 2261}, {Price: 4093.09, Vol: 3814}, {Price: 4093.1, Vol: 1977}, {Price: 4093.09, Vol: 5039}, {Price: 4093.09, Vol: 4064}, {Price: 4093.09, Vol: 3223},
			}},
			args:    args{filePath: "test_datas/get_minute_time_data/default"},
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
			r := &GetMinuteTimeDataResponse{}
			if err := r.Unmarshal(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("GetMinuteTimeDataResponse error = %v, wantErr %v", err, tt.wantErr)
			}
			if r.Count != tt.wantData.Count {
				t.Errorf("GetMinuteTimeDataResponse count error = %v, wantErr %v", r.Count, tt.wantData.Count)
			}
			for i, value := range r.MinuteTimeDataList {
				if diff := deep.Equal(value, tt.wantData.MinuteTimeDataList[i]); diff != nil {
					t.Error(fmt.Sprintf("deep equal[%v]", i), diff)
				}
			}
		})
	}
}
