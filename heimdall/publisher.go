package heimdall

type Publishable interface {
	GetMessage() string
	GetTable() string
	GetAction() string
}

type Publisher interface {
	Send(publishable Publishable)
}
