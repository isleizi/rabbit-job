package app

type AppEntity struct {
	appName int64
	appKey  string
}

func Save(a *AppEntity) error {
	return nil
}
