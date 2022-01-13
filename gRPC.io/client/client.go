package main

import (
	"context"
	"flag"
	"io"
	"log"
	"math/rand"
	"time"

	pb "github.com/Guilospanck/gRPC/gRPC.io/route_guide/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/examples/data"
)

var (
	serverAddr         = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverHostOverride = flag.String("server_host_override", "x.test.example.com", "The server name used to verify the hostname returned by the TLS handshake")
)

func printFeature(client pb.RouteGuideClient, point *pb.Point) {
	log.Printf("Getting feature for point (%d, %d)", point.Latitude, point.Longitude)

	// cancels operation if timeout exceeds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	feature, err := client.GetFeature(ctx, point)
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v: ", client, err)
	}
	log.Println(feature)
}

func printFeatures(client pb.RouteGuideClient, rectangle *pb.Rectangle) {
	log.Printf("Looking for features in rectangle %v", rectangle)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.ListFeatures(ctx, rectangle)
	if err != nil {
		log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
	}

	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
		}
		log.Println(feature)
	}
}

func runRecordRoute(client pb.RouteGuideClient) {
	// creates a random number of random points
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pointCount := int(r.Int31n(100)) + 2 // traverse at least two points

	var points []*pb.Point
	for i := 0; i < pointCount; i++ {
		points = append(points, randomPoint(r))
	}
	log.Printf("Traversing %d points.", len(points))

	// context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.RecordRoute(ctx)
	if err != nil {
		log.Fatalf("%v.RecordRoute(_) = _, %v", client, err)
	}

	// stream data points to server
	for _, point := range points {
		if err := stream.Send(point); err != nil {
			log.Fatalf("%v.Send(%v) = %v", stream, point, err)
		}
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	log.Printf("Route summary: %v", reply)
}

func runRouteChat(client pb.RouteGuideClient) {
	notes := []*pb.RouteNote{
		{Location: &pb.Point{Latitude: 0, Longitude: 1}, Message: "First message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 2}, Message: "Second message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 3}, Message: "Third message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 4}, Message: "Fourth message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 5}, Message: "Fifth message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 6}, Message: "Sixth message"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get stream from RouteChat
	stream, err := client.RouteChat(ctx)
	if err != nil {
		log.Fatalf("%v.RouteChat(ctx) = _, %v", client, err)
	}

	waitc := make(chan struct{})

	// reads from stream
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("%v.Recv() = _, %v", stream, err)
			}
			log.Printf("Got message %s at point(%d, %d)", in.Message, in.Location.Latitude, in.Location.Longitude)
		}
	}()

	// writes to stream
	for _, note := range notes {
		if err := stream.Send(note); err != nil {
			log.Fatalf("%v.Send(_) = %v", stream, err)
		}
	}

	stream.CloseSend()
	<-waitc
}

func randomPoint(r *rand.Rand) *pb.Point {
	lat := (r.Int31n(180) - 90) * 1e7
	long := (r.Int31n(360) - 180) * 1e7
	return &pb.Point{Latitude: lat, Longitude: long}
}

func main() {
	flag.Parse()

	var opts []grpc.DialOption

	// get TLS options
	if *tls {
		if *caFile == "" {
			*caFile = data.Path("x509/ca_cert.pem")
		}

		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// dial to server
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	// client stub
	client := pb.NewRouteGuideClient(conn)

	/*
	   In gRPC-Go, RPCs operate in blocking/synchronous mode, which means that the RPC
	   call wait for the server to respond, and wil either return a response or an error.
	*/

	// simple RPC "Get Feature" valid feature
	// printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})
	// feature missing
	// printFeature(client, &pb.Point{Latitude: 0, Longitude: 0})

	// Looking for features between 40, -75 and 42, -73.
	// printFeatures(client, &pb.Rectangle{
	// 	Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
	// 	Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
	// })

	// RecordRoute
	// runRecordRoute(client)

	// RouteChat
	runRouteChat(client)

}
