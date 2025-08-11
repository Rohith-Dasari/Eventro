package eventsrepository

type EventStorageI interface {
	//add methods
	SaveEvents()
	GetEvents()
}
