package presentation

import (
	"context"

	"github.com/Nickolasll/urlshortener/internal/app/config"

	pb "github.com/Nickolasll/urlshortener/internal/app/presentation/proto"
)

// ShortenServer поддерживает все необходимые методы сервера.
type ShortenServer struct {
	pb.UnimplementedShortenServer
}

func (s *ShortenServer) ShortenURL(ctx context.Context, req *pb.ShortenURLRequest) (*pb.ShortenURLResponse, error) {
	var response pb.ShortenURLResponse
	short, err := shorten(req.Url, req.UserId)
	if err != nil {
		slug, _ := getShortURLByOriginalURL(short.OriginalURL)
		response.Error = "StatusConflict"
		response.Result = *config.SlugEndpoint + slug
		return &response, nil
	}
	response.Result = *config.SlugEndpoint + short.ShortURL
	return &response, nil
}

func (s *ShortenServer) Expand(ctx context.Context, req *pb.ExpandRequest) (*pb.ExpandResponse, error) {
	var response pb.ExpandResponse
	value, err := expand(req.Slug)
	if err != nil {
		response.Error = "StatusNotFound"
		return &response, nil
	}
	if value.Deleted {
		response.Error = "StatusGone"
		return &response, nil
	}
	response.Location = value.OriginalURL
	return &response, nil
}

func (s *ShortenServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	response := pb.PingResponse{
		Status: "Ok",
	}
	if ping() != nil {
		response.Status = "InternalServerError"
	}
	return &response, nil
}

func (s *ShortenServer) BatchShorten(ctx context.Context, req *pb.BatchShortenRequest) (*pb.BatchShortenResponse, error) {
	var response pb.BatchShortenResponse
	var batchInput []BatchInput

	for _, requestItem := range req.Input {
		batchItem := BatchInput{
			CorrelationID: requestItem.CorrelationId,
			OriginalURL:   requestItem.OriginalUrl,
		}
		batchInput = append(batchInput, batchItem)
	}

	batchOutput := batchShorten(batchInput, req.UserId)

	for _, outputItem := range batchOutput {
		responseItem := pb.BatchShortenResponseItem{
			CorrelationId: outputItem.CorrelationID,
			ShortUrl:      outputItem.ShortURL,
		}
		response.Result = append(response.Result, &responseItem)
	}
	return &response, nil
}

func (s *ShortenServer) FindURLs(ctx context.Context, req *pb.FindURLsRequest) (*pb.FindURLsResponse, error) {
	var response pb.FindURLsResponse
	URLs := findURLs(req.UserId)
	for _, outputItem := range URLs {
		responseItem := pb.FindURLsResponseItem{
			OriginalUrl: outputItem.OriginalURL,
			ShortUrl:    outputItem.ShortURL,
		}
		response.Result = append(response.Result, &responseItem)
	}
	return &response, nil
}

func (s *ShortenServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	response := pb.DeleteResponse{
		Status: "Accepted",
	}
	go bulkDelete(req.ShortUrls, req.UserId)
	return &response, nil
}

func (s *ShortenServer) GetInternalStats(ctx context.Context, req *pb.GetInternalStatsRequest) (*pb.GetInternalStatsResponse, error) {
	stats, err := getInternalStats()
	response := pb.GetInternalStatsResponse{
		Urls:  int32(stats.URLs),
		Users: int32(stats.Users),
	}
	return &response, err
}
