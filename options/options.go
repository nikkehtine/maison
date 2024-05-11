package options

type Config struct {
	Input  string
	Output string
}

var DefaultConfig = Config{
	Input:  "./",
	Output: "./public",
}
