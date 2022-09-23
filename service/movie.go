package service

import (
	"context"
	pb "int-file-server/_proto"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type movie struct {
	logger *logrus.Logger
}

type MovieServicer interface {
	UploadMoviePosters(ctx context.Context, movieID string, imageExtension string, stream pb.FileServerMovieSvc_UploadMoviePostersServer) error
	DeleteMoviePoster(ctx context.Context, movieID string, image string) error
}

func InitiateMovieService(logger *logrus.Logger) MovieServicer {
	return &movie{
		logger: logger,
	}
}

func (m *movie) DeleteMoviePoster(ctx context.Context, movieID string, image string) error {
	return os.Remove(baseFilePath + moviePrefix + movieID + "/" + image)
}

func (m *movie) UploadMoviePosters(ctx context.Context, movieID string, imageExtension string, stream pb.FileServerMovieSvc_UploadMoviePostersServer) error {
	imageData, err := streamImage(stream, &pb.UploadSeriesPostersRequest{})
	if err != nil {
		m.logger.Error("Error while streaming image")
		return errors.Wrap(err, "Error while streaming image")
	}

	moviePrefix := moviePrefix + movieID + "/"
	postersPath, err := uploadImage(moviePrefix, imageExtension, imageData)
	if err != nil {
		m.logger.Error("Error while saving image")
		return errors.Wrap(err, "Error while saving image")
	}

	res := &pb.UploadPosterResponse{
		PosterPath: postersPath,
	}

	err = stream.SendAndClose(res)
	if err != nil {
		m.logger.Error("Error while sending response")
		return errors.Wrap(err, "Error while sending response")
	}
	return nil
}
