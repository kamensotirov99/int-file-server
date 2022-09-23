package service

import (
	"context"
	pb "int-file-server/_proto"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type celebrity struct {
	logger *logrus.Logger
}

type CelebrityServicer interface {
	UploadCelebrityPosters(ctx context.Context, celebrityID string, imageExtension string, stream pb.FileServerCelebritySvc_UploadCelebrityPostersServer) error
	DeleteCelebrityPoster(ctx context.Context, celebrityID string, image string) error
}

func InitiateCelebrityService(logger *logrus.Logger) CelebrityServicer {
	return &celebrity{
		logger: logger,
	}
}

func (c *celebrity) DeleteCelebrityPoster(ctx context.Context, celebrityID string, image string) error {
	return os.Remove(baseFilePath + celebrityPrefix + celebrityID + "/" + image)
}

func (c *celebrity) UploadCelebrityPosters(ctx context.Context, celebrityID string, imageExtension string, stream pb.FileServerCelebritySvc_UploadCelebrityPostersServer) error {
	imageData, err := streamImage(stream, &pb.UploadCelebrityPostersRequest{})
	if err != nil {
		c.logger.Error("Error while streaming image")
		return errors.Wrap(err, "Error while streaming image")
	}

	celebrityPrefix := celebrityPrefix + celebrityID + "/"
	postersPath, err := uploadImage(celebrityPrefix, imageExtension, imageData)
	if err != nil {
		c.logger.Error("Error while saving image")
		return errors.Wrap(err, "Error while saving image")
	}

	res := &pb.UploadPosterResponse{
		PosterPath: postersPath,
	}

	err = stream.SendAndClose(res)
	if err != nil {
		c.logger.Error("Error while sending response")
		return errors.Wrap(err, "Error while sending response")
	}
	return nil
}