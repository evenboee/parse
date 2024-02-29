package parse

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type StringUnmarshaler interface {
	UnmarshalString(string) error
}

func TryWith[T any](config *Config, s string) (T, error) {
	var t T
	err := config.set(&t, s)
	return t, err
}

func Try[T any](s string, opts ...Option) (T, error) {
	conf := NewConfig(opts...)
	return TryWith[T](conf, s)
}

func MustWith[T any](config *Config, s string) T {
	v, err := TryWith[T](config, s)
	if err != nil {
		panic(err)
	}
	return v
}

func Must[T any](s string, opts ...Option) T {
	conf := NewConfig(opts...)
	return MustWith[T](conf, s)
}

func (config *Config) set(obj any, v string) error {
	val := reflect.ValueOf(obj).Elem()
	err, _ := config.setValue(val, v)
	return err
}

func (config *Config) setValue(value reflect.Value, v string) (error, bool) {
	switch i := value.Addr().Interface().(type) {
	case StringUnmarshaler:
		err := i.UnmarshalString(v)
		return err, err == nil
	}

	switch value.Kind() {
	case reflect.Pointer:
		nVal := reflect.New(value.Type().Elem())

		err, ok := config.setValue(nVal.Elem(), v)
		if err != nil {
			return err, false
		}
		if ok {
			value.Set(nVal)
			return nil, true
		}
		return nil, false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch value.Interface().(type) {
		case time.Duration:
			return config.setDuration(value, v)
		}
		return config.setInt(value, v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return config.setUint(value, v)
	case reflect.Float32, reflect.Float64:
		return config.setFloat(value, v)
	case reflect.Bool:
		return config.setBool(value, v)
	case reflect.String:
		if v != "" {
			value.SetString(v)
		}
		return nil, v != ""
	case reflect.Struct:
		switch value.Interface().(type) {
		case time.Time:
			return config.setTime(value, v)
		}
		return json.Unmarshal([]byte(v), value.Addr().Interface()), true
	case reflect.Slice:
		return config.setSlice(value, v)
	case reflect.Array:
		return config.setArray(value, v), true
	}

	return fmt.Errorf("unsupported type: %v", value.Kind()), false
}

func (config *Config) setArray(value reflect.Value, v string) error {
	parts := strings.Split(v, config.SliceSeparator)

	if len(parts) != value.Len() {
		return fmt.Errorf("cannot set array: expected %d parts, got %d", value.Len(), len(parts))
	}

	for i, part := range parts {
		if err, _ := config.setValue(value.Index(i), part); err != nil {
			return err
		}
	}
	return nil
}

func (config *Config) setSlice(value reflect.Value, v string) (error, bool) {
	if v == "" {
		return nil, false
	}

	parts := strings.Split(v, config.SliceSeparator)
	slice := reflect.MakeSlice(value.Type(), len(parts), len(parts))
	anySet := false
	for i, part := range parts {
		err, ok := config.setValue(slice.Index(i), part)
		if err != nil {
			return err, false
		}
		if ok {
			anySet = true
		}
	}
	value.Set(slice)
	return nil, anySet
}

func (config *Config) setInt(value reflect.Value, v string) (error, bool) {
	if v == "" {
		return nil, false
	}

	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return err, false
	}
	value.SetInt(i)
	return nil, true
}

func (config *Config) setUint(value reflect.Value, v string) (error, bool) {
	if v == "" {
		return nil, false
	}

	i, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return err, false
	}
	value.SetUint(i)
	return nil, true
}

func (config *Config) setFloat(value reflect.Value, v string) (error, bool) {
	if v == "" {
		return nil, false
	}

	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return err, false
	}
	value.SetFloat(f)
	return nil, true
}

func (config *Config) setBool(value reflect.Value, v string) (error, bool) {
	if v == "" {
		return nil, false
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		return err, false
	}
	value.SetBool(b)
	return nil, true
}

func (config *Config) setDuration(value reflect.Value, v string) (error, bool) {
	if v == "" {
		return nil, false
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		return err, false
	}

	value.Set(reflect.ValueOf(d))
	return nil, true
}

func (config *Config) setTime(value reflect.Value, v string) (error, bool) {
	if v == "" {
		return nil, false
	}

	t, err := time.Parse(config.TimeFormat, v)
	if err != nil {
		return err, false
	}
	value.Set(reflect.ValueOf(t))
	return nil, true
}
