package document

import "time"

type GetAccessibleInput struct {
	RequesterID string
}

type GetByAuthorIDInput struct {
	AuthorID    string
	RequesterID string
}

type GetByDocumentIDInput struct {
	DocumentID  int
	RequesterID string
}

type GetDocumentOutput struct {
	ID        int
	OwnerID   string
	Name      string
	Content   string
	IsPrivate bool
	CreatedAt time.Time
	EditedAt  time.Time
}

type CreateInput struct {
	OwnerID   string
	Name      string
	Content   string
	IsPrivate bool
}

type CreateOutput struct {
	ID        int
	Name      string
	Content   string
	IsPrivate bool
	CreatedAt time.Time
}

type UpdateInput struct {
	OwnerID   string
	ID        int
	Name      *string
	Content   *string
	IsPrivate *bool
}

type UpdateOutput struct {
	ID        int
	Name      string
	Content   string
	IsPrivate bool
	EditedAt  time.Time
}

type DeleteInput struct {
	DocumentID  int
	RequesterID string
}

type GetUnsafeInput struct {
	DocumentID int
}

type UpdateUnsafeInput struct {
	ID        int
	Name      *string
	Content   *string
	IsPrivate *bool
}

type DeleteUnsafeInput struct {
	DocumentID int
}
