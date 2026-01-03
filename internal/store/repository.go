package store

import "sync"

type Repository[K comparable, V any] struct {
	mu   sync.Mutex
	data map[K]V
}

func NewRepository[K comparable, V any]() *Repository[K, V] {
	return &Repository[K, V]{
		data: make(map[K]V),
	}
}

func (r *Repository[K, V]) Post(k K, v V) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[k] = v

}

func (r *Repository[K, V]) Get(key K) (V, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	v, exists := r.data[key]
	return v, exists
}

func (r *Repository[K, V]) GetAll() []V {
	r.mu.Lock()
	defer r.mu.Unlock()
	res := make([]V, 0, len(r.data))
	for _, v := range r.data {
		res = append(res, v)
	}
	return res
}
func (r *Repository[K, V]) Delete(key K) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.data, key)
}

func (r *Repository[K, V]) Update(key K, updater func(*V) error) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	v, exists := r.data[key]

	if !exists {
		return false
	}
	err := updater(&v)
	if err != nil {
		return false
	}
	r.data[key] = v
	return true
}
