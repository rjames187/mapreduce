package plugins

type KeyValue struct {
	Key   string
	Value string
}

// a pair of arbitrary map and reduce functions
type Plugin interface {
	Map(string) []*KeyValue
	Reduce([]*KeyValue) []*KeyValue
}

var Plugins map[string]Plugin = map[string]Plugin{
	"wc": WordCount{},
}