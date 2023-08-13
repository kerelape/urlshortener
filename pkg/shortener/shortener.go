package shortener

import (
	"context"

	"github.com/kerelape/urlshortener/internal/app"
	"github.com/kerelape/urlshortener/internal/app/model"
	pb "github.com/kerelape/urlshortener/pkg/shortener/proto"
)

// Server is grpc server.
type Server struct {
	pb.UnimplementedShortenerServer

	shortener model.Shortener
}

// NewServer returns a new Server.
func NewServer(shortener model.Shortener) *Server {
	return &Server{
		shortener: shortener,
	}
}

// Shorten shortens url.
func (s *Server) Shorten(ctx context.Context, in *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	var result, err = s.shortener.Shorten(ctx, app.NewToken(), in.Url)
	if err != nil {
		return nil, err
	}
	return &pb.ShortenResponse{Result: result}, nil
}

// Reveal reveals url.
func (s *Server) Reveal(ctx context.Context, in *pb.RevealRequest) (*pb.RevealResponse, error) {
	var origin, err = s.shortener.Reveal(ctx, in.ShortUrl)
	if err != nil {
		return nil, err
	}
	return &pb.RevealResponse{OriginalUrl: origin}, nil
}
