package main

func main() {
	app := App{}
	app.Initialize(DBUser, DBPass, DBName)
	app.Run(":10000")
}
