package plugins

type KeyValue struct {
	Key   string
	Value string
}

// a plugin is a set of an arbitrary map and reduce function

type Plugin interface {
	Map(string) []*KeyValue
	Reduce([]*KeyValue) []*KeyValue
}

var Plugins map[string]Plugin = map[string]Plugin{
	"wc": WordCount{},
}