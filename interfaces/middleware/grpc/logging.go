package middleware

import (
	"context"
	"time"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"
	"google.golang.org/grpc"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LoggingUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	startTime := time.Now()

	resp, err := handler(ctx, req)

	st, _ := status.FromError(err)
	statusCode := grpcCodes.OK
	if err != nil {
		statusCode = st.Code()
	}

	log.Info(
		"Access log",
		log.Ftime("Date", startTime),
		log.Fstring("Method", info.FullMethod),
		log.Fint("StatusCode", int(statusCode)),
		log.Fduration("Duration", time.Since(startTime)),
	)

	if statusCode != grpcCodes.OK {
		log.Error(
			"Error log",
			log.Ftime("Date", startTime),
			log.Fstring("Method", info.FullMethod),
			log.Fint("StatusCode", int(statusCode)),
			log.Fstring("Error", err.Error()),
		)
	}

	return resp, err
}
