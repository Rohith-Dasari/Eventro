package venuerepository_test

import (
	"errors"
	"eventro2/models"
	venuerepository "eventro2/repository/venue_repository"
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

func TestVenueRepositoryPG(t *testing.T) {
	baseCols := []string{"id", "host_id", "name", "city", "is_blocked"}

	// ---- Create ----
	t.Run("Create", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := venuerepository.NewVenueRepositoryPG(gdb)

		venue := &models.Venue{ID: "v1", HostID: "h1", Name: "Test Venue", City: "Test City"}

		mock.ExpectBegin()
		// Adjusted to 7 arguments to match the actual query
		mock.ExpectQuery(`INSERT INTO "venues"`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("v1"))
		mock.ExpectCommit()

		if err := repo.Create(venue); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	// ---- GetByID ----
	t.Run("GetByID success", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := venuerepository.NewVenueRepositoryPG(gdb)

		rows := sqlmock.NewRows(baseCols).AddRow("v1", "h1", "Test Venue", "Test City", false)
		mock.ExpectQuery(`SELECT \* FROM "venues" WHERE id = \$1 ORDER BY "venues"."id" LIMIT \$2`).
			WithArgs("v1", 1).WillReturnRows(rows)

		venue, err := repo.GetByID("v1")
		if err != nil || venue.ID != "v1" {
			t.Errorf("expected venue v1, got %v, err=%v", venue, err)
		}
	})

	t.Run("GetByID error", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := venuerepository.NewVenueRepositoryPG(gdb)

		mock.ExpectQuery(`SELECT \* FROM "venues" WHERE id = \$1 ORDER BY "venues"."id" LIMIT \$2`).
			WithArgs("v1", 1).WillReturnError(errors.New("db error"))

		_, err := repo.GetByID("v1")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	// ---- List ----
	t.Run("List", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := venuerepository.NewVenueRepositoryPG(gdb)

		rows := sqlmock.NewRows(baseCols).
			AddRow("v1", "h1", "Venue One", "City One", false).
			AddRow("v2", "h2", "Venue Two", "City Two", false)

		mock.ExpectQuery(`SELECT \* FROM "venues"`).
			WillReturnRows(rows)

		venues, err := repo.List()
		if err != nil || len(venues) != 2 {
			t.Errorf("expected 2 venues, got %v, err=%v", venues, err)
		}
	})

	// ---- ListByHost ----
	t.Run("ListByHost", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := venuerepository.NewVenueRepositoryPG(gdb)

		rows := sqlmock.NewRows(baseCols).
			AddRow("v1", "h1", "Venue One", "City One", false).
			AddRow("v2", "h1", "Venue Two", "City Two", false)

		mock.ExpectQuery(`SELECT \* FROM "venues" WHERE host_id = \$1`).
			WithArgs("h1").WillReturnRows(rows)

		venues, err := repo.ListByHost("h1")
		if err != nil || len(venues) != 2 {
			t.Errorf("expected 2 venues, got %v, err=%v", venues, err)
		}
	})

	// ---- ListByCity ----
	t.Run("ListByCity", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := venuerepository.NewVenueRepositoryPG(gdb)

		rows := sqlmock.NewRows(baseCols).
			AddRow("v1", "h1", "Venue One", "City One", false).
			AddRow("v2", "h2", "Venue Two", "City One", false)

		mock.ExpectQuery(`SELECT \* FROM "venues" WHERE city ILIKE \$1`).
			WithArgs("City One").WillReturnRows(rows)

		venues, err := repo.ListByCity("City One")
		if err != nil || len(venues) != 2 {
			t.Errorf("expected 2 venues, got %v, err=%v", venues, err)
		}
	})

	// ---- Update ----
	t.Run("Update", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := venuerepository.NewVenueRepositoryPG(gdb)

		venue := &models.Venue{ID: "v1", HostID: "h1", Name: "Updated Venue", City: "Updated City"}

		mock.ExpectBegin()
		// Adjusted to 7 arguments to match the actual query
		mock.ExpectExec(`UPDATE "venues" SET`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		if err := repo.Update(venue); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	// ---- Delete ----
	t.Run("Delete", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := venuerepository.NewVenueRepositoryPG(gdb)

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "venues" WHERE id = \$1`).
			WithArgs("v1").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		if err := repo.Delete("v1"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	// ---- Find with filters ----
	t.Run("Find with VenueID filter", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := venuerepository.NewVenueRepositoryPG(gdb)

		rows := sqlmock.NewRows(baseCols).
			AddRow("v1", "h1", "Venue One", "City One", false)

		mock.ExpectQuery(`SELECT \* FROM "venues" WHERE id = \$1`).
			WithArgs("v1").WillReturnRows(rows)

		filter := models.VenueFilter{VenueID: "v1"}
		venues, err := repo.Find(filter)
		if err != nil || len(venues) != 1 || venues[0].ID != "v1" {
			t.Errorf("expected venue v1, got %v, err=%v", venues, err)
		}
	})

	t.Run("Find with HostID filter", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := venuerepository.NewVenueRepositoryPG(gdb)

		rows := sqlmock.NewRows(baseCols).
			AddRow("v1", "h1", "Venue One", "City One", false).
			AddRow("v2", "h1", "Venue Two", "City Two", false)

		mock.ExpectQuery(`SELECT \* FROM "venues" WHERE host_id = \$1`).
			WithArgs("h1").WillReturnRows(rows)

		filter := models.VenueFilter{HostID: "h1"}
		venues, err := repo.Find(filter)
		if err != nil || len(venues) != 2 {
			t.Errorf("expected 2 venues, got %v, err=%v", venues, err)
		}
	})

	t.Run("Find with City filter", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := venuerepository.NewVenueRepositoryPG(gdb)

		rows := sqlmock.NewRows(baseCols).
			AddRow("v1", "h1", "Venue One", "City One", false).
			AddRow("v2", "h2", "Venue Two", "City One", false)

		mock.ExpectQuery(`SELECT \* FROM "venues" WHERE city = \$1`).
			WithArgs("City One").WillReturnRows(rows)

		filter := models.VenueFilter{City: "City One"}
		venues, err := repo.Find(filter)
		if err != nil || len(venues) != 2 {
			t.Errorf("expected 2 venues, got %v, err=%v", venues, err)
		}
	})

	t.Run("Find with IsBlocked filter", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := venuerepository.NewVenueRepositoryPG(gdb)

		rows := sqlmock.NewRows(baseCols).
			AddRow("v3", "h3", "Blocked Venue", "City Three", true)

		mock.ExpectQuery(`SELECT \* FROM "venues" WHERE is_blocked = \$1`).
			WithArgs(true).WillReturnRows(rows)

		filter := models.VenueFilter{IsBlocked: true}
		venues, err := repo.Find(filter)
		if err != nil || len(venues) != 1 || !venues[0].IsBlocked {
			t.Errorf("expected 1 blocked venue, got %v, err=%v", venues, err)
		}
	})

	t.Run("Find error", func(t *testing.T) {
		gdb, mock, cleanup := setupTestDB(t)
		defer cleanup()
		repo := venuerepository.NewVenueRepositoryPG(gdb)

		mock.ExpectQuery(`SELECT \* FROM "venues"`).
			WillReturnError(errors.New("db error"))

		_, err := repo.Find(models.VenueFilter{})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
