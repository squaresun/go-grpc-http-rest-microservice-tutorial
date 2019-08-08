package main

func main() {
	if err := RunServer(); err != nil {
		panic(err)
	}
}
