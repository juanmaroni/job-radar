package server

import (
	"context"
	matching "jobradar/proto/gen/go/matching/v1"
)

type matchingServer struct {
	matching.UnimplementedMatchingServiceServer
}

func NewMatchingServer() *matchingServer {
    return &matchingServer{}
}

func (ms *matchingServer) MatchScore(ctx context.Context, req *matching.MatchScoreRequest) (*matching.MatchScoreResponse, error) {
	matchResp := matching.MatchScoreResponse{
		Score: int32(len(req.Role)) % 100,
	}

	return &matchResp, nil
}
