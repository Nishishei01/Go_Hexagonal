package gorm

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Nishishei01/Go_Hexagonal/internal/domains"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("Failed to open gorm db: %v", err)
	}

	return gormDB, mock, func() {
		mockDB.Close()
	}
}

func TestPostGormRepository_CreatePost(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewPostGormRepository(db)
	req := &domains.PostRequest{Title: "Title", Content: "Content", UserID: 1}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO "posts"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		err := repo.CreatePost(req)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure", func(t *testing.T) {
		expectedErr := errors.New("db insert error")
		mock.ExpectQuery(`INSERT INTO "posts"`).
			WillReturnError(expectedErr)

		err := repo.CreatePost(req)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestPostGormRepository_GetAllPost(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewPostGormRepository(db)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "content", "user_id", "create_at", "update_at"}).
			AddRow(1, "Title 1", "Content 1", 1, time.Now(), time.Now()).
			AddRow(2, "Title 2", "Content 2", 2, time.Now(), time.Now())

		mock.ExpectQuery(`SELECT \* FROM "posts"`).
			WillReturnRows(rows)

		posts, err := repo.GetAllPost()
		assert.NoError(t, err)
		assert.Len(t, posts, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure", func(t *testing.T) {
		expectedErr := errors.New("db query error")
		mock.ExpectQuery(`SELECT \* FROM "posts"`).
			WillReturnError(expectedErr)

		posts, err := repo.GetAllPost()
		assert.Error(t, err)
		assert.Nil(t, posts)
		assert.Equal(t, expectedErr, err)
	})
}

func TestPostGormRepository_GetPostByID(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewPostGormRepository(db)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "content", "user_id", "create_at", "update_at"}).
			AddRow(1, "Title 1", "Content 1", 1, time.Now(), time.Now())

		mock.ExpectQuery(`SELECT \* FROM "posts" WHERE id = \$1`).
			WithArgs(1, 1).
			WillReturnRows(rows)

		post, err := repo.GetPostByID(1)
		assert.NoError(t, err)
		if post != nil {
			assert.Equal(t, uint(1), post.ID)
		}
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "posts" WHERE id = \$1`).
			WithArgs(1, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		post, err := repo.GetPostByID(1)
		assert.Error(t, err)
		assert.Nil(t, post)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestPostGormRepository_UpdatePost(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewPostGormRepository(db)
	req := &domains.PostRequest{Title: "Updated", Content: "Updated Content", UserID: 1}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(`UPDATE "posts" SET`).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdatePost(1, req)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure not found", func(t *testing.T) {
		mock.ExpectExec(`UPDATE "posts" SET`).
			WillReturnResult(sqlmock.NewResult(1, 0))

		err := repo.UpdatePost(1, req)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestPostGormRepository_DeletePost(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewPostGormRepository(db)

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "posts" WHERE id = \$1`).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.DeletePost(1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "posts" WHERE id = \$1`).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 0))
		mock.ExpectRollback()

		err := repo.DeletePost(1)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("failure db error", func(t *testing.T) {
		expectedErr := errors.New("delete error")
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "posts" WHERE id = \$1`).
			WithArgs(1).
			WillReturnError(expectedErr)
		mock.ExpectRollback()

		err := repo.DeletePost(1)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
