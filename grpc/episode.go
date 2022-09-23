package grpc

import (
	"context"
	pb "int-file-server/_proto"
	"int-file-server/service"

	"github.com/pkg/errors"
)

type GrpcEpisodeSvc struct {
	service service.EpisodeServicer
	pb.UnimplementedFileServerEpisodeSvcServer
}

func NewEpisodeSvc(service service.EpisodeServicer) *GrpcEpisodeSvc {
	return &GrpcEpisodeSvc{
		service: service,
	}
}

func (c *GrpcEpisodeSvc) UploadEpisodePosters(stream pb.FileServerEpisodeSvc_UploadEpisodePostersServer) error {
	req, err := stream.Recv()
	if err != nil {
		return errors.Wrap(err, "Error while receiving image info")
	}

	seriesID := req.GetSeriesId()
	seasonID := req.GetSeasonId()
	episodeID := req.GetEpisodeId()
	imageExtension := req.GetImageExtension()

	err = c.service.UploadEpisodePosters(stream.Context(), seriesID, seasonID, episodeID, imageExtension, stream)
	if err != nil {
		return errors.New("error while uploading episode posters")
	}
	return nil
}

func (c *GrpcEpisodeSvc) DeleteEpisodePoster(ctx context.Context, req *pb.DeleteEpisodePosterRequest) (*pb.EmptyResponse, error) {
	err := c.service.DeleteEpisodePoster(ctx, req.SeriesId, req.SeasonId, req.EpisodeId, req.Image)
	if err != nil {
		return nil, errors.New("error while deleting episode poster")
	}
	return &pb.EmptyResponse{}, nil
}
