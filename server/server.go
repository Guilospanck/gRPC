package main

import (
	"context"

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

}

func (s *routeGuideServer) RecordRoute(stream pb.RouteGuide_RecordRouteServer) error {

}

func (s *routeGuideServer) RouteChat(stream pb.RouteGuide_RouteChatServer) error {

}
