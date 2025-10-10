package usecase

import (
	"context"
	"errors"
	"shortener/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Post(ctx context.Context, link model.Link) (*model.Link, error) {
	args := m.Called(ctx, link)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Link), args.Error(1)
}

func (m *MockRepository) Get(ctx context.Context, shortenUrl string) (string, error) {
	args := m.Called(ctx, shortenUrl)
	return args.String(0), args.Error(1)
}

func (m *MockRepository) FilterLinks(ctx context.Context, f model.FilterLinksInput) ([]model.Link, error) {
	args := m.Called(ctx, f)
	return args.Get(0).([]model.Link), args.Error(1)
}

func (m *MockRepository) IncrementVisits(ctx context.Context, shortenUrl string) error {
	args := m.Called(ctx, shortenUrl)
	return args.Error(0)
}

func TestCreateLink(t *testing.T) {
	testCases := []struct {
		name           string
		input          model.CreateLinkInput
		repo           func(*MockRepository)
		expectedResult *model.Link
		expectedError  error
	}{
		{
			name: "return success result after 1 attempt",
			input: model.CreateLinkInput{
				OriginalUrl: "https://google.com",
			},
			repo: func(mr *MockRepository) {
				mr.On("Post", mock.Anything, mock.MatchedBy(func(link model.Link) bool {
					return link.OriginalUrl == "https://google.com" && len(link.ShortenUrl) == ShortenLength
				})).Return(&model.Link{
					OriginalUrl: "https://google.com",
					ShortenUrl:  "abc123de",
				}, nil).Once()
			},
			expectedResult: &model.Link{
				OriginalUrl: "https://google.com",
				ShortenUrl:  "abc123de",
			},
			expectedError: nil,
		},
		{
			name: "return success result after 2 attempts (collision)",
			input: model.CreateLinkInput{
				OriginalUrl: "https://google.com",
			},
			repo: func(mr *MockRepository) {
				mr.On("Post", mock.Anything, mock.MatchedBy(func(link model.Link) bool {
					return link.OriginalUrl == "https://google.com"
				})).Return(nil, model.ErrorNonUniq).Once()

				mr.On("Post", mock.Anything, mock.MatchedBy(func(link model.Link) bool {
					return link.OriginalUrl == "https://google.com"
				})).Return(&model.Link{
					OriginalUrl: "https://google.com",
					ShortenUrl:  "xyz789ab",
				}, nil).Once()
			},
			expectedResult: &model.Link{
				OriginalUrl: "https://google.com",
				ShortenUrl:  "xyz789ab",
			},
			expectedError: nil,
		},
		{
			name: "3 failed attempts due to collisions",
			input: model.CreateLinkInput{
				OriginalUrl: "https://google.com",
			},
			repo: func(mr *MockRepository) {
				mr.On("Post", mock.Anything, mock.MatchedBy(func(link model.Link) bool {
					return link.OriginalUrl == "https://google.com"
				})).Return(nil, model.ErrorNonUniq).Times(3)
			},
			expectedResult: nil,
			expectedError:  model.ErrorAttemptsExhausted,
		},
		{
			name: "return other repository error",
			input: model.CreateLinkInput{
				OriginalUrl: "https://google.com",
			},
			repo: func(mr *MockRepository) {
				// Ошибка не связанная с коллизией
				mr.On("Post", mock.Anything, mock.MatchedBy(func(link model.Link) bool {
					return link.OriginalUrl == "https://google.com"
				})).Return(nil, errors.New("database connection failed")).Once()
			},
			expectedResult: nil,
			expectedError:  errors.New("database connection failed"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			if tc.repo != nil {
				tc.repo(mockRepo)
			}
			uc := New(mockRepo)

			result, err := uc.CreateLink(context.Background(), tc.input)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedResult, result)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestRedirect(t *testing.T) {
	testCases := []struct {
		name           string
		shortenUrl     string
		setupMocks     func(*MockRepository)
		expectedResult string
		expectedError  error
	}{
		{
			name:       "success redirect",
			shortenUrl: "abc123de",
			setupMocks: func(mr *MockRepository) {
				mr.On("Get", mock.Anything, "abc123de").Return("https://google.com", nil).Once()
				mr.On("IncrementVisits", mock.Anything, "abc123de").Return(nil).Once()
			},
			expectedResult: "https://google.com",
			expectedError:  nil,
		},
		{
			name:       "return ErrorNotFound",
			shortenUrl: "nonexistent",
			setupMocks: func(mr *MockRepository) {
				mr.On("Get", mock.Anything, "nonexistent").Return("", model.ErrorNotFound).Once()
			},
			expectedResult: "",
			expectedError:  model.ErrorNotFound,
		},
		{
			name:       "return OK if ErrorIncrementVisits only",
			shortenUrl: "abc123de",
			setupMocks: func(mr *MockRepository) {
				mr.On("Get", mock.Anything, "abc123de").Return("https://google.com", nil).Once()
				mr.On("IncrementVisits", mock.Anything, "abc123de").Return(model.ErrorIncrementVisits).Once()
			},
			expectedResult: "https://google.com",
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockRepository)

			if tc.setupMocks != nil {
				tc.setupMocks(mockRepo)
			}

			uc := New(mockRepo)
			result, err := uc.Redirect(context.Background(), tc.shortenUrl)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedResult, result)
			if tc.expectedError != nil {
				mockRepo.AssertNotCalled(t, "IncrementVisits")
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
