package bookingrepository_test

import (
	"eventro2/models"
	bookingrepository "eventro2/repository/booking_repository"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestBookingRepositoryPG_Create(t *testing.T) {
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

	repo := bookingrepository.NewBookingRepositoryPG(gormDB)

	booking := &models.Booking{
		BookingID: "b123",
		UserID:    "u1",
		ShowID:    "s1",
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "bookings"`).
			WithArgs(
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnRows(sqlmock.NewRows([]string{"booking_id"}).AddRow(booking.BookingID))
		mock.ExpectCommit()

		err := repo.Create(booking)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "bookings"`).
			WithArgs(
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnError(fmt.Errorf("insert failed"))
		mock.ExpectRollback()

		err := repo.Create(booking)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
func TestBookingRepositoryPG_GetByID(t *testing.T) {
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

	repo := bookingrepository.NewBookingRepositoryPG(gormDB)
	bookingID := "b123"

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"booking_id", "user_id", "show_id"}).
			AddRow("b123", "u1", "s1")

		mock.ExpectQuery(`SELECT .* FROM "bookings" WHERE booking_id = \$1 ORDER BY "bookings"."booking_id" LIMIT \$2`).
			WithArgs(bookingID, 1).
			WillReturnRows(rows)

		booking, err := repo.GetByID(bookingID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if booking.BookingID != "b123" {
			t.Errorf("expected booking_id b123, got %s", booking.BookingID)
		}
	})

	t.Run("not found / error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .* FROM "bookings" WHERE booking_id = \$1 ORDER BY "bookings"."booking_id" LIMIT \$2`).
			WithArgs(bookingID, 1).
			WillReturnError(fmt.Errorf("record not found"))

		_, err := repo.GetByID(bookingID)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
func TestBookingRepositoryPG_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock db: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	repo := bookingrepository.NewBookingRepositoryPG(gormDB)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"booking_id", "user_id", "show_id"}).
			AddRow("b1", "u1", "s1").
			AddRow("b2", "u2", "s2")

		mock.ExpectQuery(`SELECT \* FROM "bookings"`).
			WillReturnRows(rows)

		bookings, err := repo.List()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(bookings) != 2 {
			t.Errorf("expected 2 bookings, got %d", len(bookings))
		}
		if bookings[0].BookingID != "b1" {
			t.Errorf("expected first booking b1, got %s", bookings[0].BookingID)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "bookings"`).
			WillReturnError(fmt.Errorf("db error"))

		_, err := repo.List()
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
func TestBookingRepositoryPG_ListByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock db: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	repo := bookingrepository.NewBookingRepositoryPG(gormDB)
	userID := "u1"

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"booking_id", "user_id", "show_id"}).
			AddRow("b1", "u1", "s1").
			AddRow("b2", "u1", "s2")

		mock.ExpectQuery(`SELECT \* FROM "bookings" WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnRows(rows)

		bookings, err := repo.ListByUser(userID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(bookings) != 2 {
			t.Errorf("expected 2 bookings, got %d", len(bookings))
		}
		if bookings[0].BookingID != "b1" {
			t.Errorf("expected first booking b1, got %s", bookings[0].BookingID)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "bookings" WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnError(fmt.Errorf("db error"))

		_, err := repo.ListByUser(userID)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestBookingRepositoryPG_ListByShow(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock db: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	repo := bookingrepository.NewBookingRepositoryPG(gormDB)
	showID := "s1"

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"booking_id", "user_id", "show_id"}).
			AddRow("b1", "u1", "s1").
			AddRow("b2", "u2", "s1")

		mock.ExpectQuery(`SELECT \* FROM "bookings" WHERE show_id = \$1`).
			WithArgs(showID).
			WillReturnRows(rows)

		bookings, err := repo.ListByShow(showID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(bookings) != 2 {
			t.Errorf("expected 2 bookings, got %d", len(bookings))
		}
		if bookings[0].BookingID != "b1" {
			t.Errorf("expected first booking b1, got %s", bookings[0].BookingID)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "bookings" WHERE show_id = \$1`).
			WithArgs(showID).
			WillReturnError(fmt.Errorf("db error"))

		_, err := repo.ListByShow(showID)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
func TestBookingRepositoryPG_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock db: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	repo := bookingrepository.NewBookingRepositoryPG(gormDB)
	booking := &models.Booking{BookingID: "b1", UserID: "u1", ShowID: "s1"}

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "bookings"`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), booking.BookingID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(booking)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "bookings"`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), booking.BookingID).
			WillReturnError(fmt.Errorf("db error"))
		mock.ExpectRollback()

		err := repo.Update(booking)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
func TestBookingRepositoryPG_Find(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	repo := bookingrepository.NewBookingRepositoryPG(gormDB)

	booking := models.Booking{BookingID: "b1", UserID: "u1", ShowID: "s1"}

	t.Run("filter by bookingID", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"booking_id", "user_id", "show_id"}).
			AddRow(booking.BookingID, booking.UserID, booking.ShowID)

		mock.ExpectQuery(`SELECT .* FROM "bookings"`).
			WithArgs("b1").
			WillReturnRows(rows)

		results, err := repo.Find("b1", "", "")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(results) != 1 {
			t.Errorf("expected 1 result, got %d", len(results))
		}
	})

	t.Run("filter by userID", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"booking_id", "user_id", "show_id"}).
			AddRow(booking.BookingID, booking.UserID, booking.ShowID)

		mock.ExpectQuery(`SELECT .* FROM "bookings"`).
			WithArgs("u1").
			WillReturnRows(rows)

		results, err := repo.Find("", "u1", "")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(results) != 1 {
			t.Errorf("expected 1 result, got %d", len(results))
		}
	})

	t.Run("filter by showID", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"booking_id", "user_id", "show_id"}).
			AddRow(booking.BookingID, booking.UserID, booking.ShowID)

		mock.ExpectQuery(`SELECT .* FROM "bookings"`).
			WithArgs("s1").
			WillReturnRows(rows)

		results, err := repo.Find("", "", "s1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(results) != 1 {
			t.Errorf("expected 1 result, got %d", len(results))
		}
	})

	t.Run("filter by all fields", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"booking_id", "user_id", "show_id"}).
			AddRow(booking.BookingID, booking.UserID, booking.ShowID)

		mock.ExpectQuery(`SELECT .* FROM "bookings"`).
			WithArgs("b1", "u1", "s1").
			WillReturnRows(rows)

		results, err := repo.Find("b1", "u1", "s1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(results) != 1 {
			t.Errorf("expected 1 result, got %d", len(results))
		}
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .* FROM "bookings"`).
			WithArgs("b1").
			WillReturnError(fmt.Errorf("db error"))

		_, err := repo.Find("b1", "", "")
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
func TestBookingRepositoryPG_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock db: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	repo := bookingrepository.NewBookingRepositoryPG(gormDB)
	bookingID := "b1"

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "bookings" WHERE booking_id = \$1`).
			WithArgs(bookingID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(bookingID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "bookings" WHERE booking_id = \$1`).
			WithArgs(bookingID).
			WillReturnError(fmt.Errorf("delete failed"))
		mock.ExpectRollback()

		err := repo.Delete(bookingID)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
