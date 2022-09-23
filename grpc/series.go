package grpc

import (
	"context"
	"int-file-server/service"

	pb "int-file-server/_proto"

	"github.com/pkg/errors"
)

type GrpcSeriesSvc struct {
	service service.SeriesServicer
	pb.UnimplementedFileServerSeriesSvcServer
}

func NewSeriesSvc(service service.SeriesServicer) *GrpcSeriesSvc {
	return &GrpcSeriesSvc{
		service: service,
	}
}

func (c *GrpcSeriesSvc) UploadSeriesPosters(stream pb.FileServerSeriesSvc_UploadSeriesPostersServer) error {
	req, err := stream.Recv()
	if err != nil {
		return errors.Wrap(err, "Error while receiving image info")
	}

	seriesID := req.GetSeriesId()
	imageExtension := req.GetImageExtension()

	err = c.service.UploadSeriesPosters(stream.Context(), seriesID, imageExtension, stream)
	if err != nil {
		return errors.New("error while uploading series posters")
	}
	return nil
}

func (c *GrpcSeriesSvc) DeleteSeriesPoster(ctx context.Context, req *pb.DeleteSeriesPosterRequest) (*pb.EmptyResponse, error) {
	err := c.service.DeleteSeriesPoster(ctx, req.SeriesId, req.Image)
	if err != nil {
		return nil, errors.New("error while deleting series poster")
	}
	return &pb.EmptyResponse{}, nil
}
