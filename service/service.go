package service

import (
	"io"
	"os"

	"bytes"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

const (
	baseFilePath    = "../int-static-files"
	celebrityPrefix = "/celebrities/"
	moviePrefix     = "/movies/"
	articlePrefix   = "/articles/"
	seriesPrefix    = "/series/"
	maxImageSize    = 1 << 20
)

func uploadImage(prefix string, imageType string, image bytes.Buffer) (string, error) {
	if imageType != ".jpg" && imageType != ".png" {
		return "", errors.New("The provided file format is not allowed. Please upload a JPG or PNG image.")
	}

	if _, err := os.Stat(baseFilePath + prefix); os.IsNotExist(err) {
		err := os.MkdirAll(baseFilePath+prefix, os.ModePerm)
		if err != nil {
			return "", errors.Wrap(err, "Error while creating folder!")
		}
	}

	fileName := uuid.NewString() + imageType
	out, err := os.OpenFile(baseFilePath+prefix+"/"+fileName, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return "", errors.Wrap(err, "Error while creating image!")
	}
	defer out.Close()

	buff := make([]byte, 1024)
	for {
		_, err := image.Read(buff)
		if err != nil {
			break
		} else {
			out.Write(buff)
		}
	}

	if err != nil {
		return "", errors.Wrap(err, "Error while copying image from source to destination!")
	}
	return prefix + fileName, nil
}

func streamImage(stream grpc.ServerStream, message proto.Message) (bytes.Buffer, error) {
	imageData := bytes.Buffer{}
	imageSize := 0

	for {
		ctx := stream.Context()
		if ctx.Err() != nil {
			return bytes.Buffer{}, errors.Wrap(ctx.Err(), "Error while streaming context")
		}

		err := stream.RecvMsg(message)
		if err == io.EOF {
			break
		}
		if err != nil {
			return bytes.Buffer{}, errors.Wrap(err, "Error while receiving chunk data")
		}

		chunkDataField := message.ProtoReflect().Descriptor().Fields().ByTextName("chunkData")
		if chunkDataField == nil {
			return bytes.Buffer{}, errors.Wrap(err, "Error while getting 'chunk data' field descriptor")
		}
		chunk := message.ProtoReflect().Get(chunkDataField).Bytes()
		size := len(chunk)

		imageSize += size
		if imageSize > maxImageSize {
			return bytes.Buffer{}, errors.Wrap(err, "Image is too large")
		}

		_, err = imageData.Write(chunk)
		if err != nil {
			return bytes.Buffer{}, errors.Wrap(err, "Error while writing chunk data")
		}
	}
	return imageData, nil
}
