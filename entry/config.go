package entry

var cfg Cfg
var name = "ossync"

type Cfg struct {
	Db    string   `json:"db"`
	Oss   string   `json:"oss"`
	Files []string `json:"files"`
}
