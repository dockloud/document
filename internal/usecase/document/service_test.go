package document

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/dockloud/document/internal/domain/document"
	"github.com/dockloud/document/internal/usecase/document/mocks"
	"github.com/golang/mock/gomock"
)

var (
	createdAt = time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	editedAt  = time.Date(2025, 2, 20, 12, 0, 0, 0, time.UTC)
)

// Один и тот же экземпляр ошибки и в setup, и в wantErr — чтобы errors.Is работал.
var (
	errGetByAuthorFailed   = errors.New("get by author failed")
	errGetAccessibleFailed = errors.New("get accessible failed")
	errGetByIDFailed       = errors.New("get by id failed")
	errCreateFailed        = errors.New("create failed")
	errDeleteFailed        = errors.New("delete failed")
	errUpdateFailed        = errors.New("update failed")
	errGetUnsafeFailed     = errors.New("get unsafe failed")
	errUpdateUnsafeFailed  = errors.New("update unsafe failed")
	errDeleteUnsafeFailed  = errors.New("delete unsafe failed")
)

func strPtr(s string) *string { return &s }

func boolPtr(b bool) *bool { return &b }

// nilStrPtrMatcher матчит типизированный nil *string в interface{}.
// gomock.Nil() не подходит: он матчит только не типизированный nil.
type nilStrPtrMatcher struct{}

func (nilStrPtrMatcher) Matches(x interface{}) bool {
	s, ok := x.(*string)
	return ok && s == nil
}

func (nilStrPtrMatcher) String() string { return "nil *string" }

// nilBoolPtrMatcher матчит типизированный nil *bool в interface{}.
type nilBoolPtrMatcher struct{}

func (nilBoolPtrMatcher) Matches(x interface{}) bool {
	b, ok := x.(*bool)
	return ok && b == nil
}

func (nilBoolPtrMatcher) String() string { return "nil *bool" }

// newTestService собирает Service с моком репозитория.
// Если setup не задан — никаких вызовов репозитория не ожидается,
// и любой реальный вызов провалит тест.
func newTestService(t *testing.T, setup func(m *mocks.MockRepository)) *Service {
	t.Helper()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mocks.NewMockRepository(ctrl)
	if setup != nil {
		setup(repo)
	}
	return NewService(repo)
}

func assertError(t *testing.T, got, want error) {
	t.Helper()
	if (got == nil) != (want == nil) {
		t.Fatalf("error = %v, want %v", got, want)
	}
	if got != nil && !errors.Is(got, want) {
		t.Fatalf("error = %q, want %q", got, want)
	}
}

func TestService_GetByAuthorID(t *testing.T) {
	docs := []document.Document{
		{ID: 1, OwnerID: "author-1", Name: "Doc 1", Content: "Content 1", IsPrivate: false, CreatedAt: createdAt, EditedAt: editedAt},
		{ID: 2, OwnerID: "author-1", Name: "Doc 2", Content: "Content 2", IsPrivate: true, CreatedAt: createdAt, EditedAt: editedAt},
	}
	want := []GetDocumentOutput{
		{ID: 1, OwnerID: "author-1", Name: "Doc 1", Content: "Content 1", IsPrivate: false, CreatedAt: createdAt, EditedAt: editedAt},
		{ID: 2, OwnerID: "author-1", Name: "Doc 2", Content: "Content 2", IsPrivate: true, CreatedAt: createdAt, EditedAt: editedAt},
	}

	tests := []struct {
		name    string
		input   *GetByAuthorIDInput
		setup   func(m *mocks.MockRepository)
		want    []GetDocumentOutput
		wantErr error
	}{
		{
			name:  "success",
			input: &GetByAuthorIDInput{AuthorID: "author-1", RequesterID: "requester-1"},
			setup: func(m *mocks.MockRepository) {
				m.EXPECT().
					GetByAuthorID(gomock.Any(), "author-1", "requester-1").
					Return(docs, nil)
			},
			want: want,
		},
		{
			name:  "empty result",
			input: &GetByAuthorIDInput{AuthorID: "author-1", RequesterID: "requester-1"},
			setup: func(m *mocks.MockRepository) {
				m.EXPECT().
					GetByAuthorID(gomock.Any(), "author-1", "requester-1").
					Return(nil, nil)
			},
		},
		{
			name:  "repository failed",
			input: &GetByAuthorIDInput{AuthorID: "author-1", RequesterID: "requester-1"},
			setup: func(m *mocks.MockRepository) {
				m.EXPECT().
					GetByAuthorID(gomock.Any(), "author-1", "requester-1").
					Return(nil, errGetByAuthorFailed)
			},
			wantErr: errGetByAuthorFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(t, tt.setup)

			got, err := svc.GetByAuthorID(context.Background(), tt.input)
			assertError(t, err, tt.wantErr)
			if tt.wantErr != nil {
				if got != nil {
					t.Fatalf("got = %+v, want nil", got)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestService_GetAccessible(t *testing.T) {
	docs := []document.Document{
		{ID: 1, OwnerID: "other-1", Name: "Public doc", Content: "Content 1", IsPrivate: false, CreatedAt: createdAt, EditedAt: editedAt},
	}
	want := []GetDocumentOutput{
		{ID: 1, OwnerID: "other-1", Name: "Public doc", Content: "Content 1", IsPrivate: false, CreatedAt: createdAt, EditedAt: editedAt},
	}

	tests := []struct {
		name    string
		input   *GetAccessibleInput
		setup   func(m *mocks.MockRepository)
		want    []GetDocumentOutput
		wantErr error
	}{
		{
			name:  "success",
			input: &GetAccessibleInput{RequesterID: "requester-1"},
			setup: func(m *mocks.MockRepository) {
				m.EXPECT().
					GetAccessible(gomock.Any(), "requester-1").
					Return(docs, nil)
			},
			want: want,
		},
		{
			name:  "empty result",
			input: &GetAccessibleInput{RequesterID: "requester-1"},
			setup: func(m *mocks.MockRepository) {
				m.EXPECT().
					GetAccessible(gomock.Any(), "requester-1").
					Return(nil, nil)
			},
		},
		{
			name:  "repository failed",
			input: &GetAccessibleInput{RequesterID: "requester-1"},
			setup: func(m *mocks.MockRepository) {
				m.EXPECT().
					GetAccessible(gomock.Any(), "requester-1").
					Return(nil, errGetAccessibleFailed)
			},
			wantErr: errGetAccessibleFailed,
		},
		{
			// Даже если репозиторий вернул и документы, и ошибку,
			// сервис не мапит данные, а возвращает nil и ошибку.
			name:  "repository failed with documents",
			input: &GetAccessibleInput{RequesterID: "requester-1"},
			setup: func(m *mocks.MockRepository) {
				m.EXPECT().
					GetAccessible(gomock.Any(), "requester-1").
					Return(docs, errGetAccessibleFailed)
			},
			wantErr: errGetAccessibleFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(t, tt.setup)

			got, err := svc.GetAccessible(context.Background(), tt.input)
			assertError(t, err, tt.wantErr)
			if tt.wantErr != nil {
				if got != nil {
					t.Fatalf("got = %+v, want nil", got)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestService_GetByDocumentID(t *testing.T) {
	doc := &document.Document{ID: 7, OwnerID: "owner-1", Name: "Doc", Content: "Content", IsPrivate: true, CreatedAt: createdAt, EditedAt: editedAt}

	tests := []struct {
		name    string
		input   *GetByDocumentIDInput
		setup   func(m *mocks.MockRepository)
		want    *GetDocumentOutput
		wantErr error
	}{
		{
			name:  "success",
			input: &GetByDocumentIDInput{DocumentID: 7, RequesterID: "requester-1"},
			setup: func(m *mocks.MockRepository) {
				m.EXPECT().
					GetByID(gomock.Any(), 7, "requester-1").
					Return(doc, nil)
			},
			want: &GetDocumentOutput{ID: 7, OwnerID: "owner-1", Name: "Doc", Content: "Content", IsPrivate: true, CreatedAt: createdAt, EditedAt: editedAt},
		},
		{
			name:  "repository failed",
			input: &GetByDocumentIDInput{DocumentID: 7, RequesterID: "requester-1"},
			setup: func(m *mocks.MockRepository) {
				m.EXPECT().
					GetByID(gomock.Any(), 7, "requester-1").
					Return(nil, errGetByIDFailed)
			},
			wantErr: errGetByIDFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(t, tt.setup)

			got, err := svc.GetByDocumentID(context.Background(), tt.input)
			assertError(t, err, tt.wantErr)
			if tt.wantErr != nil {
				if got != nil {
					t.Fatalf("got = %+v, want nil", got)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestService_Create(t *testing.T) {
	created := &document.Document{ID: 5, OwnerID: "owner-1", Name: "My doc", Content: "Hello world", IsPrivate: true, CreatedAt: createdAt, EditedAt: editedAt}

	tests := []struct {
		name    string
		input   CreateInput
		setup   func(m *mocks.MockRepository)
		want    *CreateOutput
		wantErr error
	}{
		{
			name:  "success",
			input: CreateInput{OwnerID: "owner-1", Name: "My doc", Content: "Hello world", IsPrivate: true},
			setup: func(m *mocks.MockRepository) {
				m.EXPECT().
					Create(gomock.Any(), "owner-1", "My doc", "Hello world", true).
					Return(created, nil)
			},
			want: &CreateOutput{ID: 5, Name: "My doc", Content: "Hello world", IsPrivate: true, CreatedAt: createdAt},
		},
		{
			name:    "name too short",
			input:   CreateInput{OwnerID: "owner-1", Name: "ab", Content: "Hello world"},
			wantErr: document.ErrTitleTooShort,
		},
		{
			name:    "name too long",
			input:   CreateInput{OwnerID: "owner-1", Name: strings.Repeat("a", document.MAX_TITLE_LEN+1), Content: "Hello world"},
			wantErr: document.ErrTitleTooLong,
		},
		{
			name:    "content too long",
			input:   CreateInput{OwnerID: "owner-1", Name: "My doc", Content: strings.Repeat("a", document.MAX_BODY_LEN+1)},
			wantErr: document.ErrBodyTooLong,
		},
		{
			name:  "repository failed",
			input: CreateInput{OwnerID: "owner-1", Name: "My doc", Content: "Hello world", IsPrivate: true},
			setup: func(m *mocks.MockRepository) {
				m.EXPECT().
					Create(gomock.Any(), "owner-1", "My doc", "Hello world", true).
					Return(nil, errCreateFailed)
			},
			wantErr: errCreateFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(t, tt.setup)

			got, err := svc.Create(context.Background(), tt.input)
			assertError(t, err, tt.wantErr)
			if tt.wantErr != nil {
				if got != nil {
					t.Fatalf("got = %+v, want nil", got)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	deleted := &document.Document{ID: 7, OwnerID: "owner-1", Name: "Doc", Content: "Content"}

	tests := []struct {
		name    string
		input   DeleteInput
		setup   func(m *mocks.MockRepository)
		wantErr error
	}{
		{
			name:  "success",
			input: DeleteInput{DocumentID: 7, RequesterID: "requester-1"},
			setup: func(m *mocks.MockRepository) {
				m.EXPECT().
					Delete(gomock.Any(), 7, "requester-1").
					Return(deleted, nil)
			},
		},
		{
			name:  "repository failed",
			input: DeleteInput{DocumentID: 7, RequesterID: "requester-1"},
			setup: func(m *mocks.MockRepository) {
				m.EXPECT().
					Delete(gomock.Any(), 7, "requester-1").
					Return(nil, errDeleteFailed)
			},
			wantErr: errDeleteFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(t, tt.setup)

			err := svc.Delete(context.Background(), tt.input)
			assertError(t, err, tt.wantErr)
		})
	}
}

func TestService_Update(t *testing.T) {
	updated := &document.Document{ID: 7, OwnerID: "owner-1", Name: "New name", Content: "New content", IsPrivate: false, CreatedAt: createdAt, EditedAt: editedAt}
	want := &UpdateOutput{ID: 7, Name: "New name", Content: "New content", IsPrivate: false, EditedAt: editedAt}

	tests := []struct {
		name    string
		input   UpdateInput
		setup   func(m *mocks.MockRepository)
		want    *UpdateOutput
		wantErr error
	}{
		{
			name: "success: full update",
			input: UpdateInput{
				ID:        7,
				OwnerID:   "owner-1",
				Name:      strPtr("New name"),
				Content:   strPtr("New content"),
				IsPrivate: boolPtr(false),
			},
			setup: func(m *mocks.MockRepository) {
				m.EXPECT().
					Update(gomock.Any(), 7, "owner-1", gomock.Eq(strPtr("New name")), gomock.Eq(strPtr("New content")), gomock.Eq(boolPtr(false))).
					Return(updated, nil)
			},
			want: want,
		},
		{
			// Частичное обновление: непоставленные поля должны дойти до репозитория как nil.
			name: "success: partial update",
			input: UpdateInput{
				ID:      7,
				OwnerID: "owner-1",
				Name:    strPtr("New name"),
			},
			setup: func(m *mocks.MockRepository) {
				m.EXPECT().
					Update(gomock.Any(), 7, "owner-1", gomock.Eq(strPtr("New name")), nilStrPtrMatcher{}, nilBoolPtrMatcher{}).
					Return(updated, nil)
			},
			want: want,
		},
		{
			name:  "repository failed",
			input: UpdateInput{ID: 7, OwnerID: "owner-1", Name: strPtr("New name")},
			setup: func(m *mocks.MockRepository) {
				m.EXPECT().
					Update(gomock.Any(), 7, "owner-1", gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errUpdateFailed)
			},
			wantErr: errUpdateFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(t, tt.setup)

			got, err := svc.Update(context.Background(), tt.input)
			assertError(t, err, tt.wantErr)
			if tt.wantErr != nil {
				if got != nil {
					t.Fatalf("got = %+v, want nil", got)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestService_GetUnsafe(t *testing.T) {
	doc := &document.Document{ID: 7, OwnerID: "owner-1", Name: "Doc", Content: "Content", IsPrivate: true, CreatedAt: createdAt, EditedAt: editedAt}

	tests := []struct {
		name    string
		input   GetUnsafeInput
		setup   func(m *mocks.MockRepository)
		want    *GetDocumentOutput
		wantErr error
	}{
		{
			name:  "success",
			input: GetUnsafeInput{DocumentID: 7},
			setup: func(m *mocks.MockRepository) {
				m.EXPECT().
					GetUnsafe(gomock.Any(), 7).
					Return(doc, nil)
			},
			want: &GetDocumentOutput{ID: 7, OwnerID: "owner-1", Name: "Doc", Content: "Content", IsPrivate: true, CreatedAt: createdAt, EditedAt: editedAt},
		},
		{
			name:  "repository failed",
			input: GetUnsafeInput{DocumentID: 7},
			setup: func(m *mocks.MockRepository) {
				m.EXPECT().
					GetUnsafe(gomock.Any(), 7).
					Return(nil, errGetUnsafeFailed)
			},
			wantErr: errGetUnsafeFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(t, tt.setup)

			got, err := svc.GetUnsafe(context.Background(), tt.input)
			assertError(t, err, tt.wantErr)
			if tt.wantErr != nil {
				if got != nil {
					t.Fatalf("got = %+v, want nil", got)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestService_DeleteUnsafe(t *testing.T) {
	deleted := &document.Document{ID: 7, OwnerID: "owner-1", Name: "Doc", Content: "Content"}

	tests := []struct {
		name    string
		input   DeleteUnsafeInput
		setup   func(m *mocks.MockRepository)
		wantErr error
	}{
		{
			name:  "success",
			input: DeleteUnsafeInput{DocumentID: 7},
			setup: func(m *mocks.MockRepository) {
				m.EXPECT().
					DeleteUnsafe(gomock.Any(), 7).
					Return(deleted, nil)
			},
		},
		{
			name:  "repository failed",
			input: DeleteUnsafeInput{DocumentID: 7},
			setup: func(m *mocks.MockRepository) {
				m.EXPECT().
					DeleteUnsafe(gomock.Any(), 7).
					Return(nil, errDeleteUnsafeFailed)
			},
			wantErr: errDeleteUnsafeFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(t, tt.setup)

			err := svc.DeleteUnsafe(context.Background(), tt.input)
			assertError(t, err, tt.wantErr)
		})
	}
}

func TestService_UpdateUnsafe(t *testing.T) {
	updated := &document.Document{ID: 7, OwnerID: "owner-1", Name: "New name", Content: "New content", IsPrivate: true, CreatedAt: createdAt, EditedAt: editedAt}

	tests := []struct {
		name    string
		input   UpdateUnsafeInput
		setup   func(m *mocks.MockRepository)
		want    *GetDocumentOutput
		wantErr error
	}{
		{
			name: "success",
			input: UpdateUnsafeInput{
				ID:        7,
				Name:      strPtr("New name"),
				Content:   strPtr("New content"),
				IsPrivate: boolPtr(true),
			},
			setup: func(m *mocks.MockRepository) {
				m.EXPECT().
					UpdateUnsafe(gomock.Any(), 7, gomock.Eq(strPtr("New name")), gomock.Eq(strPtr("New content")), gomock.Eq(boolPtr(true))).
					Return(updated, nil)
			},
			want: &GetDocumentOutput{ID: 7, OwnerID: "owner-1", Name: "New name", Content: "New content", IsPrivate: true, CreatedAt: createdAt, EditedAt: editedAt},
		},
		{
			name:  "repository failed",
			input: UpdateUnsafeInput{ID: 7, Name: strPtr("New name")},
			setup: func(m *mocks.MockRepository) {
				m.EXPECT().
					UpdateUnsafe(gomock.Any(), 7, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errUpdateUnsafeFailed)
			},
			wantErr: errUpdateUnsafeFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(t, tt.setup)

			got, err := svc.UpdateUnsafe(context.Background(), tt.input)
			assertError(t, err, tt.wantErr)
			if tt.wantErr != nil {
				if got != nil {
					t.Fatalf("got = %+v, want nil", got)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got = %+v, want %+v", got, tt.want)
			}
		})
	}
}
