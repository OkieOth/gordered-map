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

func HasValue(m *MapThing, key string) bool {
	if m.entry.Type == ordered.OBJECT {
		if orderObject, isOk := m.entry.Value.(ordered.OrderedObject); isOk {
			for _, o := range orderObject {
				if o.Key == key {
					return true
				}
			}
		}
	}
	return false
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

func GetTypedChildArray[T any](m *MapThing, key string) ([]T, bool) {
	if m.entry.Type == ordered.OBJECT {
		if orderObject, isOk := m.entry.Value.(ordered.OrderedObject); isOk {
			for _, o := range orderObject {
				if o.Key == key {
					if tInst, isT := o.Value.Value.(ordered.OrderedArray); isT {
						ret := make([]T, 0)
						for _, orderedValue := range tInst {
							if v, ok := orderedValue.Value.(T); ok {
								ret = append(ret, v)
							}
						}
						return ret, true
					} else {
						break
					}
				}
			}
		}
	}
	return []T{}, false
}

func GetAnyTypedChildArray(m *MapThing, key string) ([]any, bool) {
	if m.entry.Type == ordered.OBJECT {
		if orderObject, isOk := m.entry.Value.(ordered.OrderedObject); isOk {
			for _, o := range orderObject {
				if o.Key == key {
					if tInst, isT := o.Value.Value.(ordered.OrderedArray); isT {
						ret := make([]any, 0)
						for _, orderedValue := range tInst {
							ret = append(ret, orderedValue.Value)
						}
						return ret, true
					} else {
						break
					}
				}
			}
		}
	}
	return []any{}, false
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

type CastFunc func(value ordered.OrderedPair) (any, bool)

func OrderedPair2Value(value ordered.OrderedPair) (any, bool) {
	return value.Value.Value, true
}

type CastFunc2[T any] func(value ordered.OrderedValue) (T, bool)

func OrderedValue2Value[T any](value ordered.OrderedValue) (T, bool) {
	if s, ok := value.Value.(T); ok {
		return s, true
	}
	return *(new(T)), false
}

func (m *MapThing) Iterate() iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		if m.entry.Type == ordered.OBJECT {
			if orderObject, isOk := m.entry.Value.(ordered.OrderedObject); isOk {
				for _, o := range orderObject {
					if !yield(o.Key, o.Value) {
						return
					}
				}
			}
		}
	}
}

// helps to interate over sub dictionaries
func (m *MapThing) IterateOverMaps() iter.Seq2[string, MapThing] {
	return func(yield func(string, MapThing) bool) {
		if m.entry.Type == ordered.OBJECT {
			if orderObject, isOk := m.entry.Value.(ordered.OrderedObject); isOk {
				for _, o := range orderObject {
					mapThing := MapThing{
						entry: o.Value,
					}
					if !yield(o.Key, mapThing) {
						return
					}
				}
			}
		}
	}
}

func ToText(v any) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func (m *MapThing) IterateToValue(castFunc CastFunc) iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		if m.entry.Type == ordered.OBJECT {
			if orderObject, isOk := m.entry.Value.(ordered.OrderedObject); isOk {
				for _, o := range orderObject {
					if item, shouldReturned := castFunc(o); shouldReturned {
						if !yield(o.Key, item) {
							return
						}
					}
				}
			}
		}
	}
}

func IterateOverArray[T any](array any, castFunc CastFunc2[T]) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		if arrayType, ok := array.(ordered.OrderedArray); ok {
			for i, o := range arrayType {
				if item, shouldReturned := castFunc(*o); shouldReturned {
					if !yield(i, item) {
						return
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

func (m *MapThing) SerializeJSONFile(fileName string) error {
	outputFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error while creating output file (%s): %v", fileName, err)
	}
	defer outputFile.Close()
	jsonData, err := m.Serialize()
	if err != nil {
		return err
	}
	_, err = outputFile.Write(jsonData)
	if err != nil {
		return fmt.Errorf("error while writing output file (%s): %v", fileName, err)
	}

	return nil
}

func (m *MapThing) Serialize() ([]byte, error) {
	jsonData, err := json.MarshalIndent(m.entry, "", "  ")
	if err != nil {
		return []byte{}, fmt.Errorf("error while marshal data: %v", err)
	}
	return jsonData, nil
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
