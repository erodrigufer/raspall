package server

import "time"

func (app *Application) checkIfResponseCached(sourceName string) bool {
	_, found := app.cache.Get(sourceName)
	return found
}

func (app *Application) cacheResponse(sourceName string) {
	app.cache.Set(sourceName, true, 5*time.Minute)
}
