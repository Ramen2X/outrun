package requests

type SendApolloRequest struct {
	Base
	Type  int64    `json:"type,string"`
	Value []string `json:"value"`
}

type SetNoahIDRequest struct {
	Base
	NoahID int64 `json:"noahId,string"`
}
