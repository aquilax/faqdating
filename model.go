package main

type Model interface {
	Init(c *Config) error
}
