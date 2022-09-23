package service

import (
	"context"
	"os"

	pb "int-file-server/_proto"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type season struct {
	logger *logrus.Logger
}

type SeasonServicer interface {
	UploadSeasonPosters(ctx context.Context, seriesID string, seasonID string, imageExtension string, stream pb.FileServerSeasonSvc_UploadSeasonPostersServer) error
	DeleteSeasonPoster(ctx context.Context, seriesID string, seasonID string, image string) error
}

func InitiateSeasonService(logger *logrus.Logger) SeasonServicer {
	return &season{
		logger: logger,
	}
}

func (s *season) DeleteSeasonPoster(ctx context.Context, seriesID string, seasonID string, image string) error {
	return os.Remove(baseFilePath + seriesPrefix + seriesID + "/" + seasonID + "/" + image)
}

func (s *season) UploadSeasonPosters(ctx context.Context, seriesID string, seasonID string, imageExtension string, stream pb.FileServerSeasonSvc_UploadSeasonPostersServer) error {
	imageData, err := streamImage(stream, &pb.UploadSeasonPostersRequest{})
	if err != nil {
		s.logger.Error("Error while streaming image")
		return errors.Wrap(err, "Error while streaming image")
	}

	seasonPrefix := seriesPrefix + seriesID + "/" + seasonID + "/"
	postersPath, err := uploadImage(seasonPrefix, imageExtension, imageData)
	if err != nil {
		s.logger.Error("Error while saving image")
		return errors.Wrap(err, "Error while saving image")
	}

	res := &pb.UploadPosterResponse{
		PosterPath: postersPath,
	}

	err = stream.SendAndClose(res)
	if err != nil {
		s.logger.Error("Error while sending response")
		return errors.Wrap(err, "Error while sending response")
	}
	return nil
}
