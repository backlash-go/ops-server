package service

import (
	"fmt"
	"ops-server/entity"
	"reflect"
	"testing"
)

func TestQueryPermissionList(t *testing.T) {
	type args struct {
		req *entity.PermissionInfoListRequest
	}
	tests := []struct {
		name               string
		args               args
		wantPermissionInfo []entity.PermissionInfoList
		wantTotalCount     int64
		wantErr            bool
	}{
		{
			name:    "TestQueryPermissionList",
			args: args{req:&entity.PermissionInfoListRequest{PageSize:5,Page:2}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPermissionInfo, gotTotalCount, err := QueryPermissionList(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("QuerryPermissonList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPermissionInfo, tt.wantPermissionInfo) {
				t.Errorf("QuerryPermissonList() gotPermissionInfo = %v, want %v", gotPermissionInfo, tt.wantPermissionInfo)
			}
			if gotTotalCount != tt.wantTotalCount {
				t.Errorf("QuerryPermissonList() gotTotalCount = %v, want %v", gotTotalCount, tt.wantTotalCount)
			}

			fmt.Printf("gotTotalCount is %+v", gotTotalCount)

			fmt.Printf("gotPermissionInfo is %+v", gotPermissionInfo)

		})
	}
}

func TestQueryPermissionListAndRoles(t *testing.T) {
	type args struct {
		permissionId []uint64
	}
	tests := []struct {
		name               string
		args               args
		wantPermissionRole []entity.PermissionIdRoleContact
		wantErr            bool
	}{
		{
			name:    "TestQueryPermissionListAndRoles",
			args: args{permissionId:[]uint64{1,2,3,4,5}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPermissionRole, err := QueryPermissionListAndRoles(tt.args.permissionId)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryPermissionListAndRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPermissionRole, tt.wantPermissionRole) {
				t.Errorf("QueryPermissionListAndRoles() gotPermissionRole = %v, want %v", gotPermissionRole, tt.wantPermissionRole)
			}

			fmt.Printf("gotPermissionRole is :  %+v\n", gotPermissionRole)
		})
	}
}