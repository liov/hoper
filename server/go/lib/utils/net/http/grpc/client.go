package grpci

import (
	"crypto/tls"
	"github.com/liov/hoper/server/go/lib/utils/net/http/grpc/stats"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"strings"
)

func GetDefaultClient(target string) (*grpc.ClientConn, error) {
	// Set up a connection to the server.
	return grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(&stats.ClientHandler{}))
}

func GetTlsClient(target string) (*grpc.ClientConn, error) {
	// Set up a connection to the server.
	return grpc.Dial(target, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{ServerName: strings.Split(target, ":")[0], InsecureSkipVerify: true})),
		grpc.WithStatsHandler(&stats.ClientHandler{}))
}
