package index_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/brogand93/index-replicate/pkg/index"
)

type HttpMock struct {
	res http.Response
}

func (c *HttpMock) Do(req *http.Request) (*http.Response, error) {
	return &c.res, nil
}

func Test_Get(t *testing.T) {
	type args struct {
		index string
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
		mockRes http.Response
		want    index.Index
	}{
		{
			name: "Simple index request should be successful",
			args: args{
				index: "sample-index",
			},
			wantErr: nil,
			mockRes: http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewReader([]byte(
					`{"companyList":[{"symbol":"SI","name":"sample-index","weight":100,"rank":1,"changePercent":0.0250020835069531,"lastPrice":100,"netChange":0.03,"changeColor":"green"}]}`,
				))),
			},
			want: index.Index{
				Components: []index.Component{
					{
						Name:   "sample-index",
						Price:  100.0,
						Rank:   1,
						Symbol: "SI",
						Weight: 100.0,
						Value:  0,
					},
				},
			},
		}, {
			name: "invalid index returns NotFound error",
			args: args{
				index: "invalid-index",
			},
			wantErr: index.NotFound{
				Message: "invalid-index is not a valid index for this data source",
			},
			mockRes: http.Response{
				StatusCode: 404,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
			},
			want: index.Index{},
		}, {
			name: "index data source error returns DataSourceUnavailable error",
			args: args{
				index: "error-index",
			},
			wantErr: index.DataSourceUnavailable{
				Message: "https://svcga.com/sc/index data source is either unavailable or misconfigured",
			},
			mockRes: http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
			},
			want: index.Index{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := index.Client{
				Http: &HttpMock{
					res: tt.mockRes,
				},
			}
			got, err := mockClient.Get(tt.args.index)

			if err != tt.wantErr {
				t.Errorf("Get() = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, &tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
