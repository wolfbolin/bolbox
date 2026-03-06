package configs

import (
	"fmt"
	"os"
	"reflect"

	flag "github.com/spf13/pflag"
	"github.com/wolfbolin/bolbox/pkg/errors"
)

// parseFlags 解析命令行参数并更新配置
func (m *Manager[T]) parseFlags() error {
	t := m.confElem.Type()
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)
	flagSet.ParseErrorsAllowlist.UnknownFlags = true

	changes := make([]func(), 0)
	for i := 0; i < t.NumField(); i++ {
		fieldKey := t.Field(i)
		conf, _ := m.Conf(fieldKey.Name)
		err := parseFlag(flagSet, fieldKey, conf, &changes)
		if err != nil {
			return err
		}
	}

	//flagSet.SetOutput(io.Discard)
	err := flagSet.Parse(os.Args[1:])
	if errors.Is(err, flag.ErrHelp) {
		if m.Options.ExitOnHelp {
			os.Exit(0)
		}
		return errors.Join(ErrPrintUsage, err)
	}
	if err != nil {
		return errors.Wrapf(ErrParseFlags, "Parse with flag set failed. %s", err.Error())
	}
	for _, change := range changes {
		change()
	}
	return nil
}

// parseFlag 为单个配置字段注册命令行参数
func parseFlag(flagSet *flag.FlagSet, fieldKey reflect.StructField, conf *Config, changes *[]func()) error {
	flagName := fieldKey.Tag.Get("flag")
	flagDesc := fieldKey.Tag.Get("desc")

	if flagName == "" {
		return nil
	}
	if flagDesc == "" {
		flagDesc = fmt.Sprintf("Flag for %s", flagName)
	}

	switch conf.val.Kind() {
	case reflect.String:
		valAddr := flagSet.String(flagName, conf.val.String(), flagDesc)
		*changes = append(*changes, func() {
			_ = conf.SetByValue(*valAddr)
		})
	case reflect.Bool:
		valAddr := flagSet.Bool(flagName, conf.val.Bool(), flagDesc)
		*changes = append(*changes, func() {
			_ = conf.SetByValue(*valAddr)
		})
	case reflect.Int, reflect.Int32, reflect.Int64:
		valAddr := flagSet.Int64(flagName, conf.val.Int(), flagDesc)
		*changes = append(*changes, func() {
			_ = conf.SetByValue(*valAddr)
		})
	case reflect.Float32, reflect.Float64:
		valAddr := flagSet.Float64(flagName, conf.val.Float(), flagDesc)
		*changes = append(*changes, func() {
			_ = conf.SetByValue(*valAddr)
		})
	case reflect.Map:
		valAddr := flagSet.String(flagName, conf.val.String(), flagDesc)
		*changes = append(*changes, func() {
			_ = conf.SetByValue(*valAddr)
		})
	default:
		return errors.Wrapf(ErrParseFlags, "Not suppose parse process flag[%s] for var[%s]", flagName, fieldKey.Name)
	}
	return nil
}
