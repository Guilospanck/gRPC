package main

import (
	"context"
	"math"

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
