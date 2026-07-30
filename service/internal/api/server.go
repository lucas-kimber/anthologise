package api

type server struct {
	store Store
}

func newServer(store Store) *server {
	return &server{
		store,
	}
}
