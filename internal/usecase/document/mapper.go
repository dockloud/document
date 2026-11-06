package document

import "github.com/dockloud/document/internal/domain/document"

func DomainToOutputDTO(doc document.Document) *GetDocumentOutput {
	return &GetDocumentOutput{
		ID:        doc.ID,
		OwnerID:   doc.OwnerID,
		Name:      doc.Name,
		Content:   doc.Content,
		IsPrivate: doc.IsPrivate,
		CreatedAt: doc.CreatedAt,
		EditedAt:  doc.EditedAt,
	}
}

func DomainToCreateDTO(doc document.Document) *CreateOutput {
	return &CreateOutput{
		ID:        doc.ID,
		Name:      doc.Name,
		Content:   doc.Content,
		IsPrivate: doc.IsPrivate,
		CreatedAt: doc.CreatedAt,
	}
}

func DomainToUpdateDTO(doc document.Document) *UpdateOutput {
	return &UpdateOutput{
		ID:        doc.ID,
		Name:      doc.Name,
		Content:   doc.Content,
		IsPrivate: doc.IsPrivate,
		EditedAt:  doc.EditedAt,
	}
}
