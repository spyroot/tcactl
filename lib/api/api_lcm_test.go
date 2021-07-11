package api

import (
	"context"
	"github.com/spyroot/tcactl/lib/client"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/client/specs"
	"testing"
)

func TestGetInstance(t *testing.T) {
	tests := []struct {
		name         string
		rest         *client.RestClient
		wantErr      bool
		instanceName string
		poolName     string
		clusterName  string
		vimName      string
	}{
		{
			name:         "Must return test instance",
			rest:         rest,
			wantErr:      false,
			instanceName: "unit_test_instance",
			poolName:     getTestNodePoolName(),
			vimName:      getTestCloudProvider(),
			clusterName:  getTestWorkloadClusterName(),
		},
		{
			name:         "Must return error",
			rest:         rest,
			wantErr:      true,
			instanceName: "unit_test_instance01",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				instance *response.LcmInfo
				err      error
			)

			a := getTcaApi(t, tt.rest, false)
			ctx := context.Background()
			if instance, err = a.GetInstance(ctx, tt.instanceName); (err != nil) != tt.wantErr {
				t.Errorf("CreateCnfInstance() error = %v, wantOnGetErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err == nil {
				t.Errorf("CreateCnfInstance() err must not be nil")
				return
			}
			if tt.wantErr && err != nil {
				return
			}
			if instance == nil {
				t.Errorf("CreateCnfInstance() no error instance must not be nil")
				return
			}
			if instance.VnfInstanceName != tt.instanceName {
				t.Errorf("CreateCnfInstance() return wrong instance")
				return
			}

		})
	}
}

func TestCreateCnfInstance(t *testing.T) {

	tests := []struct {
		name            string
		rest            *client.RestClient
		wantOnGetErr    bool
		wantOnCreateErr bool
		instanceName    string
		poolName        string
		clusterName     string
		vimName         string
	}{
		{
			name:            "Create unit_test instance must fail",
			rest:            rest,
			wantOnGetErr:    false,
			wantOnCreateErr: true,
			instanceName:    "unit_test_instance",
			poolName:        "wrong pool name",
			vimName:         getTestCloudProvider(),
			clusterName:     getTestWorkloadClusterName(),
		},
		{
			name:            "Create unit_test instance",
			rest:            rest,
			wantOnGetErr:    false,
			wantOnCreateErr: true,
			instanceName:    "unit_test_instance",
			poolName:        getTestNodePoolName(),
			vimName:         "wrong cluster",
			clusterName:     getTestWorkloadClusterName(),
		},
		{
			name:            "Create unit_test instance must pass",
			rest:            rest,
			wantOnGetErr:    false,
			wantOnCreateErr: false,
			instanceName:    "unit_test_instance",
			poolName:        getTestNodePoolName(),
			vimName:         getTestCloudProvider(),
			clusterName:     getTestWorkloadClusterName(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a := getTcaApi(t, tt.rest, false)
			ctx := context.Background()

			instance, err := a.GetInstance(ctx, tt.instanceName)
			if tt.wantOnGetErr && err == nil {
				t.Errorf("TestCreateCnfInstance() must return error")
				return
			}
			if tt.wantOnGetErr && err != nil {
				return
			}

			if instance == nil {
				t.Errorf("TestCreateCnfInstance() instance must not be nil")
				return
			}

			if instance.IsFailed() {
				err := a.RollbackCnf(ctx, tt.instanceName, true, true)
				if err != nil {
					t.Errorf("TestCreateCnfInstance() rollback failed %v", err)
					return
				}
			}

			if err := a.CreateCnfInstance(ctx, &CreateInstanceApiReq{
				InstanceName: tt.instanceName,
				PoolName:     tt.poolName,
				VimName:      tt.vimName,
				ClusterName:  tt.clusterName,
				Namespace:    "default",
				RepoUsername: getTestRepoUsername(),
				RepoPassword: getTestRepoPassword(),
				RepoUrl:      getTestRepoUrl(),
				IsBlocking:   false,
				IsVerbose:    false,
				AdditionalParam: &specs.AdditionalParams{
					DisableGrant:        true,
					IgnoreGrantFailure:  false,
					DisableAutoRollback: false,
				},
			}); (err != nil) != tt.wantOnCreateErr {
				t.Errorf("CreateCnfInstance() error = %v, wantOnGetErr %v", err, tt.wantOnCreateErr)
			}
		})
	}
}

func TestCnfUpdateState(t *testing.T) {

	tests := []struct {
		name         string
		rest         *client.RestClient
		wantErr      bool
		instanceName string
		poolName     string
		clusterName  string
		vimName      string
		doBlock      bool
		doVerbose    bool
	}{
		{
			name:         "Update unit_test instance",
			rest:         rest,
			instanceName: "unit_test_instance",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a := getTcaApi(t, tt.rest, false)
			ctx := context.Background()

			instance, err := a.GetInstance(ctx, tt.instanceName)
			if tt.wantErr && err == nil {
				t.Errorf("TestCnfUpdateState must not be nil")
				return
			}
			if tt.wantErr && err != nil {
				return
			}
			if instance == nil {
				t.Errorf("GetInstance() must not be nil")
				return
			}

			actions, err := a.GetLcmActions(ctx, tt.instanceName)
			if err != nil {
				return
			}

			// we can only reset
			if len(actions.Instantiate.Href) == 0 && len(actions.UpdateState.Href) > 0 {
				err := a.ResetState(ctx, &ResetInstanceApiReq{
					InstanceName: tt.instanceName,
					IsBlocking:   tt.doBlock,
					IsVerbose:    tt.doVerbose,
				})
				if err != nil {
					return
				}
			}

			if len(actions.Instantiate.Href) > 0 {
				instantiateReq := specs.LcmInstantiateRequest{
					FlavourID: "default",
					AdditionalVduParams: &specs.AdditionalParams{
						DisableGrant:        true,
						IgnoreGrantFailure:  false,
						DisableAutoRollback: false,
					},
					VimConnectionInfo: instance.VimConnectionInfo,
				}

				if _, err := a.UpdateCnfState(ctx, &UpdateInstanceApiReq{
					InstanceName: tt.instanceName,
					ClusterName:  tt.clusterName,
					IsBlocking:   tt.doBlock,
					IsVerbose:    tt.doVerbose,
					UpdateReq:    &instantiateReq,
				}); (err != nil) != tt.wantErr {
					t.Errorf("CreateCnfInstance() error = %v, wantOnGetErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestDeleteCnf(t *testing.T) {

	tests := []struct {
		name         string
		rest         *client.RestClient
		wantErr      bool
		instanceName string
	}{
		{
			name:         "Instance not found",
			rest:         rest,
			wantErr:      true,
			instanceName: "invalid name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			a := getTcaApi(t, tt.rest, false)
			if err := a.DeleteCnf(ctx, tt.instanceName); (err != nil) != tt.wantErr {
				t.Errorf("DeleteCnf() error = %v, wantOnGetErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteCnfInstance(t *testing.T) {

	tests := []struct {
		name         string
		rest         *client.RestClient
		wantErr      bool
		instanceName string
		vimName      string
		isForce      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			a := getTcaApi(t, tt.rest, false)

			if err := a.DeleteCnfInstance(ctx, tt.instanceName, tt.vimName, tt.isForce); (err != nil) != tt.wantErr {
				t.Errorf("DeleteCnfInstance() error = %v, wantOnGetErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRollbackCnf(t *testing.T) {

	tests := []struct {
		name         string
		rest         *client.RestClient
		instanceName string
		doBlock      bool
		doVerbose    bool
		wantErr      bool
	}{
		{
			name:         "rollback cnf",
			rest:         rest,
			instanceName: getTestInstanceName(),
			doBlock:      true,
			doVerbose:    true,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			a := getTcaApi(t, tt.rest, false)
			if err := a.RollbackCnf(ctx, tt.instanceName, tt.doBlock, tt.doVerbose); (err != nil) != tt.wantErr {
				t.Errorf("RollbackCnf() error = %v, wantOnGetErr %v", err, tt.wantErr)
			}
		})
	}
}

//TestTerminateCnfInstance
func TestTerminateCnfInstance(t *testing.T) {
	tests := []struct {
		name         string
		rest         *client.RestClient
		wantErr      bool
		instanceName string
		clusterName  string
		doBlock      bool
		doVerbose    bool
	}{
		{
			name:         "Instance not found",
			rest:         rest,
			wantErr:      true,
			instanceName: "invalid name",
		},
		{
			name:         "Instance must terminate",
			rest:         rest,
			wantErr:      false,
			instanceName: getTestInstanceName(),
			clusterName:  getTestWorkloadClusterName(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()
			a := getTcaApi(t, tt.rest, false)
			//req := TerminateInstanceApiReq{
			//	InstanceName: tt.InstanceName,
			//	ClusterName:  tt.ClusterName,
			//	IsBlocking:   tt.doBlock,
			//	IsVerbose:    tt.doVerbose,
			//}

			if !tt.wantErr {
				instance, err := a.GetInstance(ctx, tt.instanceName)
				if tt.wantErr && err == nil {
					t.Errorf("TestTerminateCnfInstance() must return error")
					return
				}

				if instance == nil {
					t.Errorf("TestTerminateCnfInstance() failed to get instance")
					return
				}
				if instance.IsFailed() {
					t.Log("Instance failed state")
					state, err := a.UpdateCnfState(ctx, &UpdateInstanceApiReq{
						InstanceName: tt.instanceName,
						ClusterName:  tt.clusterName,
						IsBlocking:   true,
						IsVerbose:    true,
						UpdateReq:    nil,
					})
					if err != nil {
						t.Errorf("TestTerminateCnfInstance() failed update state %s", err)
						return
					}
					if state == nil {
						t.Errorf("TestTerminateCnfInstance() update state must not return nil")
						return
					}
				}
			}

			//if err := a.TerminateCnfInstance(ctx, &req); (err != nil) != tt.wantOnGetErr {
			//	t.Errorf("TestTerminateCnfInstance() error = %v, wantOnGetErr %v", err, tt.wantOnGetErr)
			//}
		})
	}
}
