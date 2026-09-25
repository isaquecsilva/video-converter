package flows

type Flow interface {
	Run() (string, error)
}
