package service

import "sync"

// RatingStore is an interface to store laptop ratings
type RatingStore interface {
	Add(laptopId string, score float64) (*Rating, error)
}

// Rating contains the rating information for a laptop
type Rating struct {
	Count uint32
	Sum   float64
}

// InMemoryRatingStore stores laptop ratings in memory
type InMemoryRatingStore struct {
	mutex   sync.Mutex
	ratings map[string]*Rating
}

func NewInMemoryRatingStore() *InMemoryRatingStore {
	return &InMemoryRatingStore{
		ratings: make(map[string]*Rating),
	}
}

func (store *InMemoryRatingStore) Add(laptopId string, score float64) (*Rating, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	rating := store.ratings[laptopId]
	if rating == nil {
		rating = &Rating{
			Count: 1,
			Sum:   score,
		}
	} else {
		rating.Count += 1
		rating.Sum += score
	}

	store.ratings[laptopId] = rating

	return rating, nil
}
