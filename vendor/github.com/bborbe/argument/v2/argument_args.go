// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package argument

import (
	"context"
	"flag"
	"reflect"
	"strconv"
	"time"

	"github.com/bborbe/errors"
	libtime "github.com/bborbe/time"
)

// ParseArgs into the given struct.
func ParseArgs(ctx context.Context, data interface{}, args []string) error {
	values, err := argsToValues(ctx, data, args)
	if err != nil {
		return errors.Wrapf(ctx, err, "args to values failed")
	}
	if err := Fill(ctx, data, values); err != nil {
		return errors.Wrapf(ctx, err, "fill failed")
	}
	return nil
}

func argsToValues(ctx context.Context, data interface{}, args []string) (map[string]interface{}, error) {
	e := reflect.ValueOf(data).Elem()
	t := e.Type()
	values := make(map[string]interface{})
	for i := 0; i < e.NumField(); i++ {
		tf := t.Field(i)
		ef := e.Field(i)
		argName, ok := tf.Tag.Lookup("arg")
		if !ok {
			continue
		}
		defaultString, found := tf.Tag.Lookup("default")
		usage := tf.Tag.Get("usage")
		switch ef.Interface().(type) {
		case string:
			values[tf.Name] = flag.CommandLine.String(argName, defaultString, usage)
		case bool:
			defaultValue, _ := strconv.ParseBool(defaultString)
			values[tf.Name] = flag.CommandLine.Bool(argName, defaultValue, usage)
		case int:
			defaultValue, _ := strconv.Atoi(defaultString)
			values[tf.Name] = flag.CommandLine.Int(argName, defaultValue, usage)
		case int64:
			defaultValue, _ := strconv.ParseInt(defaultString, 10, 0)
			values[tf.Name] = flag.CommandLine.Int64(argName, defaultValue, usage)
		case uint:
			defaultValue, _ := strconv.ParseUint(defaultString, 10, 0)
			values[tf.Name] = flag.CommandLine.Uint(argName, uint(defaultValue), usage)
		case uint64:
			defaultValue, _ := strconv.ParseUint(defaultString, 10, 0)
			values[tf.Name] = flag.CommandLine.Uint64(argName, defaultValue, usage)
		case int32:
			defaultValue, _ := strconv.ParseInt(defaultString, 10, 0)
			values[tf.Name] = flag.CommandLine.Int(argName, int(defaultValue), usage)
		case float64:
			defaultValue, _ := strconv.ParseFloat(defaultString, 64)
			values[tf.Name] = flag.CommandLine.Float64(argName, defaultValue, usage)
		case *float64:
			if found {
				defaultValue, _ := strconv.ParseFloat(defaultString, 64)
				values[tf.Name] = defaultValue
			}
			flag.CommandLine.Func(argName, usage, func(s string) error {
				if s == "" {
					return nil
				}
				v, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return errors.Wrapf(ctx, err, "parse float failed")
				}
				values[tf.Name] = v
				return nil
			})
		case time.Duration:
			if found {
				defaultValue, _ := libtime.ParseDuration(ctx, defaultString)
				if defaultValue != nil {
					values[tf.Name] = defaultValue.Duration()
				}
			}
			flag.CommandLine.Func(argName, usage, func(value string) error {
				if value == "" {
					return nil
				}
				duration, err := libtime.ParseDuration(ctx, value)
				if err != nil {
					return errors.Wrapf(ctx, err, "parse duration failed")
				}
				values[tf.Name] = duration.Duration()
				return nil
			})
		default:
			return nil, errors.Errorf(ctx, "field %s with type %T is unsupported", tf.Name, ef.Interface())
		}
	}
	if err := flag.CommandLine.Parse(args); err != nil {
		return nil, errors.Wrapf(ctx, err, "parse commandline failed")
	}
	return values, nil
}
