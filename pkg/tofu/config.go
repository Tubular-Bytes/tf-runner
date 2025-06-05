package tofu

var (
	debug bool = false
)

func SetDebug(d bool) {
	debug = d
}

func Debug() bool {
	return debug
}
