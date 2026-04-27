package services

import (
	"errors"
	"testing"
	"time"

	"github.com/Nishishei01/Go_Hexagonal/internal/domains"
	"github.com/stretchr/testify/assert"
)

type mockPostRepo struct {
	CreatePostFunc  func(postRequest *domains.PostRequest) error
	GetAllPostFunc  func() ([]*domains.Post, error)
	GetPostByIDFunc func(id uint) (*domains.Post, error)
	UpdatePostFunc  func(id uint, postRequest *domains.PostRequest) error
	DeletePostFunc  func(id uint) error
}

func (m *mockPostRepo) CreatePost(postRequest *domains.PostRequest) error {
	if m.CreatePostFunc != nil {
		return m.CreatePostFunc(postRequest)
	}
	return nil
}

func (m *mockPostRepo) GetAllPost() ([]*domains.Post, error) {
	if m.GetAllPostFunc != nil {
		return m.GetAllPostFunc()
	}
	return nil, nil
}

func (m *mockPostRepo) GetPostByID(id uint) (*domains.Post, error) {
	if m.GetPostByIDFunc != nil {
		return m.GetPostByIDFunc(id)
	}
	return nil, nil
}

func (m *mockPostRepo) UpdatePost(id uint, postRequest *domains.PostRequest) error {
	if m.UpdatePostFunc != nil {
		return m.UpdatePostFunc(id, postRequest)
	}
	return nil
}

func (m *mockPostRepo) DeletePost(id uint) error {
	if m.DeletePostFunc != nil {
		return m.DeletePostFunc(id)
	}
	return nil
}

type mockPostCacheRepo struct {
	SetFunc    func(key string, value any, expiration time.Duration) error
	GetFunc    func(key string, dest any) error
	DeleteFunc func(key string) error
}

func (m *mockPostCacheRepo) Set(key string, value any, expiration time.Duration) error {
	if m.SetFunc != nil {
		return m.SetFunc(key, value, expiration)
	}
	return nil
}

func (m *mockPostCacheRepo) Get(key string, dest any) error {
	if m.GetFunc != nil {
		return m.GetFunc(key, dest)
	}
	return nil
}

func (m *mockPostCacheRepo) Delete(key string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(key)
	}
	return nil
}

func TestCreatePost(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockPostRepo{
			CreatePostFunc: func(postRequest *domains.PostRequest) error {
				return nil
			},
		}

		cacheRepo := &mockPostCacheRepo{
			DeleteFunc: func(key string) error {
				return nil
			},
		}
		service := NewPostService(repo, cacheRepo)

		err := service.Create(&domains.PostRequest{Title: "Unit Test", Content: "Unit Test", UserID: 1})
		assert.NoError(t, err)
	})
	t.Run("failure from repository", func(t *testing.T) {
		expectedErr := errors.New("database error")
		repo := &mockPostRepo{
			CreatePostFunc: func(postRequest *domains.PostRequest) error {
				return expectedErr
			},
		}
		cacheRepo := &mockPostCacheRepo{}
		service := NewPostService(repo, cacheRepo)

		err := service.Create(&domains.PostRequest{Title: "Unit Test", Content: "Unit Test", UserID: 1})
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestGetAllPost(t *testing.T) {
	t.Run("cache hit", func(t *testing.T) {
		expectedPosts := []*domains.Post{
			{
				ID:    1,
				Title: "Cached Post",
			},
		}
		cacheRepo := &mockPostCacheRepo{
			GetFunc: func(key string, dest any) error {
				if p, ok := dest.(*[]*domains.Post); ok {
					*p = expectedPosts
				}
				return nil
			},
		}
		repo := &mockPostRepo{}

		service := NewPostService(repo, cacheRepo)
		posts, err := service.GetAll()

		assert.NoError(t, err)
		assert.Equal(t, expectedPosts, posts)
	})

	t.Run("cache miss, fetch from repo success", func(t *testing.T) {
		expectedPosts := []*domains.Post{
			{ID: 1, Title: "Repo Post"},
		}
		cacheRepo := &mockPostCacheRepo{
			GetFunc: func(key string, dest any) error {
				return errors.New("cache miss")
			},
			SetFunc: func(key string, value any, expiration time.Duration) error {
				return nil
			},
		}
		repo := &mockPostRepo{
			GetAllPostFunc: func() ([]*domains.Post, error) {
				return expectedPosts, nil
			},
		}

		service := NewPostService(repo, cacheRepo)
		posts, err := service.GetAll()

		assert.NoError(t, err)
		assert.Equal(t, expectedPosts, posts)
	})

	t.Run("cache miss, fetch from repo failure", func(t *testing.T) {
		expectedErr := errors.New("db error")
		cacheRepo := &mockPostCacheRepo{
			GetFunc: func(key string, dest any) error {
				return errors.New("cache miss")
			},
		}
		repo := &mockPostRepo{
			GetAllPostFunc: func() ([]*domains.Post, error) {
				return nil, expectedErr
			},
		}

		service := NewPostService(repo, cacheRepo)
		posts, err := service.GetAll()

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, posts)
	})
}

func TestGetPostByID(t *testing.T) {
	t.Run("cache hit", func(t *testing.T) {
		expectedPost := &domains.Post{ID: 1, Title: "Cached Post"}
		cacheRepo := &mockPostCacheRepo{
			GetFunc: func(key string, dest any) error {
				if p, ok := dest.(**domains.Post); ok {
					*p = expectedPost
				}
				return nil
			},
		}
		repo := &mockPostRepo{}

		service := NewPostService(repo, cacheRepo)
		post, err := service.GetByID(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedPost, post)
	})

	t.Run("cache miss, fetch from repo success", func(t *testing.T) {
		expectedPost := &domains.Post{ID: 1, Title: "Repo Post"}
		cacheRepo := &mockPostCacheRepo{
			GetFunc: func(key string, dest any) error {
				return errors.New("cache miss")
			},
			SetFunc: func(key string, value any, expiration time.Duration) error {
				return nil
			},
		}
		repo := &mockPostRepo{
			GetPostByIDFunc: func(id uint) (*domains.Post, error) {
				return expectedPost, nil
			},
		}

		service := NewPostService(repo, cacheRepo)
		post, err := service.GetByID(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedPost, post)
	})

	t.Run("cache miss, fetch from repo failure", func(t *testing.T) {
		expectedErr := errors.New("db error")
		cacheRepo := &mockPostCacheRepo{
			GetFunc: func(key string, dest any) error {
				return errors.New("cache miss")
			},
		}
		repo := &mockPostRepo{
			GetPostByIDFunc: func(id uint) (*domains.Post, error) {
				return nil, expectedErr
			},
		}

		service := NewPostService(repo, cacheRepo)
		post, err := service.GetByID(1)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, post)
	})
}

func TestUpdatePost(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockPostRepo{
			UpdatePostFunc: func(id uint, postRequest *domains.PostRequest) error {
				return nil
			},
		}
		cacheRepo := &mockPostCacheRepo{
			DeleteFunc: func(key string) error {
				return nil
			},
		}
		service := NewPostService(repo, cacheRepo)

		err := service.Update(1, &domains.PostRequest{Title: "Updated", Content: "Updated", UserID: 1})
		assert.NoError(t, err)
	})

}

func TestDeletePost(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockPostRepo{
			DeletePostFunc: func(id uint) error {
				return nil
			},
		}
		cacheRepo := &mockPostCacheRepo{
			DeleteFunc: func(key string) error {
				return nil
			},
		}
		service := NewPostService(repo, cacheRepo)

		err := service.Delete(1)
		assert.NoError(t, err)
	})
}
