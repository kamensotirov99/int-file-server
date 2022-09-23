package grpc

import (
	"context"
	pb "int-file-server/_proto"
	"int-file-server/service"

	"github.com/pkg/errors"
)

type GrpcSeasonSvc struct {
	service service.SeasonServicer
	pb.UnimplementedFileServerSeasonSvcServer
}

func NewSeasonSvc(service service.SeasonServicer) *GrpcSeasonSvc {
	return &GrpcSeasonSvc{
		service: service,
	}
}

func (c *GrpcSeasonSvc) UploadSeasonPosters(stream pb.FileServerSeasonSvc_UploadSeasonPostersServer) error {
	req, err := stream.Recv()
	if err != nil {
		return errors.Wrap(err, "Error while receiving image info")
	}

	seriesID := req.GetSeriesId()
	seasonID := req.GetSeasonId()
	imageExtension := req.GetImageExtension()

	err = c.service.UploadSeasonPosters(stream.Context(), seriesID, seasonID, imageExtension, stream)
	if err != nil {
		return errors.New("error while uploading season posters")
	}
	return nil
}

func (c *GrpcSeasonSvc) DeleteSeasonPoster(ctx context.Context, req *pb.DeleteSeasonPosterRequest) (*pb.EmptyResponse, error) {
	err := c.service.DeleteSeasonPoster(ctx, req.SeriesId, req.SeasonId, req.Image)
	if err != nil {
		return nil, errors.New("error while deleting season poster")
	}
	return &pb.EmptyResponse{}, nil
}
