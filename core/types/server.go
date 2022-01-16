package types

/*
Server is a representation of an instance of a server of stat.ink or salmon-stats/api.
*/
type Server struct {
	Address   string `json:"address"`
	APIKey    string `json:"api_key"`
	ShortName string `json:"short_name"`
}
