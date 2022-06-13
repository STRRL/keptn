package sdk

import (
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cp-common/api"
	"net/http"
	"reflect"
	"testing"
)

func Test_createAPI(t *testing.T) {

	apiInit := Initializer{
		Internal: func(client *http.Client, apiMappings ...api.InClusterAPIMappings) (*api.InternalAPISet, error) {
			return &api.InternalAPISet{}, nil
		},
		Remote: func(baseURL string, options ...func(*keptnapi.APISet)) (*keptnapi.APISet, error) {
			return &keptnapi.APISet{}, nil
		},
	}

	tests := []struct {
		name         string
		env          envConfig
		wantInternal bool
		wantErr      bool
	}{
		{
			name:         "test no env internal NATS ",
			env:          envConfig{},
			wantInternal: true,
			wantErr:      false,
		},
		{
			name: "test FAIL for no http address",
			env: envConfig{
				KeptnAPIEndpoint: "ssh://mynotsogoodendpoint",
			},
			wantErr:      true,
			wantInternal: false,
		},
		{
			name: "test FAIL for no good address",
			env: envConfig{
				KeptnAPIEndpoint: ":///MALFORMEDendpoint",
			},
			wantErr:      true,
			wantInternal: false,
		},
		{
			name: "test PASS for http address",
			env: envConfig{
				KeptnAPIEndpoint: "http://endpoint",
			},
			wantErr:      false,
			wantInternal: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createAPI(nil, tt.env, apiInit)
			if (err != nil) != tt.wantErr {
				t.Errorf("createAPI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.wantInternal && !reflect.DeepEqual(got, &api.InternalAPISet{}) {
				t.Errorf("createAPI() got = %v, wanted internal API", got)
			} else if err == nil && !tt.wantInternal && !reflect.DeepEqual(got, &keptnapi.APISet{}) {
				t.Errorf("createAPI() got = %v, want remote execution plane", got)
			}

		})
	}
}