package grpc

import (
	"context"
	"int-file-server/service"

	pb "int-file-server/_proto"

	"github.com/pkg/errors"
)

type GrpcMovieSvc struct {
	service service.MovieServicer
	pb.UnimplementedFileServerMovieSvcServer
}

func NewMovieSvc(service service.MovieServicer) *GrpcMovieSvc {
	return &GrpcMovieSvc{
		service: service,
	}
}

func (c *GrpcMovieSvc) UploadMoviePosters(stream pb.FileServerMovieSvc_UploadMoviePostersServer) error {
	req, err := stream.Recv()
	if err != nil {
		return errors.Wrap(err, "Error while receiving image info")
	}

	movieID := req.GetMovieId()
	imageExtension := req.GetImageExtension()

	err = c.service.UploadMoviePosters(stream.Context(), movieID, imageExtension, stream)
	if err != nil {
		return errors.New("error while uploading movie posters")
	}
	return nil
}

func (c *GrpcMovieSvc) DeleteMoviePoster(ctx context.Context, req *pb.DeleteMoviePosterRequest) (*pb.EmptyResponse, error) {
	err := c.service.DeleteMoviePoster(ctx, req.MovieId, req.Image)
	if err != nil {
		return nil, errors.New("error while deleting movie poster")
	}
	return &pb.EmptyResponse{}, nil
}
