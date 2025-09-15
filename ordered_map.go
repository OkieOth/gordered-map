package omap

import (
	"encoding/json"
	"fmt"
	"iter"
	"os"

	"github.com/okieoth/gordered-map/ordered"
)

type MapThing struct {
	entry *ordered.OrderedValue
}

func GetValue[T any](m *MapThing, key string) (T, bool) {
	if m.entry.Type == ordered.OBJECT {
		if orderObject, isOk := m.entry.Value.(ordered.OrderedObject); isOk {
			for _, o := range orderObject {
				if o.Key == key {
					if tInst, isT := o.Value.Value.(T); isT {
						return tInst, true
					} else {
						break
					}
				}
			}
		}
	}
	return *(new(T)), false
}

func GetChildMap(m *MapThing, key string) (*MapThing, bool) {
	if m.entry.Type == ordered.OBJECT {
		if orderObject, isOk := m.entry.Value.(ordered.OrderedObject); isOk {
			for _, o := range orderObject {
				if o.Key == key {
					if tInst, isT := o.Value.Value.(ordered.OrderedObject); isT {
						return &MapThing{
							entry: &ordered.OrderedValue{
								Type:  ordered.OBJECT,
								Value: tInst,
							},
						}, true
					} else {
						break
					}
				}
			}
		}
	}
	return nil, false
}

func GetChildArray(m *MapThing, key string) (ordered.OrderedArray, bool) {
	if m.entry.Type == ordered.OBJECT {
		if orderObject, isOk := m.entry.Value.(ordered.OrderedObject); isOk {
			for _, o := range orderObject {
				if o.Key == key {
					if tInst, isT := o.Value.Value.(ordered.OrderedArray); isT {
						return tInst, true
					} else {
						break
					}
				}
			}
		}
	}
	return nil, false
}

func Set[T any](m *MapThing, key string, value T) error {
	if m.entry.Type == ordered.OBJECT {
		if orderObject, isOk := m.entry.Value.(ordered.OrderedObject); isOk {
			orderedValue, err := ordered.NewOrderedValue[T](value)
			if err != nil {
				return fmt.Errorf("couldn't create new value for key (%s): %v", key, err)
			}
			for i, o := range orderObject {
				if o.Key == key {
					// existing value will be replaced
					if _, isT := o.Value.Value.(T); isT {
						orderObject[i].Value = &orderedValue
						return nil
					} else {
						return fmt.Errorf("key (%s) has not the requested generic type: %T", key, value)
					}
				}
			}
			// new entry is created
			m.entry.Value = append(orderObject, ordered.OrderedPair{
				Key:   key,
				Value: &orderedValue,
			})
		}
	}
	return nil
}

func Iterate[T any](m *MapThing) iter.Seq2[string, T] {
	return func(yield func(string, T) bool) {
		if m.entry.Type == ordered.OBJECT {
			if orderObject, isOk := m.entry.Value.(ordered.OrderedObject); isOk {
				for _, o := range orderObject {
					if tInst, isT := o.Value.Value.(T); isT {
						if !yield(o.Key, tInst) {
							return
						}
					}
				}
			}
		}
	}
}

func NewFromJSON(data []byte) (*MapThing, error) {
	var orderedValue ordered.OrderedValue
	if err := json.Unmarshal(data, &orderedValue); err != nil {
		return nil, fmt.Errorf("error while unmarshalling data: %v", err)
	}
	return &MapThing{
		entry: &orderedValue,
	}, nil
}

func NewFromJSONFile(fileName string) (*MapThing, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error while reading file (%s): %v", fileName, err)
	}
	return NewFromJSON(data)
}

func GetArrayLen(a ordered.OrderedArray) int {
	return len(a)
}

func GetValueAt[T any](a ordered.OrderedArray, index int) (T, error) {
	if index < 0 || index >= len(a) {
		return *(new(T)), fmt.Errorf("index out of array dimensions")
	}
	tmp := any(*new(T))
	switch tmp.(type) {
	case int, int32, int64, float32, float64, bool, string:
		if v, isOk := a[index].Value.(T); isOk {
			return v, nil
		} else {
			return *(new(T)), fmt.Errorf("value at index doesn't have the expected type")
		}
	default:
		return *(new(T)), fmt.Errorf("try get value of an unsupported type: %T", tmp)
	}
}
