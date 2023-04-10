package store

type Store interface {
	User() UserRepo
	Account() AccountRepo
	Transaction() TransactionRepo
}
