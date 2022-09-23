package app

import (
	transport_grpc "int-file-server/grpc"
	"int-file-server/service"
	"net"
	"net/http"

	pb "int-file-server/_proto"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type app struct {
	router *mux.Router
	logger *logrus.Logger
}

const (
	HttpPort = "2003"
	GrpcPort = "2004"
)

func InitializeHttpApp(logger *logrus.Logger) *app {
	a := app{}
	a.router = mux.NewRouter()

	a.logger = logger
	a.registerRoutes()
	return &a
}

func InitializeGrpcApp(logger *logrus.Logger) {
	a := app{}
	a.logger = logger
	a.createGrpcServer(GrpcPort)
}

func (a *app) createGrpcServer(serverPort string) {
	celebrityService := service.InitiateCelebrityService(a.logger)
	grpcCelebrityServer := transport_grpc.NewCelebritySvc(celebrityService)

	articleService := service.InitiateArticleService(a.logger)
	grpcArticleServer := transport_grpc.NewArticleSvc(articleService)

	seriesService := service.InitiateSeriesService(a.logger)
	grpcSeriesServer := transport_grpc.NewSeriesSvc(seriesService)

	seasonService := service.InitiateSeasonService(a.logger)
	episodeService := service.InitiateEpisodeService(a.logger)

	grpcSeasonServer := transport_grpc.NewSeasonSvc(seasonService)
	grpcEpisodeServer := transport_grpc.NewEpisodeSvc(episodeService)

	movieService := service.InitiateMovieService(a.logger)
	grpcMovieServer := transport_grpc.NewMovieSvc(movieService)

	listen, err := net.Listen("tcp", ":"+serverPort)
	if err != nil {
		a.logger.WithError(err).Fatal("Error while starting grpc server")
	}

	s := grpc.NewServer()

	pb.RegisterFileServerCelebritySvcServer(s, grpcCelebrityServer)
	pb.RegisterFileServerArticleSvcServer(s, grpcArticleServer)
	pb.RegisterFileServerSeriesSvcServer(s, grpcSeriesServer)
	pb.RegisterFileServerSeasonSvcServer(s, grpcSeasonServer)
	pb.RegisterFileServerEpisodeSvcServer(s, grpcEpisodeServer)
	pb.RegisterFileServerMovieSvcServer(s, grpcMovieServer)

	reflection.Register(s)
	err = s.Serve(listen)
	if err != nil {
		a.logger.WithError(err).Fatal("Error while serving grpc server")
	}
}

func (a *app) Run() error {
	// TODO allow origin only from your gateway
	origins := handlers.AllowedOrigins([]string{"*"})
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE"})

	err := http.ListenAndServe(":"+HttpPort, handlers.CORS(headers, origins, methods)(a.router))
	if err != nil {
		a.logger.Error("CONNECTION ERROR")
		return errors.Wrap(err, "connection error")
	}
	return nil
}

func (a *app) registerRoutes() {
	fs := http.FileServer(http.Dir("../int-static-files"))
	a.router.PathPrefix("/").
		Handler(http.StripPrefix("/", fs))
}
