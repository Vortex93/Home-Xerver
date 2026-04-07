package arraylist

// ========================================
// Types
// ========================================
type ArrayList[T any] struct {
    data  []T
    limit int
}

// ========================================
// Constructors
// ========================================

func New[T any]() *ArrayList[T] {
    return &ArrayList[T]{
        data:  make([]T, 0),
        limit: -1, // unlimited
    }
}

func NewWithLimit[T any](limit int) *ArrayList[T] {
    return &ArrayList[T]{
        data:  make([]T, 0, limit),
        limit: limit,
    }
}

// ========================================
// Methods
// ========================================

func (al *ArrayList[T]) Add(element T) {
    // unlimited
    if al.limit <= 0 {
        al.data = append(al.data, element)
        return
    }

    // not full
    if len(al.data) < al.limit {
        al.data = append(al.data, element)
        return
    }

    // full → shift left and replace last
    copy(al.data, al.data[1:])
    al.data[len(al.data)-1] = element
}

func (al *ArrayList[T]) Get(index int) (T, bool) {
    if index < 0 || index >= len(al.data) {
        var zero T
        return zero, false
    }
    return al.data[index], true
}

func (al *ArrayList[T]) Set(index int, value T) bool {
    if index < 0 || index >= len(al.data) {
        return false
    }
    al.data[index] = value
    return true
}

func (al *ArrayList[T]) Remove(index int) (T, bool) {
    if index < 0 || index >= len(al.data) {
        var zero T
        return zero, false
    }

    removed := al.data[index]
    al.data = append(al.data[:index], al.data[index+1:]...)
    return removed, true
}

func (al *ArrayList[T]) Size() int {
    return len(al.data)
}

func (al *ArrayList[T]) IsEmpty() bool {
    return len(al.data) == 0
}

func (al *ArrayList[T]) IsFull() bool {
    if al.limit <= 0 {
        return false
    }
    return len(al.data) >= al.limit
}

func (al *ArrayList[T]) Clear() {
    if al.limit > 0 {
        al.data = make([]T, 0, al.limit)
        return
    }
    al.data = make([]T, 0)
}

func (al *ArrayList[T]) ToSlice() []T {
    return al.data
}

func (al *ArrayList[T]) Limit() int {
    return al.limit
}

func (al *ArrayList[T]) All(predicate func(T) bool) bool {
    for _, element := range al.data {
        if !predicate(element) {
            return false
        }
    }
    return true
}
