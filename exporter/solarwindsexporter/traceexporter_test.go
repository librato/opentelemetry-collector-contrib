// Copyright OpenTelemetry Authors
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

package solarwindsexporter

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.uber.org/zap"
)

func TestTraceExporter(t *testing.T) {

	time.Sleep(3 * time.Second)
	factory := NewFactory()
	exporterCfg := factory.CreateDefaultConfig().(*Config)
	params := component.ExporterCreateParams{Logger: zap.NewNop()}
	te, err := factory.CreateTraceExporter(context.Background(), params, exporterCfg)

	assert.NoError(t, err)
	assert.NotNil(t, te, "failed to create trace exporter")
	te.Start(context.Background(), nil)

	traces := pdata.NewTraces()
	resourceSpans := traces.ResourceSpans()
	resourceSpans.Resize(1)
	resourceSpans.At(0).InitEmpty()
	resourceSpans.At(0).InstrumentationLibrarySpans().Resize(1)
	resourceSpans.At(0).InstrumentationLibrarySpans().At(0).Spans().Resize(1)
	span := resourceSpans.At(0).InstrumentationLibrarySpans().At(0).Spans().At(0)
	span.SetName("foobar2")
	err = te.ConsumeTraces(context.Background(), traces)
	assert.NoError(t, err)
	assert.NoError(t, te.Shutdown(context.Background()))
}
