package repo

import "fmt"

type AppRepoImpl struct {
}

func (a *AppRepoImpl) Save() error {
	fmt.Println("this is my impl")
}
