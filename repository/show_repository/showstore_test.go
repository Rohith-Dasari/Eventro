package showrepository_test

import (
	"errors"
	"eventro2/models"
	showrepository "eventro2/repository/show_repository"
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

func TestShowRepositoryPG(t *testing.T) {
	baseCols := []string{"id", "event_id", "host_id", "venue_id"}

	// create
	t.Run("Create", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := showrepository.NewShowRepositoryPG(gdb)

		show := &models.Show{ID: "s1", EventID: "e1", HostID: "h1", VenueID: "v1"}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "shows"`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("s1"))
		mock.ExpectCommit()

		if err := repo.Create(show); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	// get by id
	t.Run("GetByID success", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := showrepository.NewShowRepositoryPG(gdb)

		rows := sqlmock.NewRows(baseCols).AddRow("s1", "e1", "h1", "v1")
		// Match exactly what GORM produces
		mock.ExpectQuery(`SELECT \* FROM "shows" WHERE id = \$1 ORDER BY "shows"."id" LIMIT \$2`).
			WithArgs("s1", 1).WillReturnRows(rows)

		show, err := repo.GetByID("s1")
		if err != nil || show.ID != "s1" {
			t.Errorf("expected show s1, got %v, err=%v", show, err)
		}
	})

	t.Run("GetByID error", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := showrepository.NewShowRepositoryPG(gdb)

		mock.ExpectQuery(`SELECT \* FROM "shows" WHERE id = \$1 ORDER BY "shows"."id" LIMIT \$2`).
			WithArgs("s1", 1).WillReturnError(errors.New("db error"))

		_, err := repo.GetByID("s1")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	// list
	t.Run("List", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := showrepository.NewShowRepositoryPG(gdb)

		rows := sqlmock.NewRows(baseCols).
			AddRow("s1", "e1", "h1", "v1").
			AddRow("s2", "e2", "h2", "v2")

		mock.ExpectQuery(`SELECT \* FROM "shows"`).
			WillReturnRows(rows)

		shows, err := repo.List()
		if err != nil || len(shows) != 2 {
			t.Errorf("expected 2 shows, got %v, err=%v", shows, err)
		}
	})

	t.Run("ListByEvent", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := showrepository.NewShowRepositoryPG(gdb)

		rows := sqlmock.NewRows(baseCols).AddRow("s1", "e1", "h1", "v1")
		mock.ExpectQuery(`SELECT \* FROM "shows" WHERE event_id = \$1`).
			WithArgs("e1").WillReturnRows(rows)

		shows, err := repo.ListByEvent("e1")
		if err != nil || len(shows) != 1 {
			t.Errorf("expected 1 show, got %v, err=%v", shows, err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := showrepository.NewShowRepositoryPG(gdb)

		show := &models.Show{ID: "s1", EventID: "e1", HostID: "h1", VenueID: "v1"}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "shows" SET`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		if err := repo.Update(show); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	// delete
	t.Run("Delete", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := showrepository.NewShowRepositoryPG(gdb)

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "shows" WHERE id = \$1`).
			WithArgs("s1").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		if err := repo.Delete("s1"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	// find
	t.Run("Find with filters", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := showrepository.NewShowRepositoryPG(gdb)

		rows := sqlmock.NewRows(baseCols).AddRow("s1", "e1", "h1", "v1")

		mock.ExpectQuery(`SELECT \* FROM "shows" WHERE id = \$1 AND event_id = \$2 AND host_id = \$3 AND venue_id = \$4`).
			WithArgs("s1", "e1", "h1", "v1").
			WillReturnRows(rows)

		filter := models.ShowFilter{ShowID: "s1", EventID: "e1", HostID: "h1", VenueID: "v1"}
		shows, err := repo.Find(filter)
		if err != nil || len(shows) != 1 {
			t.Errorf("expected 1 show, got %v, err=%v", shows, err)
		}
	})

	t.Run("Find error", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := showrepository.NewShowRepositoryPG(gdb)

		mock.ExpectQuery(`SELECT \* FROM "shows"`).
			WillReturnError(errors.New("db error"))

		_, err := repo.Find(models.ShowFilter{})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
