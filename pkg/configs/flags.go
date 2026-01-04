package configs

import (
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/wolfbolin/bolbox/pkg/log"
)

// parseFlags 解析命令行参数并更新配置
// 该函数会遍历配置结构体的所有字段，查找带有 "flag" 标签的字段，
// 并根据字段类型自动注册对应的命令行参数。
// 支持的字段类型包括：string、bool、int、int32、int64、float32、float64
// 解析失败时会记录警告日志，但不会中断程序执行
// 如果用户显式调用 --help 或 -h，会打印帮助信息并退出程序
func (m *Manager[T]) parseFlags() {
	t := m.confElem.Type()
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)
	hookList := make([]func(), 0)
	for i := 0; i < t.NumField(); i++ {
		fieldKey := t.Field(i)
		fieldVal := m.confElem.Field(i)
		if !fieldVal.CanSet() {
			continue
		}
		parseFlag(flagSet, fieldKey, fieldVal, &hookList)
	}

	// 在函数开始处检查并显示帮助信息
	checkAndShowHelp(flagSet)

	// 非帮助模式，禁止输出到标准输出
	//flagSet.Usage = func() {}
	flagSet.SetOutput(io.Discard)

	// 忽略解析错误，不输出任何内容
	_ = flagSet.Parse(os.Args[1:])
	for _, hook := range hookList {
		hook()
	}
}

// parseFlag 为单个配置字段注册命令行参数
// 参数说明：
//   - flagSet: 命令行参数集合，用于注册新的参数
//   - fieldKey: 结构体字段的元信息，包含标签信息
//   - fieldVal: 字段的值，用于设置默认值和接收解析结果
//   - hookList: 钩子函数列表，用于处理需要类型转换的字段（如 int32、float32）
//
// 该函数会从字段的 "flag" 标签中获取参数名，从 "desc" 标签中获取参数描述。
// 如果未提供描述，将使用默认描述 "Flag for {参数名}"。
// 对于 int32 和 float32 类型，由于 Go 标准库 flag 包的限制，需要使用临时变量和钩子函数进行类型转换。
func parseFlag(flagSet *flag.FlagSet, fieldKey reflect.StructField, fieldVal reflect.Value, hookList *[]func()) {
	flagName := fieldKey.Tag.Get("flag")
	flagDesc := fieldKey.Tag.Get("desc")

	if flagName == "" {
		return
	}
	if flagDesc == "" {
		flagDesc = fmt.Sprintf("Flag for %s", flagName)
	}

	switch fieldVal.Kind() {
	case reflect.String:
		flagSet.StringVar(fieldVal.Addr().Interface().(*string), flagName, fieldVal.String(), flagDesc)
	case reflect.Bool:
		flagSet.BoolVar(fieldVal.Addr().Interface().(*bool), flagName, fieldVal.Bool(), flagDesc)
	case reflect.Int:
		flagSet.IntVar(fieldVal.Addr().Interface().(*int), flagName, int(fieldVal.Int()), flagDesc)
	case reflect.Int32:
		var temp int64
		flagSet.Int64Var(&temp, flagName, fieldVal.Int(), flagDesc)
		*hookList = append(*hookList, func() {
			fieldVal.SetInt(temp)
		})
	case reflect.Float32:
		var temp float64
		flagSet.Float64Var(&temp, flagName, fieldVal.Float(), flagDesc)
		*hookList = append(*hookList, func() {
			fieldVal.SetFloat(temp)
		})
	case reflect.Int64:
		flagSet.Int64Var(fieldVal.Addr().Interface().(*int64), flagName, fieldVal.Int(), flagDesc)
	case reflect.Float64:
		flagSet.Float64Var(fieldVal.Addr().Interface().(*float64), flagName, fieldVal.Float(), flagDesc)
	default:
		log.Warnf("Not suppose parse process flag[%s] for var[%s]", flagName, fieldKey.Name)
	}
}

// checkAndShowHelp 检查是否显式调用了 --help 或 -h，如果是则显示帮助信息并退出程序
func checkAndShowHelp(flagSet *flag.FlagSet) {
	for _, arg := range os.Args[1:] {
		if arg == "--help" || arg == "-h" {
			flagSet.Usage()
			os.Exit(0)
		}
	}
}
