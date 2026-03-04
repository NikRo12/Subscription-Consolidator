package store

type Store interface {
	User() UserRepository
	Sub() SubRepository
}
