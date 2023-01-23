package streamer

import (
	"context"
	"net/http"

	connect_go "github.com/bufbuild/connect-go"
	"github.com/joematpal/go-server/example/server/internal/logger"
	streamer_v1 "github.com/joematpal/go-server/pkg/streamer/v1"
	"github.com/joematpal/go-server/pkg/streamer/v1/streamer_v1connect"
)

type RouteGuide struct {
	logger logger.Logger
	streamer_v1connect.UnimplementedStreamerServiceHandler
}

func New(logr logger.Logger) *RouteGuide {
	return &RouteGuide{
		logger: logr,
	}
}

func (rg *RouteGuide) StreamPoint(ctx context.Context, stream *connect_go.ClientStream[streamer_v1.Point]) (*connect_go.Response[streamer_v1.Status], error) {

	for stream.Receive() {
		point := stream.Msg()
		rg.logger.Infof("point: %+v", point)
	}

	if err := stream.Err(); err != nil {
		return nil, err
	}

	return &connect_go.Response[streamer_v1.Status]{
		Msg: &streamer_v1.Status{
			Status:  http.StatusAccepted,
			Message: "ok",
		},
	}, nil
}
