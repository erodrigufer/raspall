package server

import (
	"fmt"

	"github.com/patrickmn/go-cache"
)

type Deliverable interface {
	CreateHash() (string, error)
}

func limit[O any](limit int, objects []O) []O {
	if limit < 1 {
		return objects
	}

	if len(objects) >= limit {
		return objects[:limit]
	}

	return objects
}

func checkDeliveryStatus(element Deliverable, c *cache.Cache) (bool, error) {
	hashId, err := element.CreateHash()
	if err != nil {
		return false, fmt.Errorf("unable to create a hash id: %w", err)
	}

	_, found := c.Get(hashId)
	if found {
		return true, nil
	}

	c.Set(hashId, true, cache.DefaultExpiration)
	return false, nil

}

func getUndeliveredObjects[O Deliverable](objects []O, c *cache.Cache) ([]O, error) {
	output := make([]O, 0, 30)

	for _, object := range objects {
		delivered, err := checkDeliveryStatus(object, c)
		if err != nil {
			return nil, fmt.Errorf("unable to get undelivered objects: %w", err)
		}
		if !delivered {
			output = append(output, object)
		}
	}

	return output, nil
}
