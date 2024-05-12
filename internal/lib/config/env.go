package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type envProvider struct {
	prefix         string
	valueProviders map[string]converterSet
}

type converterSet struct {
	converters map[string]converterFunc
}

type converterFunc func(string) error

type SupportsEnvironmentVariables interface {
	SetupEnvironmentVariables(source EnvironmentVariableSource) error
}

type EnvironmentVariableSource interface {
	AddBool(target *bool, name string)
	AddInt8(target *int8, name string)
	AddInt16(target *int16, name string)
	AddInt32(target *int32, name string)
	AddInt64(target *int64, name string)
	AddUInt8(target *uint8, name string)
	AddUInt16(target *uint16, name string)
	AddUInt32(target *uint32, name string)
	AddUInt64(target *uint64, name string)
	AddFloat32(target *float32, name string)
	AddFloat64(target *float64, name string)
	AddString(target *string, name string)
	AddRaw(target *interface{}, name string, converter func(value string) (interface{}, error))
}

func NewEnvProvider(prefix string) Provider {
	return &envProvider{prefix: prefix}
}

func (e *envProvider) Name() string {
	return "env"
}

func (e *envProvider) ReadConfigSectionDefinition(name string, v interface{}) error {
	if _, ok := e.valueProviders[name]; ok {
		return fmt.Errorf("cannot override existing config section %s", name)
	}

	if supports, ok := v.(SupportsEnvironmentVariables); ok {
		source := &converterSetSource{}
		err := supports.SetupEnvironmentVariables(source)
		if err != nil {
			return err
		}
		e.valueProviders[name] = converterSet{source.converters}
	}

	return nil
}

func (e *envProvider) Parse() error {
	for section, set := range e.valueProviders {
		for value, converter := range set.converters {
			source := e.getEnvironmentVariableName(section, value)
			value := os.Getenv(source)
			if value == "" {
				continue
			}

			err := converter(value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *envProvider) getEnvironmentVariableName(section string, value string) string {
	data := fmt.Sprintf("%s_%s_%s", e.prefix, section, value)
	return strings.ToUpper(data)
}

type converterSetSource struct {
	converters map[string]converterFunc
}

func (c *converterSetSource) AddBool(target *bool, name string) {
	c.converters[name] = func(s string) error {
		t, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		*target = t
		return nil
	}
}
func (c *converterSetSource) AddInt8(target *int8, name string) {
	c.converters[name] = func(s string) error {
		tmp, err := strconv.ParseInt(s, 0, 8)
		if err != nil {
			return err
		}
		t := int8(tmp)
		*target = t
		return nil
	}
}
func (c *converterSetSource) AddInt16(target *int16, name string) {
	c.converters[name] = func(s string) error {
		tmp, err := strconv.ParseInt(s, 0, 16)
		if err != nil {
			return err
		}
		t := int16(tmp)
		*target = t
		return nil
	}
}
func (c *converterSetSource) AddInt32(target *int32, name string) {
	c.converters[name] = func(s string) error {
		tmp, err := strconv.ParseInt(s, 0, 32)
		if err != nil {
			return err
		}
		t := int32(tmp)
		*target = t
		return nil
	}
}
func (c *converterSetSource) AddInt64(target *int64, name string) {
	c.converters[name] = func(s string) error {
		tmp, err := strconv.ParseInt(s, 0, 64)
		if err != nil {
			return err
		}
		t := int64(tmp)
		*target = t
		return nil
	}
}
func (c *converterSetSource) AddUInt8(target *uint8, name string) {
	c.converters[name] = func(s string) error {
		tmp, err := strconv.ParseUint(s, 0, 8)
		if err != nil {
			return err
		}
		t := uint8(tmp)
		*target = t
		return nil
	}
}
func (c *converterSetSource) AddUInt16(target *uint16, name string) {
	c.converters[name] = func(s string) error {
		tmp, err := strconv.ParseUint(s, 0, 16)
		if err != nil {
			return err
		}
		t := uint16(tmp)
		*target = t
		return nil
	}
}
func (c *converterSetSource) AddUInt32(target *uint32, name string) {
	c.converters[name] = func(s string) error {
		tmp, err := strconv.ParseUint(s, 0, 32)
		if err != nil {
			return err
		}
		t := uint32(tmp)
		*target = t
		return nil
	}
}
func (c *converterSetSource) AddUInt64(target *uint64, name string) {
	c.converters[name] = func(s string) error {
		tmp, err := strconv.ParseUint(s, 0, 64)
		if err != nil {
			return err
		}
		t := uint64(tmp)
		*target = t
		return nil
	}
}
func (c *converterSetSource) AddFloat32(target *float32, name string) {
	c.converters[name] = func(s string) error {
		tmp, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return err
		}
		t := float32(tmp)
		*target = t
		return nil
	}
}
func (c *converterSetSource) AddFloat64(target *float64, name string) {
	c.converters[name] = func(s string) error {
		tmp, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		t := float64(tmp)
		*target = t
		return nil
	}
}
func (c *converterSetSource) AddString(target *string, name string) {
	c.converters[name] = func(s string) error {
		*target = s
		return nil
	}
}
func (c *converterSetSource) AddRaw(target *interface{}, name string, converter func(value string) (interface{}, error)) {
	c.converters[name] = func(s string) error {
		t, err := converter(s)
		if err != nil {
			return err
		}
		*target = t
		return nil
	}
}
