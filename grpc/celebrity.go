package grpc

import (
	"int-file-server/service"

	"context"
	pb "int-file-server/_proto"

	"github.com/pkg/errors"
)

type GrpcCelebritySvc struct {
	service service.CelebrityServicer
	pb.UnimplementedFileServerCelebritySvcServer
}

func NewCelebritySvc(service service.CelebrityServicer) *GrpcCelebritySvc {
	return &GrpcCelebritySvc{
		service: service,
	}
}

func (c *GrpcCelebritySvc) UploadCelebrityPosters(stream pb.FileServerCelebritySvc_UploadCelebrityPostersServer) error {
	req, err := stream.Recv()
	if err != nil {
		return errors.Wrap(err, "Error while receiving image info")
	}

	celebrityID := req.GetCelebrityId()
	imageExtension := req.GetImageExtension()

	err = c.service.UploadCelebrityPosters(stream.Context(), celebrityID, imageExtension, stream)
	if err != nil {
		return errors.New("error while uploading celebrity posters")
	}
	return nil
}

func (c *GrpcCelebritySvc) DeleteCelebrityPoster(ctx context.Context, req *pb.DeleteCelebrityPosterRequest) (*pb.EmptyResponse, error) {
	err := c.service.DeleteCelebrityPoster(ctx, req.CelebrityId, req.Image)
	if err != nil {
		return nil, errors.New("error while deleting celebrity poster")
	}
	return &pb.EmptyResponse{}, nil
}
