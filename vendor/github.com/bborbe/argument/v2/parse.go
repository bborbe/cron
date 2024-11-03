// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package argument

import (
	"context"
	"os"

	"github.com/bborbe/errors"
)

// Parse combines all functionality. It parse env args, fills it the struct, print all arguments and validate required fields are set.
func Parse(ctx context.Context, data interface{}) error {
	argsValues, err := argsToValues(ctx, data, os.Args[1:])
	if err != nil {
		return errors.Wrapf(ctx, err, "arg to values failed")
	}
	envValues, err := envToValues(ctx, data, os.Environ())
	if err != nil {
		return errors.Wrapf(ctx, err, "env to values failed")
	}
	defaultValues, err := DefaultValues(ctx, data)
	if err != nil {
		return errors.Wrapf(ctx, err, "default values failed")
	}
	if err := Fill(ctx, data, mergeValues(defaultValues, argsValues, envValues)); err != nil {
		return errors.Wrapf(ctx, err, "fill failed")
	}
	if err := Print(ctx, data); err != nil {
		return errors.Wrapf(ctx, err, "print failed")
	}
	if err := ValidateRequired(ctx, data); err != nil {
		return errors.Wrapf(ctx, err, "validate required failed")
	}
	return nil
}

func mergeValues(list ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, values := range list {
		for k, v := range values {
			result[k] = v
		}
	}
	return result
}
