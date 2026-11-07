package mapper

import (
	docv1 "github.com/dockloud/contract/gen/document/v1"
	"github.com/dockloud/document/internal/usecase/document"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// --- OUTPUT MAPPERS ---

func ToGetDocumentByIDResponse(out *document.GetDocumentOutput) *docv1.GetDocumentByIDResponse {
	return &docv1.GetDocumentByIDResponse{
		DocumentId: int64(out.ID),
		OwnerUid:   out.OwnerID,
		Name:       out.Name,
		Content:    out.Content,
		IsPrivate:  out.IsPrivate,
		CreatedAt:  timestamppb.New(out.CreatedAt),
		EditedAt:   timestamppb.New(out.EditedAt),
	}
}

func ToGetAccessibleDocumentsResponse(docs []document.GetDocumentOutput) *docv1.GetAccessibleDocumentsResponse {
	res := make([]*docv1.GetDocumentByIDResponse, len(docs))
	for i, doc := range docs {
		res[i] = ToGetDocumentByIDResponse(&doc)
	}
	return &docv1.GetAccessibleDocumentsResponse{
		Documents: res,
	}
}

func ToGetDocumentsByAuthorIDResponse(docs []document.GetDocumentOutput) *docv1.GetDocumentsByAuthorIDResponse {
	res := make([]*docv1.GetDocumentByIDResponse, len(docs))
	for i, doc := range docs {
		res[i] = ToGetDocumentByIDResponse(&doc)
	}
	return &docv1.GetDocumentsByAuthorIDResponse{
		Documents: res,
	}
}

func ToCreateDocumentResponse(out *document.CreateOutput) *docv1.CreateDocumentResponse {
	return &docv1.CreateDocumentResponse{
		DocumentId: int64(out.ID),
		Name:       new(out.Name),
		Content:    new(out.Content),
		IsPrivate:  new(out.IsPrivate),
		CreatedAt:  timestamppb.New(out.CreatedAt),
	}
}

func ToUpdateDocumentResponse(out *document.UpdateOutput) *docv1.UpdateDocumentResponse {
	return &docv1.UpdateDocumentResponse{
		DocumentId: int64(out.ID),
		Name:       new(out.Name),
		Content:    new(out.Content),
		IsPrivate:  new(out.IsPrivate),
		EditedAt:   timestamppb.New(out.EditedAt),
	}
}

func ToDeleteDocumentResponse() *docv1.DeleteDocumentResponse {
	return &docv1.DeleteDocumentResponse{}
}

func ToGetUnsafeDocumentResponse(out *document.GetDocumentOutput) *docv1.GetUnsafeDocumentResponse {
	return &docv1.GetUnsafeDocumentResponse{
		DocumentId: int64(out.ID),
		OwnerUid:   out.OwnerID,
		Name:       out.Name,
		Content:    out.Content,
		IsPrivate:  out.IsPrivate,
		CreatedAt:  timestamppb.New(out.CreatedAt),
		EditedAt:   timestamppb.New(out.EditedAt),
	}
}

func ToUpdateUnsafeDocumentResponse(out *document.UpdateOutput) *docv1.UpdateUnsafeDocumentResponse {
	return &docv1.UpdateUnsafeDocumentResponse{
		DocumentId: int64(out.ID),
		Name:       new(out.Name),
		Content:    new(out.Content),
		IsPrivate:  new(out.IsPrivate),
		EditedAt:   timestamppb.New(out.EditedAt),
	}
}

func ToDeleteUnsafeDocumentResponse() *docv1.DeleteUnsafeDocumentResponse {
	return &docv1.DeleteUnsafeDocumentResponse{}
}
