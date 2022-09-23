package service

import (
	"context"
	pb "int-file-server/_proto"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type article struct {
	logger *logrus.Logger
}

type ArticleServicer interface {
	UploadArticlePosters(ctx context.Context, articleID string, imageExtension string, stream pb.FileServerArticleSvc_UploadArticlePostersServer) error
	DeleteArticlePoster(ctx context.Context, articleID string, image string) error
}

func InitiateArticleService(logger *logrus.Logger) ArticleServicer {
	return &article{
		logger: logger,
	}
}

func (a *article) DeleteArticlePoster(ctx context.Context, articleID string, image string) error {
	return os.Remove(baseFilePath + articlePrefix + articleID + "/" + image)
}

func (a *article) UploadArticlePosters(ctx context.Context, articleID string, imageExtension string, stream pb.FileServerArticleSvc_UploadArticlePostersServer) error {
	imageData, err := streamImage(stream, &pb.UploadSeriesPostersRequest{})
	if err != nil {
		a.logger.Error("Error while streaming image")
		return errors.Wrap(err, "Error while streaming image")
	}

	articlePrefix := articlePrefix + articleID + "/"
	postersPath, err := uploadImage(articlePrefix, imageExtension, imageData)
	if err != nil {
		a.logger.Error("Error while saving image")
		return errors.Wrap(err, "Error while saving image")
	}

	res := &pb.UploadPosterResponse{
		PosterPath: postersPath,
	}

	err = stream.SendAndClose(res)
	if err != nil {
		a.logger.Error("Error while sending response")
		return errors.Wrap(err, "Error while sending response")
	}
	return nil
}
