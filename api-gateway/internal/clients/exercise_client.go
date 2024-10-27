package clients

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	exercisepb "healthApi/api-gateway/genproto/exercises"
	"healthApi/api-gateway/internal/models"
	"healthApi/api-gateway/internal/utils"
)

type ExerciseClient struct {
	client exercisepb.ExerciseServiceClient
	conn   *grpc.ClientConn
}

func NewExerciseClient(serverAddr string) (ExerciseClient, error) {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return ExerciseClient{}, err
	}

	return ExerciseClient{
		client: exercisepb.NewExerciseServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *ExerciseClient) Close() error {
	return c.conn.Close()
}

func (c *ExerciseClient) CreateExercise(ctx context.Context, userId uint, req models.ExerciseRequest) (*exercisepb.Exercise, error) {
	resp, err := c.client.CreateExercise(ctx, &exercisepb.CreateExerciseRequest{
		UserId:      int32(userId),
		Date:        utils.TimeToProtoTimestamp(req.Date),
		Description: req.Description,
		Duration:    int32(req.Duration),
		Intensity:   req.Intensity,
		Type:        req.Type,
	})
	if err != nil {
		return nil, utils.NewInternalServerError("Failed to create exercise : " + err.Error())
	}
	return resp.Exercise, nil
}

func (c *ExerciseClient) GetAllExercise(ctx context.Context, userId uint) ([]*exercisepb.Exercise, error) {
	resp, err := c.client.GetAllExercise(ctx, &exercisepb.GetAllExerciseRequest{
		UserId: int32(userId),
	})
	if err != nil {
		return nil, utils.NewInternalServerError("Failed to fetch all exercises : " + err.Error())
	}
	return resp.Exercise, nil
}
