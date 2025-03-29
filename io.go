package main

type BaseIOSchema interface {
	ToJson() (string, error)
}
