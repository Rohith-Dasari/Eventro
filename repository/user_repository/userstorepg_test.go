package userrepository_test

import (
	"errors"
	"eventro2/models"
	userrepository "eventro2/repository/user_repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	cleanup := func() { db.Close() }
	return gdb, mock, cleanup
}

func TestUserRepositoryPG(t *testing.T) {
	baseCols := []string{"user_id", "email", "password", "username", "is_blocked"}

	// create
	t.Run("Create", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := userrepository.NewUserRepositoryPG(gdb)

		user := &models.User{UserID: "u1", Email: "test@example.com", Username: "Test User"}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "users"`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow("u1"))
		mock.ExpectCommit()

		if err := repo.Create(user); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("GetByID success", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := userrepository.NewUserRepositoryPG(gdb)

		rows := sqlmock.NewRows(baseCols).AddRow("u1", "test@example.com", "password", "Test User", false)
		mock.ExpectQuery(`SELECT \* FROM "users" WHERE user_id = \$1 ORDER BY "users"."user_id" LIMIT \$2`).
			WithArgs("u1", 1).WillReturnRows(rows)

		user, err := repo.GetByID("u1")
		if err != nil || user.UserID != "u1" {
			t.Errorf("expected user u1, got %v, err=%v", user, err)
		}
	})

	t.Run("GetByID error", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := userrepository.NewUserRepositoryPG(gdb)

		mock.ExpectQuery(`SELECT \* FROM "users" WHERE user_id = \$1 ORDER BY "users"."user_id" LIMIT \$2`).
			WithArgs("u1", 1).WillReturnError(errors.New("db error"))

		_, err := repo.GetByID("u1")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("GetByEmail success", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := userrepository.NewUserRepositoryPG(gdb)

		rows := sqlmock.NewRows(baseCols).AddRow("u1", "test@example.com", "password", "Test User", false)
		mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 ORDER BY "users"."user_id" LIMIT \$2`).
			WithArgs("test@example.com", 1).WillReturnRows(rows)

		user, err := repo.GetByEmail("test@example.com")
		if err != nil || user.Email != "test@example.com" {
			t.Errorf("expected user with email test@example.com, got %v, err=%v", user, err)
		}
	})

	t.Run("GetByEmail error", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := userrepository.NewUserRepositoryPG(gdb)

		mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 ORDER BY "users"."user_id" LIMIT \$2`).
			WithArgs("test@example.com", 1).WillReturnError(errors.New("db error"))

		_, err := repo.GetByEmail("test@example.com")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("GetUsers", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := userrepository.NewUserRepositoryPG(gdb)

		rows := sqlmock.NewRows(baseCols).
			AddRow("u1", "user1@example.com", "password1", "User One", false).
			AddRow("u2", "user2@example.com", "password2", "User Two", false)

		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WillReturnRows(rows)

		users, err := repo.GetUsers()
		if err != nil || len(users) != 2 {
			t.Errorf("expected 2 users, got %v, err=%v", users, err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := userrepository.NewUserRepositoryPG(gdb)

		user := &models.User{UserID: "u1", Email: "updated@example.com", Username: "Updated User"}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users" SET`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		if err := repo.Update(user); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := userrepository.NewUserRepositoryPG(gdb)

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "users" WHERE user_id = \$1`).
			WithArgs("u1").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		if err := repo.Delete("u1"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("GetBlockedUsers", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := userrepository.NewUserRepositoryPG(gdb)

		rows := sqlmock.NewRows(baseCols).
			AddRow("u3", "blocked@example.com", "password3", "Blocked User", true)

		mock.ExpectQuery(`SELECT \* FROM "users" WHERE is_blocked = \$1`).
			WithArgs(true).
			WillReturnRows(rows)

		users, err := repo.GetBlockedUsers()
		if err != nil || len(users) != 1 {
			t.Errorf("expected 1 blocked user, got %v, err=%v", users, err)
		}
	})
}
