package rpc

import (
    db "github.com/SwanHtetAungPhyo/lmssystem/internal/repo"
    coursesv "github.com/SwanHtetAungPhyo/lmssystem/app/rpc/course"
    coursepb "github.com/SwanHtetAungPhyo/lmssystem/protogen/course"
    tenantsv "github.com/SwanHtetAungPhyo/lmssystem/app/rpc/tenant"
    tenantpb "github.com/SwanHtetAungPhyo/lmssystem/protogen/tenant"
    "github.com/sirupsen/logrus"
    "google.golang.org/grpc"
    "net"
)

// Server implements the gRPC services
type Server struct {
    coursepb.UnimplementedCourseServiceServer
    tenantpb.UnimplementedTenantServiceServer
    store  *db.Store
    logger *logrus.Logger
}

func NewServer(
    db *db.Store,
    logger *logrus.Logger,
) *Server {
    return &Server{
        store:  db,
        logger: logger,
    }
}

// Run starts the gRPC server
func (s *Server) Run() error {
    listener, err := net.Listen("tcp", ":9001")
    if err != nil {
        panic(err.Error())
    }
    grpcServer := grpc.NewServer()

    courseService := coursesv.NewCourseService()
    tenantService := tenantsv.NewTenantService()

    // Register services with gRPC server
    coursepb.RegisterCourseServiceServer(grpcServer, courseService)
    tenantpb.RegisterTenantServiceServer(grpcServer, tenantService)

    s.logger.Println("Starting server on port 9001")
    err = grpcServer.Serve(listener)
    if err != nil {
        s.logger.Fatal(err.Error())
        return err
    }
    return nil
}