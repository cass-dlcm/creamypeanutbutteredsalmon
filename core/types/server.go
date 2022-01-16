package types

/*
Server is a representation of an instance of a server of stat.ink or salmon-stats/api.
*/
type Server struct {
	Address   string `json:"address"`
	APIKey    string `json:"api_key,omitempty"`
	ShortName string `json:"short_name"`
}
