package plugins

// a plugin is a set of an arbitrary map and reduce function

type Plugin interface {
	Map() [][]string
	Reduce() string
}