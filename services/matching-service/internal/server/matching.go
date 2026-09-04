package server

import (
	"context"
	matchingproto "jobradar/proto/gen/go/matching/v1"
)

type matchingServer struct {
	matchingproto.UnimplementedMatchingServiceServer
}

func NewMatchingServer() *matchingServer {
    return &matchingServer{}
}

func (ms *matchingServer) MatchScore(ctx context.Context, req *matchingproto.MatchScoreRequest) (*matchingproto.MatchScoreResponse, error) {
	matchResp := matchingproto.MatchScoreResponse{
		Score: int32(len(req.Role)) % 100,
	}

	return &matchResp, nil
}
