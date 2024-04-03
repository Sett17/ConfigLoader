package fieldsetter

import (
	"fmt"
	"testing"
)

type TestObject struct {
	StringField string
	IntField    int
	FloatField  float64
	BoolField   bool
	NestedField struct {
		StringField string
	}
	ArrayField   []string
	MapField     map[string]string
	PointerField *string
}

func TestSetValue(t *testing.T) {
	testObject := &TestObject{
		ArrayField: make([]string, 1),
		MapField:   make(map[string]string),
	}

	otherString := "other"
	_ = otherString

	type args struct {
		obj       any
		path      string
		value     any
		validator func(obj any) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Set string field",
			args: args{
				obj:   testObject,
				path:  "StringField",
				value: "test",
				validator: func(obj any) error {
					if obj.(*TestObject).StringField != "test" {
						return fmt.Errorf("expected StringField to be 'test', got %s", obj.(*TestObject).StringField)
					}
					return nil
				},
			},
		},
		{
			name: "Set int field",
			args: args{
				obj:   testObject,
				path:  "IntField",
				value: 42,
				validator: func(obj any) error {
					if obj.(*TestObject).IntField != 42 {
						return fmt.Errorf("expected IntField to be 42, got %d", obj.(*TestObject).IntField)
					}
					return nil
				},
			},
		},
		{
			name: "Set float field",
			args: args{
				obj:   testObject,
				path:  "FloatField",
				value: 3.14,
				validator: func(obj any) error {
					if obj.(*TestObject).FloatField != 3.14 {
						return fmt.Errorf("expected FloatField to be 3.14, got %f", obj.(*TestObject).FloatField)
					}
					return nil
				},
			},
		},
		{
			name: "Set bool field",
			args: args{
				obj:   testObject,
				path:  "BoolField",
				value: true,
				validator: func(obj any) error {
					if obj.(*TestObject).BoolField != true {
						return fmt.Errorf("expected BoolField to be true, got %t", obj.(*TestObject).BoolField)
					}
					return nil
				},
			},
		},
		{
			name: "Set nested field",
			args: args{
				obj:   testObject,
				path:  "NestedField.StringField",
				value: "nested",
				validator: func(obj any) error {
					if obj.(*TestObject).NestedField.StringField != "nested" {
						return fmt.Errorf("expected NestedField.StringField to be 'nested', got %s", obj.(*TestObject).NestedField.StringField)
					}
					return nil
				},
			},
		},
		{
			name: "Set array field",
			args: args{
				obj:   testObject,
				path:  "ArrayField.0",
				value: "first",
				validator: func(obj any) error {
					if obj.(*TestObject).ArrayField[0] != "first" {
						return fmt.Errorf("expected ArrayField[0] to be 'first', got %s", obj.(*TestObject).ArrayField[0])
					}
					return nil
				},
			},
		},
		{
			name: "Set array field out of bounds",
			args: args{
				obj:   testObject,
				path:  "ArrayField.1",
				value: "second",
			},
			wantErr: true,
		},
		{
			name: "Set map field",
			args: args{
				obj:   testObject,
				path:  "MapField.key",
				value: "value",
				validator: func(obj any) error {
					if obj.(*TestObject).MapField["key"] != "value" {
						return fmt.Errorf("expected MapField[\"key\"] to be 'value', got %s", obj.(*TestObject).MapField["key"])
					}
					return nil
				},
			},
		},
		{
			name: "Set map field with invalid value",
			args: args{
				obj:   testObject,
				path:  "MapField.key",
				value: 42,
			},
			wantErr: true,
		},
		{
			name: "Set pointer field",
			args: args{
				obj:   testObject,
				path:  "PointerField",
				value: &otherString,
				validator: func(obj any) error {
					if *obj.(*TestObject).PointerField != "other" {
						return fmt.Errorf("expected PointerField to be 'other', got %s", *obj.(*TestObject).PointerField)
					}
					return nil
				},
			},
		},
		{
			name: "Set pointer field with nil value",
			args: args{
				obj:   testObject,
				path:  "PointerField",
				value: nil,
				validator: func(obj any) error {
					if obj.(*TestObject).PointerField != nil {
						return fmt.Errorf("expected PointerField to be nil, got %v", obj.(*TestObject).PointerField)
					}
					return nil
				},
			},
		},
		{
			name: "Set pointer field with invalid value",
			args: args{
				obj:   testObject,
				path:  "PointerField",
				value: 42,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SetValue(tt.args.obj, tt.args.path, tt.args.value)
			if (result != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", result, tt.wantErr)
			} else if tt.args.validator != nil && tt.args.validator(tt.args.obj) != nil {
				t.Errorf("SetValue() failed validation: %v", tt.args.validator(tt.args.obj))
			}
		})
	}
}
