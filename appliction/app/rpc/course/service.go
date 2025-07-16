package course

import (
	coursepb "github.com/SwanHtetAungPhyo/mmmmm/protogen/course"
)

type CourseService struct {
	coursepb.UnimplementedCourseServiceServer
}

func NewCourseService() *CourseService {
	return &CourseService{}
}
