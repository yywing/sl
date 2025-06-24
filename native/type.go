package native

import (
	"reflect"
	"sync"

	"github.com/yywing/sl/ast"
)

type NativeSelectorType[T ast.Value] struct {
	*ast.PrimitiveType
	once   sync.Once
	types  map[string]ast.ValueType
	fields map[string][]int
}

func NewNativeSelectorType[T ast.Value](kind string) *NativeSelectorType[T] {
	r := &NativeSelectorType[T]{PrimitiveType: ast.NewPrimitiveType(kind, ast.SelectorType)}
	r.init()
	return r
}

func (t *NativeSelectorType[T]) init() {
	t.once.Do(func() {
		t.types, t.fields = t.getFields()
	})
}

func (t *NativeSelectorType[T]) getFields() (map[string]ast.ValueType, map[string][]int) {
	types := make(map[string]ast.ValueType)
	fields := make(map[string][]int)

	// 获取泛型类型 T 的反射类型
	var zero T
	tType := reflect.TypeOf(zero)

	// 如果是指针类型，获取其元素类型
	if tType.Kind() == reflect.Pointer {
		tType = tType.Elem()
	}

	// 遍历结构体的所有字段
	for i := 0; i < tType.NumField(); i++ {
		field := tType.Field(i)

		// 获取 cel tag
		celTag := field.Tag.Get("sl")
		if celTag == "" {
			continue
		}

		types[celTag] = goTypeToValueType(field.Type)
		fields[celTag] = field.Index
	}

	// 如果没有找到匹配的字段，返回 nil
	return types, fields
}

func (t *NativeSelectorType[T]) Member(name string) ast.ValueType {
	return t.types[name]
}

func (t *NativeSelectorType[T]) Get(v T, key string) (ast.Value, bool) {
	index, ok := t.fields[key]
	if !ok {
		return nil, false
	}

	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	field := value.FieldByIndex(index)

	if field.IsValid() {
		return ValueFromGo(field.Interface()), true
	}
	return nil, false
}
