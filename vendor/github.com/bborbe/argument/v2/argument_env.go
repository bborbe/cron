// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package argument

import (
	"context"
	"reflect"
	"strconv"
	"time"

	"github.com/bborbe/errors"
	libtime "github.com/bborbe/time"
)

// ParseEnv into the given struct.
func ParseEnv(ctx context.Context, data interface{}, environ []string) error {
	values, err := envToValues(ctx, data, environ)
	if err != nil {
		return errors.Wrapf(ctx, err, "env to values failed")
	}
	if err := Fill(ctx, data, values); err != nil {
		return errors.Wrapf(ctx, err, "fill failed")
	}
	return nil
}

func envToValues(ctx context.Context, data interface{}, environ []string) (map[string]interface{}, error) {
	var err error
	envValues := make(map[string]string)
	for _, env := range environ {
		for i := 0; i < len(env); i++ {
			if env[i] == '=' {
				envValues[env[:i]] = env[i+1:]
			}
		}
	}
	values := make(map[string]interface{})
	e := reflect.ValueOf(data).Elem()
	t := e.Type()
	for i := 0; i < e.NumField(); i++ {
		tf := t.Field(i)
		ef := e.Field(i)
		argName, ok := tf.Tag.Lookup("env")
		if !ok {
			continue
		}
		value, ok := envValues[argName]
		if !ok {
			continue
		}
		switch ef.Interface().(type) {
		case string:
			values[tf.Name] = value
		case bool:
			values[tf.Name], err = strconv.ParseBool(value)
			if err != nil {
				return nil, errors.Errorf(ctx, "parse field %s as %T failed: %v", tf.Name, ef.Interface(), err)
			}
		case int:
			values[tf.Name], err = strconv.Atoi(value)
			if err != nil {
				return nil, errors.Errorf(ctx, "parse field %s as %T failed: %v", tf.Name, ef.Interface(), err)
			}
		case int64:
			values[tf.Name], err = strconv.ParseInt(value, 10, 0)
			if err != nil {
				return nil, errors.Errorf(ctx, "parse field %s as %T failed: %v", tf.Name, ef.Interface(), err)
			}
		case uint:
			values[tf.Name], err = strconv.ParseUint(value, 10, 0)
			if err != nil {
				return nil, errors.Errorf(ctx, "parse field %s as %T failed: %v", tf.Name, ef.Interface(), err)
			}
		case uint64:
			values[tf.Name], err = strconv.ParseUint(value, 10, 0)
			if err != nil {
				return nil, errors.Errorf(ctx, "parse field %s as %T failed: %v", tf.Name, ef.Interface(), err)
			}
		case int32:
			v, err := strconv.ParseInt(value, 10, 0)
			if err != nil {
				return nil, errors.Errorf(ctx, "parse field %s as %T failed: %v", tf.Name, ef.Interface(), err)
			}
			values[tf.Name] = int32(v)
		case float64:
			values[tf.Name], err = strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, errors.Errorf(ctx, "parse field %s as %T failed: %v", tf.Name, ef.Interface(), err)
			}
		case time.Duration:
			duration, err := libtime.ParseDuration(ctx, value)
			if err != nil {
				return nil, errors.Errorf(ctx, "parse field %s as %T failed: %v", tf.Name, ef.Interface(), err)
			}
			values[tf.Name] = duration.Duration()
		default:
			return nil, errors.Errorf(ctx, "field %s with type %T is unsupported", tf.Name, ef.Interface())
		}
	}
	return values, nil
}
