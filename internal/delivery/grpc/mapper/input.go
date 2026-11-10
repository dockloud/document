package mapper

import (
	docv1 "github.com/dockloud/contract/gen/document/v1"
	"github.com/dockloud/document/internal/usecase/document"
)

func ToGetAccessibleInput(req *docv1.GetAccessibleDocumentsRequest, uid string) *document.GetAccessibleInput {
	return &document.GetAccessibleInput{
		RequesterID: uid,
	}
}

func ToGetByAuthorIDInput(req *docv1.GetDocumentsByAuthorIDRequest, uid string) *document.GetByAuthorIDInput {
	return &document.GetByAuthorIDInput{
		AuthorID:    req.AuthorUid,
		RequesterID: uid,
	}
}

func ToGetByDocumentIDInput(req *docv1.GetDocumentByIDRequest, uid string) *document.GetByDocumentIDInput {
	return &document.GetByDocumentIDInput{
		DocumentID:  int(req.DocumentId),
		RequesterID: uid,
	}
}

func ToCreateInput(req *docv1.CreateDocumentRequest, uid string) *document.CreateInput {
	return &document.CreateInput{
		OwnerID:   uid,
		Name:      req.Name,
		Content:   req.Content,
		IsPrivate: req.IsPrivate,
	}
}

func ToUpdateInput(req *docv1.UpdateDocumentRequest, uid string) *document.UpdateInput {
	return &document.UpdateInput{
		OwnerID:   uid,
		ID:        int(req.DocumentId),
		Name:      req.Name,
		Content:   req.Content,
		IsPrivate: req.IsPrivate,
	}
}

func ToDeleteInput(req *docv1.DeleteDocumentRequest, uid string) *document.DeleteInput {
	return &document.DeleteInput{
		DocumentID:  int(req.DocumentId),
		RequesterID: uid,
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
