package document

import (
	"context"

	"github.com/dockloud/document/internal/domain/document"
)

type Repository interface {
	GetByAuthorID(ctx context.Context, authorID string, requesterID string) ([]document.Document, error)
	GetAccessible(ctx context.Context, uid string) ([]document.Document, error)
	GetByID(ctx context.Context, documentID int, requesterID string) (*document.Document, error)
	Create(ctx context.Context, ownerID string, name string, content string, isPrivate bool) (*document.Document, error)
	Delete(ctx context.Context, documentID int, requesterID string) (*document.Document, error)
	Update(ctx context.Context, documentID int, ownerID string, name *string, content *string, isPrivate *bool) (*document.Document, error)

	GetUnsafe(ctx context.Context, documentID int) (*document.Document, error)
	UpdateUnsafe(ctx context.Context, documentID int, name *string, content *string, isPrivate *bool) (*document.Document, error)
	DeleteUnsafe(ctx context.Context, documentID int) (*document.Document, error)
}

type Service struct {
	documentRepo Repository
}

func NewService(documentRepo Repository) *Service {
	return &Service{
		documentRepo: documentRepo,
	}
}

func (s *Service) GetByAuthorID(ctx context.Context, input *GetByAuthorIDInput) ([]GetDocumentOutput, error) {
	var mapped []GetDocumentOutput
	res, err := s.documentRepo.GetByAuthorID(ctx, input.AuthorID, input.RequesterID)
	if err != nil {
		return nil, err
	}
	for _, doc := range res {
		output := DomainToOutputDTO(doc)
		mapped = append(mapped, *output)
	}
	return mapped, err
}

func (s *Service) GetAccessible(ctx context.Context, input *GetAccessibleInput) ([]GetDocumentOutput, error) {
	var mapped []GetDocumentOutput
	res, err := s.documentRepo.GetAccessible(ctx, input.RequesterID)
	if err != nil {
		return nil, err
	}
	for _, doc := range res {
		output := DomainToOutputDTO(doc)
		mapped = append(mapped, *output)
	}
	return mapped, err
}

func (s *Service) GetByDocumentID(ctx context.Context, input *GetByDocumentIDInput) (*GetDocumentOutput, error) {
	doc, err := s.documentRepo.GetByID(ctx, input.DocumentID, input.RequesterID)
	if err != nil {
		return nil, err
	}
	return DomainToOutputDTO(*doc), nil
}

func (s *Service) Create(ctx context.Context, input *CreateInput) (*CreateOutput, error) {
	if err := document.ValidateName(input.Name); err != nil {
		return nil, err
	} else if err = document.ValidateContent(input.Content); err != nil {
		return nil, err
	}

	res, err := s.documentRepo.Create(ctx, input.OwnerID, input.Name, input.Content, input.IsPrivate)
	if err != nil {
		return nil, err
	}
	return DomainToCreateDTO(*res), err
}

func (s *Service) Delete(ctx context.Context, input *DeleteInput) error {
	_, err := s.documentRepo.Delete(ctx, input.DocumentID, input.RequesterID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Update(ctx context.Context, input *UpdateInput) (*UpdateOutput, error) {
	res, err := s.documentRepo.Update(ctx, input.ID, input.OwnerID, input.Name, input.Content, input.IsPrivate)
	if err != nil {
		return nil, err
	}
	return DomainToUpdateDTO(*res), err
}

func (s *Service) GetUnsafe(ctx context.Context, input *GetUnsafeInput) (*GetDocumentOutput, error) {
	res, err := s.documentRepo.GetUnsafe(ctx, input.DocumentID)
	if err != nil {
		return nil, err
	}
	return DomainToOutputDTO(*res), err
}

func (s *Service) DeleteUnsafe(ctx context.Context, input *DeleteUnsafeInput) error {
	_, err := s.documentRepo.DeleteUnsafe(ctx, input.DocumentID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateUnsafe(ctx context.Context, input *UpdateUnsafeInput) (*UpdateOutput, error) {
	res, err := s.documentRepo.UpdateUnsafe(ctx, input.ID, input.Name, input.Content, input.IsPrivate)
	if err != nil {
		return nil, err
	}
	return DomainToUpdateDTO(*res), err
}
