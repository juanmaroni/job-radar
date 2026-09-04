package matching

import (
	matchingproto "jobradar/proto/gen/go/matching/v1"
	"jobradar/tracker-service/internal/config/environment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGrpcClient() (matchingproto.MatchingServiceClient, error) {

	matchingHost, err := environment.GetEnvVar("TRACKER_SERVICE_MATCHING_HOST")

	if err != nil {
		return nil, err
	}

	conn, err := grpc.NewClient(matchingHost, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}
	
	return matchingproto.NewMatchingServiceClient(conn), nil
}
