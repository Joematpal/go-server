package streamer

import (
	"io"
	"net/http"

	"github.com/joematpal/go-server/example/server/internal/logger"
	streamer "github.com/joematpal/go-server/pkg/streamer/v1"
)

type RouteGuide struct {
	logger logger.Logger
	streamer.UnimplementedStreamerServer
}

func New(logr logger.Logger) *RouteGuide {
	return &RouteGuide{
		logger: logr,
	}
}

func (rg *RouteGuide) StreamPoint(stream streamer.Streamer_StreamPointServer) error {
	for {
		point, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&streamer.Status{
				Status:  http.StatusOK,
				Message: "success",
			})
		}
		if err != nil {
			return err
		}

		rg.logger.Infof("point: %+v", point)
	}
}
