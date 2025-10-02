package eventrepository_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"eventro2/models"
	eventrepository "eventro2/repository/event_repository"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTest(t *testing.T) (*eventrepository.EventRepositoryPG, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock db: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	repo := eventrepository.NewEventRepositoryPG(gormDB)
	cleanup := func() { db.Close() }

	return repo, mock, cleanup
}

func TestEventRepositoryPG_Create(t *testing.T) {
	repo, mock, cleanup := setupTest(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		ev := &models.Event{ID: "evt_1", Name: "Sample Event"}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "events"`).
			// For Postgres, GORM uses RETURNING "id"
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(ev.ID))
		mock.ExpectCommit()

		if err := repo.Create(ev); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("db error", func(t *testing.T) {
		ev := &models.Event{ID: "evt_2", Name: "Broken Insert"}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "events"`).
			WillReturnError(fmt.Errorf("insert failed"))
		mock.ExpectRollback()

		if err := repo.Create(ev); err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestEventRepositoryPG_GetByID(t *testing.T) {
	repo, mock, cleanup := setupTest(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		eventID := "evt_123"

		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(eventID, "Music Festival")

		mock.ExpectQuery(`SELECT .* FROM "events"`).
			WithArgs(eventID, 1). // gorm always adds LIMIT 1
			WillReturnRows(rows)

		ev, err := repo.GetByID(eventID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if ev.ID != eventID {
			t.Errorf("expected id %s, got %s", eventID, ev.ID)
		}
	})

	t.Run("error - not found", func(t *testing.T) {
		eventID := "evt_not_found"

		mock.ExpectQuery(`SELECT .* FROM "events"`).
			WithArgs(eventID, 1).
			WillReturnError(errors.New("record not found"))

		_, err := repo.GetByID(eventID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestEventRepositoryPG_List(t *testing.T) {
	repo, mock, cleanup := setupTest(t)
	defer cleanup()

	t.Run("success - returns events", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "category"}).
			AddRow("evt1", "Music Festival", "Music").
			AddRow("evt2", "Food Carnival", "Food")

		mock.ExpectQuery(`SELECT .* FROM "events"`).
			WillReturnRows(rows)

		events, err := repo.List()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(events) != 2 {
			t.Errorf("expected 2 events, got %d", len(events))
		}
		if events[0].Name != "Music Festival" {
			t.Errorf("expected first event 'Music Festival', got %s", events[0].Name)
		}
	})

	t.Run("error - db failure", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .* FROM "events"`).
			WillReturnError(errors.New("db failure"))

		_, err := repo.List()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestEventRepositoryPG_Update(t *testing.T) {
	repo, mock, cleanup := setupTest(t)
	defer cleanup()

	t.Run("success - updates event is_blocked", func(t *testing.T) {
		event := &models.Event{ID: "evt1", IsBlocked: true}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "events" SET "is_blocked"=`).
			WithArgs(event.IsBlocked, event.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(event)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("error - db failure", func(t *testing.T) {
		event := &models.Event{ID: "evt2", IsBlocked: false}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "events" SET "is_blocked"=`).
			WithArgs(event.IsBlocked, event.ID).
			WillReturnError(errors.New("update failed"))
		mock.ExpectRollback()

		err := repo.Update(event)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestEventRepositoryPG_Delete(t *testing.T) {
	repo, mock, cleanup := setupTest(t)
	defer cleanup()

	t.Run("success - deletes event", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "events" WHERE id = \$1`).
			WithArgs("evt1").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete("evt1")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("error - db failure", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "events" WHERE id = \$1`).
			WithArgs("evt2").
			WillReturnError(errors.New("delete failed"))
		mock.ExpectRollback()

		err := repo.Delete("evt2")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestEventRepositoryPG_AddEventArtist(t *testing.T) {
	repo, mock, cleanup := setupTest(t)
	defer cleanup()

	ea := &models.EventArtist{
		EventID:  "evt1",
		ArtistID: "art1",
	}

	t.Run("success - insert event artist", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "event_artists"`).
			WithArgs(ea.EventID, ea.ArtistID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.AddEventArtist(ea)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("error - db insert failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "event_artists"`).
			WithArgs(ea.EventID, ea.ArtistID).
			WillReturnError(errors.New("insert failed"))
		mock.ExpectRollback()

		err := repo.AddEventArtist(ea)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestEventRepositoryPG_GetArtistsByEventID(t *testing.T) {
	repo, mock, cleanup := setupTest(t)
	defer cleanup()

	eventID := "evt1"

	t.Run("success - returns artists", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow("art1", "Artist One").
			AddRow("art2", "Artist Two")

		mock.ExpectQuery(`SELECT (.+) FROM "artists"`).
			WithArgs(eventID).
			WillReturnRows(rows)

		artists, err := repo.GetArtistsByEventID(eventID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(artists) != 2 {
			t.Errorf("expected 2 artists, got %d", len(artists))
		}
	})

	t.Run("success - no artists found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"}) // empty

		mock.ExpectQuery(`SELECT (.+) FROM "artists"`).
			WithArgs(eventID).
			WillReturnRows(rows)

		artists, err := repo.GetArtistsByEventID(eventID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(artists) != 0 {
			t.Errorf("expected 0 artists, got %d", len(artists))
		}
	})

	t.Run("error - query fails", func(t *testing.T) {
		mock.ExpectQuery(`SELECT (.+) FROM "artists"`).
			WithArgs(eventID).
			WillReturnError(errors.New("db error"))

		_, err := repo.GetArtistsByEventID(eventID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestEventRepositoryPG_GetEventsByCity(t *testing.T) {
	repo, mock, cleanup := setupTest(t)
	defer cleanup()

	city := "New York"
	lowerCity := strings.ToLower(city)
	queryRegex := `SELECT DISTINCT e\.\* FROM events e JOIN shows s ON s\.event_id = e\.id JOIN venues v ON s\.venue_id = v\.id WHERE LOWER\(v\.city\) = \$1`

	t.Run("success - events found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "category", "is_blocked"}).
			AddRow("evt1", "Concert Night", "Music", false).
			AddRow("evt2", "Comedy Show", "Comedy", false)

		mock.ExpectQuery(queryRegex).
			WithArgs(lowerCity).
			WillReturnRows(rows)

		events, err := repo.GetEventsByCity(city)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(events) != 2 {
			t.Errorf("expected 2 events, got %d", len(events))
		}
	})

	t.Run("success - no events found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "category", "is_blocked"})

		mock.ExpectQuery(queryRegex).
			WithArgs(lowerCity).
			WillReturnRows(rows)

		events, err := repo.GetEventsByCity(city)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(events) != 0 {
			t.Errorf("expected 0 events, got %d", len(events))
		}
	})

	t.Run("error - query failure", func(t *testing.T) {
		mock.ExpectQuery(queryRegex).
			WithArgs(lowerCity).
			WillReturnError(errors.New("db error"))

		_, err := repo.GetEventsByCity(city)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestEventRepositoryPG_GetFilteredEvents(t *testing.T) {
	repo, mock, cleanup := setupTest(t)
	defer cleanup()

	baseCols := []string{"id", "name", "category", "is_blocked"}

	testCases := []struct {
		name      string
		filter    models.EventFilter
		setupMock func()
		expectLen int
		expectErr bool
	}{
		{
			name:   "filter by name",
			filter: models.EventFilter{Name: "Rock"},
			setupMock: func() {
				rows := sqlmock.NewRows(baseCols).AddRow("evt1", "Rock Night", "Music", false)
				mock.ExpectQuery(`SELECT .* FROM "events" WHERE LOWER\(name\) LIKE \$1`).
					WithArgs("%rock%").
					WillReturnRows(rows)
			},
			expectLen: 1,
			expectErr: false,
		},
		{
			name:   "filter by category",
			filter: models.EventFilter{Category: "Comedy"},
			setupMock: func() {
				rows := sqlmock.NewRows(baseCols).AddRow("evt2", "Comedy Night", "Comedy", false)
				mock.ExpectQuery(`SELECT .* FROM "events" WHERE category = \$1`).
					WithArgs("Comedy").
					WillReturnRows(rows)
			},
			expectLen: 1,
			expectErr: false,
		},
		{
			name:   "filter by location",
			filter: models.EventFilter{Location: "New York"},
			setupMock: func() {
				rows := sqlmock.NewRows(baseCols).AddRow("evt3", "Jazz Fest", "Music", false)
				mock.ExpectQuery(`SELECT .* FROM "events" JOIN shows s ON s.event_id = events.id JOIN venues v ON s.venue_id = v.id WHERE LOWER\(v.city\) = \$1`).
					WithArgs("new york").
					WillReturnRows(rows)
			},
			expectLen: 1,
			expectErr: false,
		},
		{
			name: "filter by is_blocked",
			filter: func() models.EventFilter {
				blocked := true
				return models.EventFilter{IsBlocked: &blocked}
			}(),
			setupMock: func() {
				rows := sqlmock.NewRows(baseCols).AddRow("evt4", "Hidden Show", "Drama", true)
				mock.ExpectQuery(`SELECT .* FROM "events" WHERE is_blocked = \$1`).
					WithArgs(true).
					WillReturnRows(rows)
			},
			expectLen: 1,
			expectErr: false,
		},
		{
			name:   "filter by artist name",
			filter: models.EventFilter{ArtistName: "Adele"},
			setupMock: func() {
				rows := sqlmock.NewRows(baseCols).AddRow("evt5", "Pop Concert", "Music", false)
				mock.ExpectQuery(`SELECT .* FROM "events" JOIN event_artists ea ON ea.event_id = events.id JOIN artists a ON ea.artist_id = a.id WHERE LOWER\(a.name\) LIKE \$1`).
					WithArgs("%adele%").
					WillReturnRows(rows)
			},
			expectLen: 1,
			expectErr: false,
		},
		{
			name:   "filter by eventID",
			filter: models.EventFilter{EventID: "evt6"},
			setupMock: func() {
				rows := sqlmock.NewRows(baseCols).AddRow("evt6", "Single Event", "Special", false)
				mock.ExpectQuery(`SELECT .* FROM "events" WHERE id = \$1`).
					WithArgs("evt6").
					WillReturnRows(rows)
			},
			expectLen: 1,
			expectErr: false,
		},
		{
			name:   "no filters - return all",
			filter: models.EventFilter{},
			setupMock: func() {
				rows := sqlmock.NewRows(baseCols).
					AddRow("evt7", "Mega Fest", "Festival", false).
					AddRow("evt8", "Standup Show", "Comedy", false)
				mock.ExpectQuery(`SELECT .* FROM "events"`).
					WillReturnRows(rows)
			},
			expectLen: 2,
			expectErr: false,
		},
		{
			name:   "db error",
			filter: models.EventFilter{},
			setupMock: func() {
				mock.ExpectQuery(`SELECT .* FROM "events"`).
					WillReturnError(errors.New("db error"))
			},
			expectLen: 0,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			events, err := repo.GetFilteredEvents(tc.filter)

			if tc.expectErr && err == nil {
				t.Error("expected error, got nil")
			}

			if !tc.expectErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if !tc.expectErr && len(events) != tc.expectLen {
				t.Errorf("expected %d events, got %d", tc.expectLen, len(events))
			}
		})
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
