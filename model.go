package main

type Model interface {
	Init(c *Config) error
	RegisterUser(email, password string) (int, error)
	LoginUser(email, password string) (int, error)
}
