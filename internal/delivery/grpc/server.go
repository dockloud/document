package grpc

import (
	"context"
	"log"

	docv1 "github.com/dockloud/contract/gen/document/v1"
	"github.com/dockloud/document/internal/delivery/grpc/mapper"
	"github.com/dockloud/document/internal/usecase/document"
)

type documentServer struct {
	docv1.UnimplementedDocumentServiceServer
	service *document.Service
	logger  *log.Logger
}

func NewDocumentServer(service document.Service) docv1.DocumentServiceServer {
	return &documentServer{
		service: &service,
		logger:  log.Default(),
	}
}

func (s *documentServer) GetDocumentsByAuthorID(ctx context.Context, req *docv1.GetDocumentsByAuthorIDRequest) (*docv1.GetDocumentsByAuthorIDResponse, error) {
	s.logger.Printf("Received GetDocumentsByAuthorIDRequest %+v", req)

	input := mapper.ToGetByAuthorIDInput(req)
	res, err := s.service.GetByAuthorID(ctx, input)

	if err != nil {
		s.logger.Printf("Failed to register %+v", req) // Примечание: текст лога сохранен из вашего шаблона
		return nil, err
	}

	return mapper.ToGetDocumentsByAuthorIDResponse(res), nil
}

func (s *documentServer) GetAccessibleDocuments(ctx context.Context, req *docv1.GetAccessibleDocumentsRequest) (*docv1.GetAccessibleDocumentsResponse, error) {
	s.logger.Printf("Received GetAccessibleDocumentsRequest %+v", req)

	input := mapper.ToGetAccessibleInput(req)
	res, err := s.service.GetAccessible(ctx, input)

	if err != nil {
		s.logger.Printf("Failed to get accessible documents %+v", req)
		return nil, err
	}

	return mapper.ToGetAccessibleDocumentsResponse(res), nil
}

func (s *documentServer) GetDocumentByID(ctx context.Context, req *docv1.GetDocumentByIDRequest) (*docv1.GetDocumentByIDResponse, error) {
	s.logger.Printf("Received GetDocumentByIDRequest %+v", req)

	input := mapper.ToGetByDocumentIDInput(req)
	res, err := s.service.GetByDocumentID(ctx, input)

	if err != nil {
		s.logger.Printf("Failed to get document by ID %+v", req)
		return nil, err
	}

	return mapper.ToGetDocumentByIDResponse(res), nil
}

func (s *documentServer) CreateDocument(ctx context.Context, req *docv1.CreateDocumentRequest) (*docv1.CreateDocumentResponse, error) {
	s.logger.Printf("Received CreateDocumentRequest %+v", req)

	input := mapper.ToCreateInput(req)
	res, err := s.service.Create(ctx, input)

	if err != nil {
		s.logger.Printf("Failed to create document %+v", req)
		return nil, err
	}

	return mapper.ToCreateDocumentResponse(res), nil
}

func (s *documentServer) UpdateDocument(ctx context.Context, req *docv1.UpdateDocumentRequest) (*docv1.UpdateDocumentResponse, error) {
	s.logger.Printf("Received UpdateDocumentRequest %+v", req)

	input := mapper.ToUpdateInput(req)
	res, err := s.service.Update(ctx, input)

	if err != nil {
		s.logger.Printf("Failed to update document %+v", req)
		return nil, err
	}

	return mapper.ToUpdateDocumentResponse(res), nil
}

func (s *documentServer) DeleteDocument(ctx context.Context, req *docv1.DeleteDocumentRequest) (*docv1.DeleteDocumentResponse, error) {
	s.logger.Printf("Received DeleteDocumentRequest %+v", req)

	input := mapper.ToDeleteInput(req)
	err := s.service.Delete(ctx, input)

	if err != nil {
		s.logger.Printf("Failed to delete document %+v", req)
		return nil, err
	}

	return mapper.ToDeleteDocumentResponse(), nil
}

func (s *documentServer) GetUnsafeDocument(ctx context.Context, req *docv1.GetUnsafeDocumentRequest) (*docv1.GetUnsafeDocumentResponse, error) {
	s.logger.Printf("Received GetUnsafeDocumentRequest %+v", req)

	input := mapper.ToGetUnsafeInput(req)
	res, err := s.service.GetUnsafe(ctx, input)

	if err != nil {
		s.logger.Printf("Failed to get unsafe document %+v", req)
		return nil, err
	}

	return mapper.ToGetUnsafeDocumentResponse(res), nil
}

func (s *documentServer) UpdateUnsafeDocument(ctx context.Context, req *docv1.UpdateUnsafeDocumentRequest) (*docv1.UpdateUnsafeDocumentResponse, error) {
	s.logger.Printf("Received UpdateUnsafeDocumentRequest %+v", req)

	input := mapper.ToUpdateUnsafeInput(req)
	res, err := s.service.UpdateUnsafe(ctx, input)

	if err != nil {
		s.logger.Printf("Failed to update unsafe document %+v", req)
		return nil, err
	}

	return mapper.ToUpdateUnsafeDocumentResponse(res), nil
}

func (s *documentServer) DeleteUnsafeDocument(ctx context.Context, req *docv1.DeleteUnsafeDocumentRequest) (*docv1.DeleteUnsafeDocumentResponse, error) {
	s.logger.Printf("Received DeleteUnsafeDocumentRequest %+v", req)

	input := mapper.ToDeleteUnsafeInput(req)
	err := s.service.DeleteUnsafe(ctx, input)

	if err != nil {
		s.logger.Printf("Failed to delete unsafe document %+v", req)
		return nil, err
	}

	return mapper.ToDeleteUnsafeDocumentResponse(), nil
}
