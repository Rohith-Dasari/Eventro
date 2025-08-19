package models

type EventCategory string

const (
	Movie    EventCategory = "movie"
	Sports   EventCategory = "sports"
	Concert  EventCategory = "concert"
	Workshop EventCategory = "workshop"
	Party    EventCategory = "party"
)

type Event struct {
	ID          string        `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string        `gorm:"type:text;not null"`
	Description string        `gorm:"type:text"`
	HypeMeter   int           `gorm:"default:0"`
	Duration    string        `gorm:"type:text"`
	Category    EventCategory `gorm:"type:text;not null"`
	IsBlocked   bool          `gorm:"default:false"`
}

type EventArtist struct {
	EventID  string `gorm:"primaryKey;type:uuid"`
	ArtistID string `gorm:"primaryKey;type:uuid"`

	Event  Event  `gorm:"foreignKey:EventID;references:ID;constraint:OnDelete:CASCADE"`
	Artist Artist `gorm:"foreignKey:ArtistID;references:ID;constraint:OnDelete:CASCADE"`
}
