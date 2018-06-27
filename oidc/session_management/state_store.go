package session_management

type StateKey string
type StateValue string


//go:generate counterfeiter . Store
type Store interface {
	Put(StateKey, StateValue) error
	Get(StateKey) StateValue
}

