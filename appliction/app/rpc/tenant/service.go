package tenant

import (
	tenantpb "github.com/SwanHtetAungPhyo/mmmmm/protogen/tenant"
)

type TenantService struct {
	tenantpb.UnimplementedTenantServiceServer
}

func NewTenantService() *TenantService {
	return &TenantService{}
}
