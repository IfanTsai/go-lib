package promise

type TypeState string

const (
	TypeStatePending   TypeState = "pending"
	TypeStateFulfilled TypeState = "fulfilled"
	TypeStateRejected  TypeState = "rejected"
)
