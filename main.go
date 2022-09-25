package main

import "github.com/JECSand/go-rest-api-boilerplate/cmd"

func main() {
	var app cmd.App
	err := app.Initialize()
	if err != nil {
		panic(err)
	}
	app.Run()
}
