package gorm

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Nishishei01/Go_Hexagonal/internal/domains"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
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

func TestRoleGormRepository_CreateRole(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRoleGormRepository(db)
	req := &domains.RoleRequest{Role: "Unit Test"}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO "roles"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		err := repo.CreateRole(req)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure", func(t *testing.T) {
		expectedErr := errors.New("db query error")
		mock.ExpectQuery(`INSERT INTO "roles"`).
			WillReturnError(expectedErr)

		err := repo.CreateRole(req)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestRoleGormRepository_GetAllRole(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRoleGormRepository(db)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "role", "create_at", "update_at"}).
			AddRow(1, "Unit Test", time.Now(), time.Now()).
			AddRow(2, "Unit Test2", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles"`)).
			WillReturnRows(rows)

		userRolesMock := sqlmock.NewRows([]string{"role_id", "user_id"})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_roles"`)).
			WillReturnRows(userRolesMock)

		roles, err := repo.GetAllRole()
		assert.NoError(t, err)
		assert.Len(t, roles, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure", func(t *testing.T) {
		expectedErr := errors.New("db query error")
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles"`)).
			WillReturnError(expectedErr)

		roles, err := repo.GetAllRole()
		assert.Error(t, err)
		assert.Nil(t, roles)
		assert.Equal(t, expectedErr, err)
	})
}

func TestRoleGormRepository_GetByIDRole(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRoleGormRepository(db)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "role", "create_at", "update_at"}).
			AddRow(1, "Unit Test", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles" WHERE id = $1`)).
			WithArgs(1, 1).
			WillReturnRows(rows)

		userRolesMock := sqlmock.NewRows([]string{"role_id", "user_id"})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_roles"`)).
			WillReturnRows(userRolesMock)

		role, err := repo.GetRoleByID(1)
		assert.NoError(t, err)
		if role != nil {
			assert.Equal(t, uint(1), role.ID)
		}
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure", func(t *testing.T) {
		expectedErr := errors.New("db query error")
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles" WHERE id = $1`)).
			WillReturnError(expectedErr)

		roles, err := repo.GetRoleByID(1)
		assert.Error(t, err)
		assert.Nil(t, roles)
		assert.Equal(t, expectedErr, err)
	})
}

func TestRoleGormRepository_UpdateRole(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewRoleGormRepository(db)
	req := &domains.RoleRequest{Role: "Unit Test"}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "roles" SET`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateRole(1, req)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure not found", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "roles" SET`)).
			WillReturnResult(sqlmock.NewResult(1, 0))

		err := repo.UpdateRole(1, req)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

}

func TestRoleGormRepository_DeleteRole(t *testing.T) {
	db, mock, clearup := setupTestDB(t)
	defer clearup()

	repo := NewRoleGormRepository(db)

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "user_roles" WHERE "user_roles"."role_id" = $1`)).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "roles" WHERE id = $1`)).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := repo.DeleteRole(1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure not found", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "user_roles" WHERE "user_roles"."role_id" = $1`)).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 0))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "roles" WHERE id = $1`)).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 0))

		mock.ExpectRollback()

		err := repo.DeleteRole(1)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("failure db error", func(t *testing.T) {
		expectedErr := errors.New("delete error")
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "user_roles" WHERE "user_roles"."role_id" = $1`)).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 0))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "roles" WHERE id = $1`)).
			WithArgs(1).
			WillReturnError(expectedErr)

		mock.ExpectRollback()

		err := repo.DeleteRole(1)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
