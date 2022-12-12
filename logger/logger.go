package logger

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"github.com/google/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var logs *logger.Logger
var console *logger.Logger

func Init() {
	service := os.Getenv("SERVICE_NAME")
	lf, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}

	logs = logger.Init(service, false, false, lf)
	console = logger.Init(service, true, false, ioutil.Discard)
}

// LogMessage - function to log colored messages
func LogMessage(level string, message string, args ...interface{}) {
	switch level {
	case "info":
		consolemessage := "\033[1;34m" + message + "\033[0m"
		console.Infof(consolemessage, args...)
		logs.Infof(message, args...)
	case "debug":
		consolemessage := "\033[1;33m-----" + message + "-----\033[0m"
		console.Warningf(consolemessage, args...)
		logs.Warningf(message, args...)
	case "error":
		consolemessage := "\033[1;31m" + message + "\033[0m"
		console.Errorf(consolemessage, args...)
		logs.Errorf(message, args...)
	default:
		consolemessage := "\033[1;36m" + message + "\033[0m"
		console.Infof(consolemessage, args...)
		logs.Infof(message, args...)
	}
}

func GRPCLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	level := "info"
	startTime := time.Now()
	result, err := handler(ctx, req)
	if err != nil {
		level = "error"
	}
	elapsedTime := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	LogMessage(level, "%s, %d, %s, %s", info.FullMethod, int(statusCode), statusCode.String(), elapsedTime)
	return result, err
}
