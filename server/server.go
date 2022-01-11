package main

import (
	"context"
	"io"
	"math"
	"time"

	pb "github.com/Guilospanck/gRPC/route_guide/proto"
	"github.com/golang/protobuf/proto"
)

// implements the generated RouteGuideServer (line 148 of route_guide_grpc.pb.go)
type routeGuideServer struct {
	savedFeatures []*pb.Feature
}

func (s *routeGuideServer) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
	for _, feature := range s.savedFeatures {
		if proto.Equal(feature.Location, point) {
			return feature, nil
		}
	}

	// No feature was found, return an unnamed feature
	return &pb.Feature{Location: point}, nil
}

func (s *routeGuideServer) ListFeatures(rect *pb.Rectangle, stream pb.RouteGuide_ListFeaturesServer) error {
	for _, feature := range s.savedFeatures {
		if inRange(feature.Location, rect) {
			if err := stream.Send(feature); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *routeGuideServer) RecordRoute(stream pb.RouteGuide_RecordRouteServer) error {
	var pointCount, featureCount, distance int32
	var lastPoint *pb.Point

	startTime := time.Now()

	for {

		// get streaming of points
		point, err := stream.Recv()
		if err == io.EOF { // end of file (stream ended)
			endTime := time.Now()
			return stream.SendAndClose(&pb.RouteSummary{
				PointCount:   pointCount,
				FeatureCount: featureCount,
				Distance:     distance,
				ElapsedTime:  int32(endTime.Sub(startTime).Seconds()),
			})
		}

		if err != nil {
			return err
		}

		// increases number of points
		pointCount++

		// verifies for features
		for _, feature := range s.savedFeatures {
			if proto.Equal(feature.Location, point) {
				featureCount++
			}
		}

		// calculates distance
		if lastPoint != nil {
			distance += calcDistance(lastPoint, point)
		}

		// updates lastPoint
		lastPoint = point
	}

}

func (s *routeGuideServer) RouteChat(stream pb.RouteGuide_RouteChatServer) error {

}

// Helper function to verify if point is inside rectangle
func inRange(point *pb.Point, rect *pb.Rectangle) bool {
	left := math.Min(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	right := math.Max(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	bottom := math.Min(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))
	top := math.Max(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))

	if float64(point.Longitude) >= left &&
		float64(point.Longitude) <= right &&
		float64(point.Latitude) >= bottom &&
		float64(point.Latitude) <= top {
		return true
	}

	return false
}

func toRadians(num float64) float64 {
	// pi = 180º
	return num * math.Pi / float64(180)
}

// calcDistance calculates the distance between two points using the "haversine" formula.
/*
	One can derive Haversine formula to calculate distance between two as:

	a = sin²(ΔlatDifference/2) + cos(lat1).cos(lt2).sin²(ΔlonDifference/2)
	c = 2.atan2(√a, √(1−a))
	d = R.c

	where,

	ΔlatDifference = lat1 – lat2 (difference of latitude)
	ΔlonDifference = lon1 – lon2 (difference of longitude)
	R is radius of earth i.e 6371 KM or 3961 miles
	and d is the distance computed between two points.
*/
func calcDistance(point1, point2 *pb.Point) int32 {
	const CordFactor float64 = 1e7
	const R = float64(6371000) // earth radius in meters

	lat1 := toRadians(float64(point1.Latitude) / CordFactor)
	lat2 := toRadians(float64(point2.Latitude) / CordFactor)
	long1 := toRadians(float64(point1.Longitude) / CordFactor)
	long2 := toRadians(float64(point2.Longitude) / CordFactor)

	dLat := lat2 - lat1
	dLong := long2 - long1

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dLong/2)*math.Sin(dLong/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c

	return int32(distance)
}
