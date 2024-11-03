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

// DefaultValues returns all default values of the given struct.
func DefaultValues(ctx context.Context, data interface{}) (map[string]interface{}, error) {
	var err error
	e := reflect.ValueOf(data).Elem()
	t := e.Type()
	values := make(map[string]interface{})
	for i := 0; i < e.NumField(); i++ {
		tf := t.Field(i)
		ef := e.Field(i)
		value, ok := tf.Tag.Lookup("default")
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
		case *float64:
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
