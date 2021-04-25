package main

var repos Repos

func init() {
	repos = NewRepos()
}

func main() {
	startAPIServer()
}
