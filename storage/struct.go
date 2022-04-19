package storage

// User defined for user enity.
type User struct {
	Login    string
	FullName string
}

// Message defined for message enity.
type Message struct {
	LoginFrom string
	LoginTo   string
	CreatedAt int32
	Body      string
}
