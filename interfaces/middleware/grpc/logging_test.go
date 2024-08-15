package middleware

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_LoggingUnaryInterceptor(t *testing.T) {
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)

	tests := []struct {
		name            string
		handler         func(context.Context, interface{}) (interface{}, error)
		expectedStatus  codes.Code
		expectAccessLog bool
		expectErrorLog  bool
	}{
		{
			name: "successful request",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return "OK", nil
			},
			expectedStatus:  codes.OK,
			expectAccessLog: true,
			expectErrorLog:  false,
		},
		{
			name: "client error request",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, status.Error(codes.InvalidArgument, "Invalid argument")
			},
			expectedStatus:  codes.InvalidArgument,
			expectAccessLog: true,
			expectErrorLog:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoggingUnaryInterceptor(
				context.Background(),
				nil, // dummy request
				&grpc.UnaryServerInfo{
					FullMethod: "/test.TestService/TestMethod",
				},
				tt.handler,
			)

			if status.Code(err) != tt.expectedStatus {
				t.Errorf("expected status code %v, got %v", tt.expectedStatus, status.Code(err))
			}

			logOutput := logBuf.String()
			if tt.expectAccessLog && !strings.Contains(logOutput, "Access log") {
				t.Errorf("expected access log, got %s", logOutput)
			}

			if tt.expectErrorLog && !strings.Contains(logOutput, "Error log") {
				t.Errorf("expected error log, got %s", logOutput)
			}

			logBuf.Reset()
		})
	}
}
