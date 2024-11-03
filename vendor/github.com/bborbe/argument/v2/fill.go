// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package argument

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/bborbe/errors"
)

// Fill the given map into the struct.
func Fill(ctx context.Context, data interface{}, values map[string]interface{}) error {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(values); err != nil {
		return errors.Wrap(ctx, err, "encode json failed")
	}
	if err := json.NewDecoder(buf).Decode(data); err != nil {
		return errors.Wrap(ctx, err, "decode json failed")
	}
	return nil
}
