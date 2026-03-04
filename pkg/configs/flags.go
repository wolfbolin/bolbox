package configs

import (
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/wolfbolin/bolbox/pkg/errors"
)

// parseFlags 解析命令行参数并更新配置
func (m *Manager[T]) parseFlags() error {
	t := m.confElem.Type()
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)
	changes := make([]func(), 0)
	for i := 0; i < t.NumField(); i++ {
		fieldKey := t.Field(i)
		conf, _ := m.Conf(fieldKey.Name)
		err := parseFlag(flagSet, fieldKey, conf, &changes)
		if err != nil {
			return err
		}
	}

	checkAndShowHelp(flagSet)
	flagSet.Usage = func() {}
	flagSet.SetOutput(io.Discard)

	args := os.Args[1:]
	filteredArgs := filterKnownArgs(flagSet, args)

	err := flagSet.Parse(filteredArgs)
	if err != nil {
		return errors.Wrapf(ParseFlagsError, "Parse with flag set failed. %s", err.Error())
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
		return errors.Wrapf(ParseFlagsError, "Not suppose parse process flag[%s] for var[%s]", flagName, fieldKey.Name)
	}
	return nil
}

// checkAndShowHelp 检查并显示帮助信息
func checkAndShowHelp(flagSet *flag.FlagSet) {
	for _, arg := range os.Args[1:] {
		if arg == "--help" || arg == "-h" {
			flagSet.Usage()
			os.Exit(0)
		}
	}
}

// filterKnownArgs 过滤掉 flagSet 中未注册的命令行参数，避免出现 "flag provided but not defined" 错误
func filterKnownArgs(flagSet *flag.FlagSet, args []string) []string {
	known := make(map[string]bool)
	flagSet.VisitAll(func(f *flag.Flag) {
		known[f.Name] = true
	})

	result := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		name, hasValue := parseFlagName(args[i])
		if name == "" {
			// 非 flag 形式（位置参数），直接保留
			result = append(result, args[i])
			continue
		}
		if !known[name] {
			// 未知 flag：若是 --flag value 形式，额外跳过紧随的 value
			if !hasValue && i+1 < len(args) && !looksLikeFlag(args[i+1]) {
				i++
			}
			continue
		}
		result = append(result, args[i])
		// 已知 flag 的 --flag value 形式：把紧随的 value 一并保留
		if !hasValue && i+1 < len(args) && !looksLikeFlag(args[i+1]) {
			i++
			result = append(result, args[i])
		}
	}
	return result
}

// parseFlagName 从单个参数中解析出 flag 名。
// 返回 (flagName, hasValue)：hasValue 表示值已内联（--flag=value 形式）。
// 若参数不是 flag（不以 '-' 开头），返回 ("", false)。
func parseFlagName(arg string) (name string, hasValue bool) {
	if len(arg) == 0 || arg[0] != '-' {
		return "", false
	}
	// 去掉 - 或 --
	name = strings.TrimPrefix(strings.TrimPrefix(arg, "--"), "-")
	// 处理 --flag=value 形式
	if before, _, found := strings.Cut(name, "="); found {
		return before, true
	}
	return name, false
}

// looksLikeFlag 判断一个参数是否是 flag（以 '-' 后跟字母开头，排除纯负数如 -123）
func looksLikeFlag(arg string) bool {
	if len(arg) < 2 || arg[0] != '-' {
		return false
	}
	return arg[1] < '0' || arg[1] > '9'
}
