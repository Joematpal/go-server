package routeguide

import (
	"io"
	"math"
	"time"

	"github.com/digital-dream-labs/go-server/example/server/internal/logger"
	"google.golang.org/grpc/examples/route_guide/routeguide"
	// "github.com/digital-dream-labs/go-server/pkg/streamer/v1"
)

type RouteGuide struct {
	logger logger.Logger
	routeguide.UnimplementedRouteGuideServer
}

func New(logr logger.Logger) *RouteGuide {
	return &RouteGuide{
		logger: logr,
	}
}

func (rg *RouteGuide) RecordRoute(stream routeguide.RouteGuide_RecordRouteServer) error {
	var pointCount, featureCount, distance int32
	var lastPoint *routeguide.Point
	startTime := time.Now()
	for {
		point, err := stream.Recv()
		if err == io.EOF {
			endTime := time.Now()
			return stream.SendAndClose(&routeguide.RouteSummary{
				PointCount:   pointCount,
				FeatureCount: featureCount,
				Distance:     distance,
				ElapsedTime:  int32(endTime.Sub(startTime).Seconds()),
			})
		}
		if err != nil {
			return err
		}

		rg.logger.Infof("point: %+v", point)
		pointCount++
		// for _, feature := range rg.savedFeatures {
		// 	if proto.Equal(feature.Location, point) {
		// 		featureCount++
		// 	}
		// }
		if lastPoint != nil {
			distance += calcDistance(lastPoint, point)
		}
		lastPoint = point
	}
}

// calcDistance calculates the distance between two points using the "haversine" formula.
// The formula is based on http://mathforum.org/library/drmath/view/51879.html.
func calcDistance(p1 *routeguide.Point, p2 *routeguide.Point) int32 {
	const CordFactor float64 = 1e7
	const R = float64(6371000) // earth radius in metres
	lat1 := toRadians(float64(p1.Latitude) / CordFactor)
	lat2 := toRadians(float64(p2.Latitude) / CordFactor)
	lng1 := toRadians(float64(p1.Longitude) / CordFactor)
	lng2 := toRadians(float64(p2.Longitude) / CordFactor)
	dlat := lat2 - lat1
	dlng := lng2 - lng1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlng/2)*math.Sin(dlng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c
	return int32(distance)
}

func toRadians(num float64) float64 {
	return num * math.Pi / float64(180)
}
