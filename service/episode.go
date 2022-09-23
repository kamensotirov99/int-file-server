package service

import (
	"context"
	"os"

	pb "int-file-server/_proto"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type episode struct {
	logger *logrus.Logger
}

type EpisodeServicer interface {
	UploadEpisodePosters(ctx context.Context, seriesID string, seasonID string, episodeID string, imageExtension string, stream pb.FileServerEpisodeSvc_UploadEpisodePostersServer) error
	DeleteEpisodePoster(ctx context.Context, seriesID string, seasonID string, episodeID string, image string) error
}

func InitiateEpisodeService(logger *logrus.Logger) EpisodeServicer {
	return &episode{
		logger: logger,
	}
}

func (e *episode) DeleteEpisodePoster(ctx context.Context, seriesID string, seasonID string, episodeID string, image string) error {
	return os.Remove(baseFilePath + seriesPrefix + seriesID + "/" + seasonID + "/" + episodeID + "/" + image)
}

func (e *episode) UploadEpisodePosters(ctx context.Context, seriesID string, seasonID string, episodeID string, imageExtension string, stream pb.FileServerEpisodeSvc_UploadEpisodePostersServer) error {
	imageData, err := streamImage(stream, &pb.UploadEpisodePostersRequest{})
	if err != nil {
		e.logger.Error("Error while streaming image")
		return errors.Wrap(err, "Error while streaming image")
	}

	episodePrefix := seriesPrefix + seriesID + "/" + seasonID + "/" + episodeID + "/"
	postersPath, err := uploadImage(episodePrefix, imageExtension, imageData)
	if err != nil {
		e.logger.Error("Error while saving image")
		return errors.Wrap(err, "Error while saving image")
	}

	res := &pb.UploadPosterResponse{
		PosterPath: postersPath,
	}

	err = stream.SendAndClose(res)
	if err != nil {
		e.logger.Error("Error while sending response")
		return errors.Wrap(err, "Error while sending response")
	}
	return nil
}