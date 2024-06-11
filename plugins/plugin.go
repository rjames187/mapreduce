package plugins

// a plugin is a set of an arbitrary map and reduce function

type Plugin interface {
	Map(string) [][]string
	Reduce(string, []string) []string
}

var Plugins map[string]Plugin = map[string]Plugin{
	"wc": WordCount{},
}