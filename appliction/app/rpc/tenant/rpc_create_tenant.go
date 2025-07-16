package tenant

import (
	"context"
	tenantpb "github.com/SwanHtetAungPhyo/mmmmm/protogen/tenant"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *TenantService) CreateTenant(ctx context.Context, req *tenantpb.CreateTenantRequest) (*tenantpb.CreateTenantResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	resp := &tenantpb.CreateTenantResponse{}
	return resp, nil
}
