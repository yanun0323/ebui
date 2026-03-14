package helper

type HashCache[T any] struct {
	nextHash  string
	valueHash string
	value     T
}

func NewHashCache[T any]() *HashCache[T] {
	return &HashCache[T]{}
}

func (c *HashCache[T]) SetNextHash(hash string) {
	c.nextHash = hash
}

func (c *HashCache[T]) IsNextHashCached() bool {
	return c.nextHash == c.valueHash
}

func (c *HashCache[T]) IsNextCacheOutdated() bool {
	return c.nextHash != c.valueHash
}

func (c *HashCache[T]) Update(value T) {
	c.valueHash = c.nextHash
	c.value = value
}

func (c *HashCache[T]) Load() T {
	return c.value
}
