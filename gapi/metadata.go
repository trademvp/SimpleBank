package gapi

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-host"
)

type Metadata struct {
	UserAgent string
	ClintIP   string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	metaData := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			metaData.UserAgent = userAgents[0]
		}
		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			metaData.UserAgent = userAgents[0]
		}
		if userAgents := md.Get(xForwardedForHeader); len(userAgents) > 0 {
			metaData.UserAgent = userAgents[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		metaData.ClintIP = p.Addr.String()
	}

	return metaData
}
