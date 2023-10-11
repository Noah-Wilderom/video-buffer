package helpers

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"time"
)

func ListenGRPC(port int, register func(s *grpc.Server)) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalln("Failed to listen for gRPC:", err)
	}

	s := grpc.NewServer()

	register(s)

	log.Println("gRPC Server started on port", port)

	if err = s.Serve(lis); err != nil {
		log.Fatalln("Failed to listen for gRPC:", err)
	}
}

func SendPayloadThroughGRPC(target string, callback func(conn *grpc.ClientConn, ctx context.Context)) error {
	conn, err := grpc.Dial(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	callback(conn, ctx)

	return nil
}
