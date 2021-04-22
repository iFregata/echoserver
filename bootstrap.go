package main

var repos Repos

func init() {
	repos = CreateRepos()
}

func main() {
	startAPIServer()
}
