package course

import (
	"context"
	coursepb "github.com/SwanHtetAungPhyo/mmmmm/protogen/course"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CourseService) CreateCourse(ctx context.Context, req *coursepb.CreateCourseRequest) (*coursepb.CreateCourseResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	resp := &coursepb.CreateCourseResponse{}
	return resp, nil
}
