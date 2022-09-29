package client

import (
	"context"
	"testing"
)

func TestRestClient_GetConsumption(t *testing.T) {

	tests := []struct {
		name    string
		client  *RestClient
		wantErr bool
	}{
		{
			name:    "basic test",
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

			consumption, err := tt.client.GetConsumption(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConsumption() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				return
			}
			if consumption == nil {
				t.Errorf("GetConsumption() shouldn't return nil.")
				return
			}
		})
	}
}
