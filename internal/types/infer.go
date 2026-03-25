package types

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"structify-go/internal/naming"
)

type Kind int

const (
	StringKind Kind = iota
	BoolKind
	IntKind
	Float64Kind
	InterfaceKind
	StructKind
	SliceKind
)

type ValueType struct {
	Kind      Kind
	StructDef *StructDef
	Element   *ValueType
	Signature string
}

type Field struct {
	Name     string
	JSONName string
	Type     *ValueType
}

type StructDef struct {
	Name      string
	Fields    []Field
	Signature string
}

type Inferrer struct {
	definitions     map[string]*StructDef
	usedStructNames map[string]int
}

func NewInferrer() *Inferrer {
	return &Inferrer{
		definitions:     make(map[string]*StructDef),
		usedStructNames: make(map[string]int),
	}
}

func (i *Inferrer) InferRoot(rootName string, value any) (*StructDef, error) {
	object, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unsupported root JSON type %T: root value must be an object", value)
	}

	return i.inferObject(rootName, object)
}

func (i *Inferrer) inferValue(suggestedName string, value any) (*ValueType, error) {
	switch typed := value.(type) {
	case nil:
		return primitiveType(InterfaceKind, "interface{}"), nil
	case string:
		return primitiveType(StringKind, "string"), nil
	case bool:
		return primitiveType(BoolKind, "bool"), nil
	case json.Number:
		if looksFloat(typed.String()) {
			return primitiveType(Float64Kind, "float64"), nil
		}

		return primitiveType(IntKind, "int"), nil
	case float64:
		return primitiveType(Float64Kind, "float64"), nil
	case map[string]any:
		definition, err := i.inferObject(suggestedName, typed)
		if err != nil {
			return nil, err
		}

		return &ValueType{
			Kind:      StructKind,
			StructDef: definition,
			Signature: definition.Signature,
		}, nil
	case []any:
		return i.inferArray(suggestedName, typed)
	default:
		return nil, fmt.Errorf("unsupported JSON value of type %T", value)
	}
}

func (i *Inferrer) inferObject(suggestedName string, object map[string]any) (*StructDef, error) {
	keys := make([]string, 0, len(object))
	for key := range object {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	fields := make([]Field, 0, len(keys))
	signatureParts := make([]string, 0, len(keys))
	usedFieldNames := make(map[string]int)

	for _, key := range keys {
		fieldType, err := i.inferValue(key, object[key])
		if err != nil {
			return nil, fmt.Errorf("field %q: %w", key, err)
		}

		fieldName := naming.UniqueName(naming.ToFieldName(key), usedFieldNames)
		fields = append(fields, Field{
			Name:     fieldName,
			JSONName: key,
			Type:     fieldType,
		})

		signatureParts = append(signatureParts, key+":"+fieldType.Signature)
	}

	signature := "struct{" + strings.Join(signatureParts, ";") + "}"
	if existing, ok := i.definitions[signature]; ok {
		return existing, nil
	}

	name := naming.UniqueName(naming.ToTypeName(suggestedName), i.usedStructNames)
	definition := &StructDef{
		Name:      name,
		Fields:    fields,
		Signature: signature,
	}

	i.definitions[signature] = definition
	return definition, nil
}

func (i *Inferrer) inferArray(suggestedName string, values []any) (*ValueType, error) {
	if len(values) == 0 {
		return &ValueType{
			Kind:      SliceKind,
			Element:   primitiveType(InterfaceKind, "interface{}"),
			Signature: "[]interface{}",
		}, nil
	}

	elementName := naming.Singularize(suggestedName)
	elementType, err := i.inferValue(elementName, values[0])
	if err != nil {
		return nil, fmt.Errorf("array index 0: %w", err)
	}

	for index := 1; index < len(values); index++ {
		nextType, err := i.inferValue(elementName, values[index])
		if err != nil {
			return nil, fmt.Errorf("array index %d: %w", index, err)
		}

		if nextType.Signature != elementType.Signature {
			return nil, fmt.Errorf(
				"unsupported mixed array element types at index %d: %s and %s",
				index,
				describeType(elementType),
				describeType(nextType),
			)
		}
	}

	return &ValueType{
		Kind:      SliceKind,
		Element:   elementType,
		Signature: "[]" + elementType.Signature,
	}, nil
}

func primitiveType(kind Kind, signature string) *ValueType {
	return &ValueType{
		Kind:      kind,
		Signature: signature,
	}
}

func looksFloat(value string) bool {
	return strings.ContainsAny(value, ".eE")
}

func describeType(valueType *ValueType) string {
	switch valueType.Kind {
	case StringKind:
		return "string"
	case BoolKind:
		return "bool"
	case IntKind:
		return "int"
	case Float64Kind:
		return "float64"
	case InterfaceKind:
		return "interface{}"
	case StructKind:
		return "object"
	case SliceKind:
		return "array"
	default:
		return "unknown"
	}
}
