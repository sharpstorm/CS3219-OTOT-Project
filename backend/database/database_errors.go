package database

type UniqueKeyError struct{}

func (m *UniqueKeyError) Error() string {
	return "Unique key violation"
}
