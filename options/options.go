/*
Copyright © 2024 nikkehtine <nikkehtine@int.pl>
*/
package options

type Config struct {
	Input  string
	Output string
}

var DefaultConfig = Config{
	Input:  "./",
	Output: "./public",
}
