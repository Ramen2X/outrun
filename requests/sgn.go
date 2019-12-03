package requests

type SendApolloRequest struct {
	Base
	Type  int64    `json:"type"`
	Value []string `json:"value"`
}
