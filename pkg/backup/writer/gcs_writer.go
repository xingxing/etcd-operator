// Copyright 2017 The etcd-operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package writer

import (
	"context"
	"fmt"
	"io"

	"github.com/coreos/etcd-operator/pkg/backup/util"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
)

var _ Writer = &gcsWriter{}

type gcsWriter struct {
	gcs *storage.Client
}

// NewGCSWriter creates a gcs writer.
func NewGCSWriter(gcs *storage.Client) Writer {
	return &gcsWriter{gcs}
}

// Write writes the backup file to the given gcs path, "<gcs-bucket-name>/<key>".
func (gcsw *gcsWriter) Write(ctx context.Context, path string, r io.Reader) (int64, error) {
	// TODO: support context.
	bucket, key, err := util.ParseBucketAndKey(path)
	if err != nil {
		return 0, err
	}

	w := gcsw.gcs.Bucket(bucket).Object(key).NewWriter(ctx)
	defer func() {
		err := w.Close()
		if err != nil {
			logrus.Errorf("failed to close GCS object writer: %v", err)
		}
	}()

	n, err := io.Copy(w, r)
	if err != nil {
		err = fmt.Errorf("failed to write GCS object: %v", err)
	}
	return n, err
}
