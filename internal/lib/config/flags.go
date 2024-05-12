package config

import (
	"encoding"
	"flag"
	"fmt"
	"os"
	"time"
)

type flagConfigProvider struct {
	flagSets map[string]*flag.FlagSet
}

type SupportsFlagConfigSectionDefinition interface {
	Configure(f FlagProvider) error
}

type FlagProvider interface {
	BoolVar(p *bool, name string, value bool, usage string)
	IntVar(p *int, name string, value int, usage string)
	Int64Var(p *int64, name string, value int64, usage string)
	UintVar(p *uint, name string, value uint, usage string)
	Uint64Var(p *uint64, name string, value uint64, usage string)
	StringVar(p *string, name string, value string, usage string)
	Float64Var(p *float64, name string, value float64, usage string)
	DurationVar(p *time.Duration, name string, value time.Duration, usage string)
	TextVar(p encoding.TextUnmarshaler, name string, value encoding.TextMarshaler, usage string)
	Var(value flag.Value, name string, usage string)
}

func NewFlagConfigProvider() Provider {
	return &flagConfigProvider{
		flagSets: make(map[string]*flag.FlagSet),
	}
}

func (f *flagConfigProvider) Name() string {
	return "flags"
}

func (f *flagConfigProvider) Parse() error {
	for key, value := range f.flagSets {
		err := value.Parse(os.Args[1:])
		if err != nil {
			return fmt.Errorf("invalid flag configuration '%s': %v", key, err)
		}
	}

	return nil
}

func (f *flagConfigProvider) ReadConfigSectionDefinition(name string, v interface{}) error {
	if _, ok := f.flagSets[name]; ok {
		return fmt.Errorf("duplicate configuration section %v", name)
	}

	supportsFlagging, ok := v.(SupportsFlagConfigSectionDefinition)
	if !ok {
		return nil
	}

	flagSet := flag.NewFlagSet(name, flag.ContinueOnError)
	builder := &flagBuilder{flagSet: flagSet, name: name}
	err := supportsFlagging.Configure(builder)
	if err != nil {
		return err
	}

	f.flagSets[name] = flagSet
	return nil
}

type flagBuilder struct {
	flagSet *flag.FlagSet
	name    string
}

func (f *flagBuilder) BoolVar(p *bool, name string, value bool, usage string) {
	f.flagSet.BoolVar(p, f.name+"-"+name, value, usage)
	flag.BoolVar(p, f.name+"-"+name, value, usage)
}

func (f *flagBuilder) IntVar(p *int, name string, value int, usage string) {
	f.flagSet.IntVar(p, f.name+"-"+name, value, usage)
	flag.IntVar(p, f.name+"-"+name, value, usage)
}

func (f *flagBuilder) Int64Var(p *int64, name string, value int64, usage string) {
	f.flagSet.Int64Var(p, f.name+"-"+name, value, usage)
	flag.Int64Var(p, f.name+"-"+name, value, usage)
}

func (f *flagBuilder) UintVar(p *uint, name string, value uint, usage string) {
	f.flagSet.UintVar(p, f.name+"-"+name, value, usage)
	flag.UintVar(p, f.name+"-"+name, value, usage)
}

func (f *flagBuilder) Uint64Var(p *uint64, name string, value uint64, usage string) {
	f.flagSet.Uint64Var(p, f.name+"-"+name, value, usage)
	flag.Uint64Var(p, f.name+"-"+name, value, usage)
}

func (f *flagBuilder) StringVar(p *string, name string, value string, usage string) {
	f.flagSet.StringVar(p, f.name+"-"+name, value, usage)
	flag.StringVar(p, f.name+"-"+name, value, usage)
}

func (f *flagBuilder) Float64Var(p *float64, name string, value float64, usage string) {
	f.flagSet.Float64Var(p, f.name+"-"+name, value, usage)
	flag.Float64Var(p, f.name+"-"+name, value, usage)
}

func (f *flagBuilder) DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	f.flagSet.DurationVar(p, f.name+"-"+name, value, usage)
	flag.DurationVar(p, f.name+"-"+name, value, usage)
}

func (f *flagBuilder) TextVar(p encoding.TextUnmarshaler, name string, value encoding.TextMarshaler, usage string) {
	f.flagSet.TextVar(p, f.name+"-"+name, value, usage)
	flag.TextVar(p, f.name+"-"+name, value, usage)
}

func (f *flagBuilder) Var(value flag.Value, name string, usage string) {
	f.flagSet.Var(value, f.name+"-"+name, usage)
	flag.Var(value, f.name+"-"+name, usage)
}
