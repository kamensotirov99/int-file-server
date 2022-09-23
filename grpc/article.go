package grpc

import (
	"context"
	"int-file-server/service"

	pb "int-file-server/_proto"

	"github.com/pkg/errors"
)

type GrpcArticleSvc struct {
	service service.ArticleServicer
	pb.UnimplementedFileServerArticleSvcServer
}

func NewArticleSvc(service service.ArticleServicer) *GrpcArticleSvc {
	return &GrpcArticleSvc{
		service: service,
	}
}

func (c *GrpcArticleSvc) UploadArticlePosters(stream pb.FileServerArticleSvc_UploadArticlePostersServer) error {
	req, err := stream.Recv()
	if err != nil {
		return errors.Wrap(err, "Error while receiving image info")
	}

	articleID := req.GetArticleId()
	imageExtension := req.GetImageExtension()

	err = c.service.UploadArticlePosters(stream.Context(), articleID, imageExtension, stream)
	if err != nil {
		return errors.New("error while uploading article posters")
	}
	return nil
}

func (c *GrpcArticleSvc) DeleteArticlePoster(ctx context.Context, req *pb.DeleteArticlePosterRequest) (*pb.EmptyResponse, error) {
	err := c.service.DeleteArticlePoster(ctx, req.ArticleId, req.Image)
	if err != nil {
		return nil, errors.New("error while deleting article poster")
	}
	return &pb.EmptyResponse{}, nil
}
