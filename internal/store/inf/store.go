package inf

type ErrCrud interface {
	error
	GetCode() int
}

type CRUD interface {
	Put(string, []byte) ErrCrud
	Post(string, []byte) ErrCrud
	Get(string) ([][]byte, ErrCrud)
	Delete(string) ErrCrud
	List(string) ([]string, ErrCrud)
}
