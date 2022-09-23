package service

import (
	"context"
	pb "int-file-server/_proto"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type series struct {
	logger *logrus.Logger
}

type SeriesServicer interface {
	UploadSeriesPosters(ctx context.Context, seriesID string, imageExtension string, stream pb.FileServerSeriesSvc_UploadSeriesPostersServer) error
	DeleteSeriesPoster(ctx context.Context, seriesID string, image string) error
}

func InitiateSeriesService(logger *logrus.Logger) SeriesServicer {
	return &series{
		logger: logger,
	}
}

func (s *series) DeleteSeriesPoster(ctx context.Context, seriesID string, image string) error {
	return os.Remove(baseFilePath + seriesPrefix + seriesID + "/" + image)
}

func (s *series) UploadSeriesPosters(ctx context.Context, seriesID string, imageExtension string, stream pb.FileServerSeriesSvc_UploadSeriesPostersServer) error {
	imageData, err := streamImage(stream, &pb.UploadSeriesPostersRequest{})
	if err != nil {
		s.logger.Error("Error while streaming image")
		return errors.Wrap(err, "Error while streaming image")
	}

	seriesPrefix := seriesPrefix + seriesID + "/"
	postersPath, err := uploadImage(seriesPrefix, imageExtension, imageData)
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
