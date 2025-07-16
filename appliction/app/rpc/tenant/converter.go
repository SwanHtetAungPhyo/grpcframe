package tenant

import (
	"github.com/SwanHtetAungPhyo/mmmmm/internal/repo"
	tenantpb "github.com/SwanHtetAungPhyo/mmmmm/protogen/tenant"
)

// ProtoToModel converts protobuf message to database model
func ProtoToModel(pb *tenantpb.Tenant) *repo.Tenant {
	if pb == nil {
		return nil
	}

	return &repo.Tenant{
		// Add your conversion fields here
		// Example:
		// ID:        pb.GetId(),
		// Name:      pb.GetName(),
		// CreatedAt: time.Unix(pb.GetCreatedAt().GetSeconds(), int64(pb.GetCreatedAt().GetNanos())),
	}
}

// ModelToProto converts database model to protobuf message
func ModelToProto(model *repo.Tenant) *tenantpb.Tenant {
	if model == nil {
		return nil
	}

	return &tenantpb.Tenant{
		// Add your conversion fields here
		// Example:
		// Id:        model.ID,
		// Name:      model.Name,
		// CreatedAt: timestamppb.New(model.CreatedAt),
	}
}

// ModelsToProtos converts slice of models to slice of protos
func ModelsToProtos(models []*repo.Tenant) []*tenantpb.Tenant {
	protos := make([]*tenantpb.Tenant, len(models))
	for i, model := range models {
		protos[i] = ModelToProto(model)
	}
	return protos
}

// ProtosToModels converts slice of protos to slice of models
func ProtosToModels(protos []*tenantpb.Tenant) []*repo.Tenant {
	models := make([]*repo.Tenant, len(protos))
	for i, pb := range protos {
		models[i] = ProtoToModel(pb)
	}
	return models
}
