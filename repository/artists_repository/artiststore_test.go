package artistrepository_test

import (
	"eventro2/models"
	artistrepository "eventro2/repository/artists_repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestArtistRepositoryPG_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	repo := artistrepository.NewArtistRepositoryPG(gormDB)

	artist := &models.Artist{Name: "Mock Artist"}
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "artists".*RETURNING "id"`).
		WithArgs(artist.Name, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err = repo.Create(artist)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestArtistRepositoryPG_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	repo := artistrepository.NewArtistRepositoryPG(gormDB)

	t.Run("success", func(t *testing.T) {
		expectedID := "123"
		expectedName := "Mock Artist"

		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(expectedID, expectedName)

		mock.ExpectQuery(`SELECT .* FROM "artists" WHERE id = \$1 ORDER BY "artists"."id" LIMIT \$2`).
			WithArgs(expectedID, 1).
			WillReturnRows(rows)

		artist, err := repo.GetByID(expectedID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if artist.ID != expectedID {
			t.Errorf("expected id %s, got %s", expectedID, artist.ID)
		}
		if artist.Name != expectedName {
			t.Errorf("expected name %s, got %s", expectedName, artist.Name)
		}
	})

	t.Run("not found / error", func(t *testing.T) {
		badID := "999"

		mock.ExpectQuery(`SELECT .* FROM "artists" WHERE id = \$1 ORDER BY "artists"."id" LIMIT \$2`).
			WithArgs(badID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		artist, err := repo.GetByID(badID)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if artist != nil {
			t.Errorf("expected nil artist, got %+v", artist)
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestArtistRepositoryPG_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	repo := artistrepository.NewArtistRepositoryPG(gormDB)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow("1", "Artist One").
			AddRow("2", "Artist Two")

		mock.ExpectQuery(`SELECT \* FROM "artists"`).
			WillReturnRows(rows)

		artists, err := repo.List()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(artists) != 2 {
			t.Errorf("expected 2 artists, got %d", len(artists))
		}
		if artists[0].ID != "1" || artists[0].Name != "Artist One" {
			t.Errorf("unexpected first artist: %+v", artists[0])
		}
		if artists[1].ID != "2" || artists[1].Name != "Artist Two" {
			t.Errorf("unexpected second artist: %+v", artists[1])
		}
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "artists"`).
			WillReturnError(gorm.ErrInvalidDB)

		artists, err := repo.List()
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if artists != nil {
			t.Errorf("expected nil artists, got %+v", artists)
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestArtistRepositoryPG_ListByEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	repo := artistrepository.NewArtistRepositoryPG(gormDB)
	eventID := "evt_123"

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "bio"}).
			AddRow("1", "Event Artist 1", "Bio1").
			AddRow("2", "Event Artist 2", "Bio2")

		mock.ExpectQuery(`SELECT .* FROM "artists" JOIN event_artists ea ON ea\.artist_id = artists\.id WHERE ea\.event_id = \$1`).
			WithArgs(eventID).
			WillReturnRows(rows)

		artists, err := repo.ListByEvent(eventID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(artists) != 2 {
			t.Errorf("expected 2 artists, got %d", len(artists))
		}
		if artists[0].Name != "Event Artist 1" {
			t.Errorf("unexpected first artist: %+v", artists[0])
		}
		if artists[1].Name != "Event Artist 2" {
			t.Errorf("unexpected second artist: %+v", artists[1])
		}
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .* FROM "artists" JOIN event_artists ea ON ea\.artist_id = artists\.id WHERE ea\.event_id = \$1`).
			WithArgs(eventID).
			WillReturnError(gorm.ErrInvalidDB)

		artists, err := repo.ListByEvent(eventID)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if artists != nil {
			t.Errorf("expected nil artists, got %+v", artists)
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestArtistRepositoryPG_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	repo := artistrepository.NewArtistRepositoryPG(gormDB)

	artist := &models.Artist{ID: "1", Name: "Updated Artist", Bio: "Updated Bio"}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "artists" SET .* WHERE "id" = \$3`).
		WithArgs(artist.Name, artist.Bio, artist.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = repo.Update(artist)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestArtistRepositoryPG_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	repo := artistrepository.NewArtistRepositoryPG(gormDB)

	artistID := "1"

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "artists" WHERE id = \$1`).
		WithArgs(artistID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = repo.Delete(artistID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestArtistRepositoryPG_GetEventsByArtistID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	repo := artistrepository.NewArtistRepositoryPG(gormDB)
	artistID := "123"

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow("evt_1", "Concert A").
			AddRow("evt_2", "Concert B")

		mock.ExpectQuery(`SELECT .* FROM "events" JOIN event_artists ea ON ea\.event_id = events\.id WHERE ea\.artist_id = \$1`).
			WithArgs(artistID).
			WillReturnRows(rows)

		events, err := repo.GetEventsByArtistID(artistID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(events) != 2 {
			t.Errorf("expected 2 events, got %d", len(events))
		}
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .* FROM "events" JOIN event_artists ea ON ea\.event_id = events\.id WHERE ea\.artist_id = \$1`).
			WithArgs(artistID).
			WillReturnError(gorm.ErrInvalidDB)

		events, err := repo.GetEventsByArtistID(artistID)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if events != nil {
			t.Errorf("expected nil events, got %+v", events)
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
func TestArtistRepositoryPG_SearchByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	repo := artistrepository.NewArtistRepositoryPG(gormDB)
	searchName := "Rockstar"
	likePattern := "%rockstar%"

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "bio"}).
			AddRow("1", "Rockstar", "Famous rockstar").
			AddRow("2", "Little Rockstar", "Upcoming star")

		mock.ExpectQuery(`SELECT .* FROM "artists" WHERE LOWER\(name\) LIKE \$1`).
			WithArgs(likePattern).
			WillReturnRows(rows)

		artists, err := repo.SearchByName(searchName)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(artists) != 2 {
			t.Errorf("expected 2 artists, got %d", len(artists))
		}
		if artists[0].Name != "Rockstar" {
			t.Errorf("expected first artist 'Rockstar', got %s", artists[0].Name)
		}
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .* FROM "artists" WHERE LOWER\(name\) LIKE \$1`).
			WithArgs(likePattern).
			WillReturnError(gorm.ErrInvalidDB)

		artists, err := repo.SearchByName(searchName)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if artists != nil {
			t.Errorf("expected nil artists, got %+v", artists)
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
