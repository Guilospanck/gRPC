package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/Guilospanck/gRPC/route_guide/proto"
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
	printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})

}
