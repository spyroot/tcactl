package client

import (
	"os"
	"testing"
)

var (
	rest = &RestClient{
		BaseURL:  os.Getenv("TCA_URL"),
		apiKey:   "",
		IsDebug:  true,
		Username: os.Getenv("TCA_USERNAME"),
		Password: os.Getenv("TCA_PASSWORD"),
		SkipSsl:  true,
		isTrace:  true,
	}

	harbor = &RestClient{
		BaseURL:               os.Getenv("HARBOR_URL"),
		apiKey:                "",
		IsDebug:               true,
		Username:              os.Getenv("HARBOR_USERNAME"),
		Password:              os.Getenv("HARBOR_PASSWORD"),
		SkipSsl:               true,
		isTrace:               true,
		isBasicAuthentication: true,
	}
)

func TestRestClient_GetInfrNetworks(t *testing.T) {

	type args struct {
		clusterId string
	}

	tests := []struct {
		name   string
		client *RestClient
		//		args    args
		//		want    *models.CloudNetworks
		wantErr bool
	}{
		{
			name:    "base",
			client:  rest,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := tt.client.GetAuthorization()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInfraNetworks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			tenants, err := tt.client.GetVimTenants()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInfraNetworks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, details := range tenants.TenantsList {
				_, err = tt.client.GetInfraNetworks(details.TenantID)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetInfraNetworks() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("GetInfraNetworks() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
