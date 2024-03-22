package main

import (
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to run")
	}
}

func run() error {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	log.Info().Msg("listening on :3000")

	server := grpc.NewServer()

	server.RegisterService(
		&grpc.ServiceDesc{
			ServiceName: "moby.buildkit.v1.Control",
			HandlerType: (*interface{})(nil),
			Streams: []grpc.StreamDesc{
				{
					StreamName: "ListWorkers",
					Handler: func(srv interface{}, serverStream grpc.ServerStream) error {
						fullMethodName, ok := grpc.MethodFromServerStream(serverStream)
						if !ok {
							return status.Errorf(codes.Internal, "lowLevelServerStream not exists in context")
						}

						log.Info().Str("fullMethodName", fullMethodName).Interface("srv", srv).Msg("ListWorkers")

						md, ok := metadata.FromIncomingContext(serverStream.Context())
						if !ok {
							return status.Errorf(codes.Internal, "metadata not exists in context")
						}

						log.Info().Interface("metadata", md).Msg("metadata")

						return nil
					},
					ServerStreams: true,
					ClientStreams: true,
				},
			},
		},
		nil,
	)

	err = server.Serve(listener)
	if err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
