package server

import (
	"context"
	"jobradar/matching-service/internal/scoring"
	matchingproto "jobradar/proto/gen/go/matching/v1"
)

type matchingServer struct {
	matchingproto.UnimplementedMatchingServiceServer
}

func NewMatchingServer() *matchingServer {
    return &matchingServer{}
}

func (ms *matchingServer) MatchScore(ctx context.Context, req *matchingproto.MatchScoreRequest) (*matchingproto.MatchScoreResponse, error) {
	return &matchingproto.MatchScoreResponse{
		Score: scoring.CalculateScore(req.Role),
	}, nil
}
