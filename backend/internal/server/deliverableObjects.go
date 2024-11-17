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

	return false, nil
}

func setDeliveredStatus(element Deliverable, c *cache.Cache) (bool, error) {
	delivered, err := checkDeliveryStatus(element, c)
	if err != nil {
		return false, fmt.Errorf("unable to set delivered status: %w", err)
	}
	hashId, err := element.CreateHash()
	if err != nil {
		return false, fmt.Errorf("unable to create a hash id: %w", err)
	}

	if !delivered {
		c.Set(hashId, true, cache.DefaultExpiration)
		return false, nil
	}
	return true, nil
}

func checkIfUndeliveredObjectsPresent[O Deliverable](objects []O, c *cache.Cache) (bool, error) {
	for _, object := range objects {
		delivered, err := checkDeliveryStatus(object, c)
		if err != nil {
			return false, fmt.Errorf("unable to check delivery status of object: %w", err)
		}
		if !delivered {
			return true, nil
		}
	}
	return false, nil
}

func getUndeliveredObjects[O Deliverable](objects []O, c *cache.Cache) ([]O, error) {
	output := make([]O, 0, 30)

	for _, object := range objects {
		delivered, err := setDeliveredStatus(object, c)
		if err != nil {
			return nil, fmt.Errorf("unable to get undelivered objects: %w", err)
		}
		if !delivered {
			output = append(output, object)
		}
	}

	return output, nil
}
