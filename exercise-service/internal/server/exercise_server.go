package server

import (
	"context"
	exercisepb "exercise-service/internal/genproto/exercises"
	"exercise-service/internal/model"
	"exercise-service/internal/service"
	"exercise-service/internal/utils"
)

type ExerciseServer struct {
	exercisepb.UnimplementedExerciseServiceServer
	service service.ExerciseService
}

func NewExerciseServer(svc service.ExerciseService) *ExerciseServer {
	return &ExerciseServer{service: svc}
}

func (s *ExerciseServer) CreateExercise(ctx context.Context, req *exercisepb.CreateExerciseRequest) (*exercisepb.CreateExerciseResponse, error) {

	err := s.service.CreateExercise(&model.Exercise{
		UserID:      uint(req.UserId),
		Type:        req.Type,
		Duration:    int(req.Duration),
		Intensity:   req.Intensity,
		Date:        utils.ProtoTimestampToTime(req.Date),
		Description: req.Intensity,
	})
	if err != nil {
		return nil, utils.NewInternalServerError("Failed to create exercise : " + err.Error())
	}

	return &exercisepb.CreateExerciseResponse{
		Exercise: &exercisepb.Exercise{
			Id:          req.UserId,
			Date:        req.Date,
			Description: req.Description,
			Duration:    req.Duration,
			Intensity:   req.Intensity,
			Type:        req.Type,
		},
	}, nil
}

func (s *ExerciseServer) GetAllExercise(ctx context.Context, req *exercisepb.GetAllExerciseRequest) (*exercisepb.GetAllExerciseResponse, error) {
	exercises, err := s.service.GetUserExercises(uint(req.UserId))
	if err != nil {
		return nil, utils.NewInternalServerError("Failed to get all exercises : " + err.Error())
	}
	var response []*exercisepb.Exercise
	for _, exercise := range exercises {
		response = append(response, &exercisepb.Exercise{
			Id:          int32(exercise.ID),
			Date:        utils.TimeToProtoTimestamp(exercise.Date),
			Description: exercise.Description,
			Duration:    int32(exercise.Duration),
			Intensity:   exercise.Intensity,
			Type:        exercise.Type,
		})
	}
	return &exercisepb.GetAllExerciseResponse{Exercise: response}, nil
}
