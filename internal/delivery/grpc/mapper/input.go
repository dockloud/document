package mapper

import (
	docv1 "github.com/dockloud/contract/gen/document/v1"
	"github.com/dockloud/document/internal/usecase/document"
)

func ToGetAccessibleInput(req *docv1.GetAccessibleDocumentsRequest) *document.GetAccessibleInput {
	return &document.GetAccessibleInput{
		RequesterID: req.RequesterUid,
	}
}

func ToGetByAuthorIDInput(req *docv1.GetDocumentsByAuthorIDRequest) *document.GetByAuthorIDInput {
	return &document.GetByAuthorIDInput{
		AuthorID:    req.AuthorUid,
		RequesterID: req.RequesterUid,
	}
}

func ToGetByDocumentIDInput(req *docv1.GetDocumentByIDRequest) *document.GetByDocumentIDInput {
	return &document.GetByDocumentIDInput{
		DocumentID:  int(req.DocumentId),
		RequesterID: req.RequesterUid,
	}
}

func ToCreateInput(req *docv1.CreateDocumentRequest) *document.CreateInput {
	return &document.CreateInput{
		OwnerID:   req.RequesterUid,
		Name:      req.Name,
		Content:   req.Content,
		IsPrivate: req.IsPrivate,
	}
}

func ToUpdateInput(req *docv1.UpdateDocumentRequest) *document.UpdateInput {
	return &document.UpdateInput{
		OwnerID:   req.RequesterUid,
		ID:        int(req.DocumentId),
		Name:      req.Name,
		Content:   req.Content,
		IsPrivate: req.IsPrivate,
	}
}

func ToDeleteInput(req *docv1.DeleteDocumentRequest) *document.DeleteInput {
	return &document.DeleteInput{
		DocumentID:  int(req.DocumentId),
		RequesterID: req.RequesterUid,
	}
}

func ToGetUnsafeInput(req *docv1.GetUnsafeDocumentRequest) *document.GetUnsafeInput {
	return &document.GetUnsafeInput{
		DocumentID: int(req.DocumentId),
	}
}

func ToUpdateUnsafeInput(req *docv1.UpdateUnsafeDocumentRequest) *document.UpdateUnsafeInput {
	return &document.UpdateUnsafeInput{
		ID:        int(req.DocumentId),
		Name:      req.Name,
		Content:   req.Content,
		IsPrivate: req.IsPrivate,
	}
}

func ToDeleteUnsafeInput(req *docv1.DeleteUnsafeDocumentRequest) *document.DeleteUnsafeInput {
	return &document.DeleteUnsafeInput{
		DocumentID: int(req.DocumentId),
	}
}
