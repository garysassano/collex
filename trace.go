// Copyright 2022 Tyler Yahn (MrAlias)
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

// Package collex provides OpenTelemetry Go Exporters that wrap OpenTelemetry
// Collector Exporters. This allows any collector exporter to be used with
// opentelemetry-go.
package collex

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/otel/sdk/trace"
)

// TracesExporter returns an OpenTelemetry-Go SpanExporter that wraps e.
func TracesExporter(e component.TracesExporter) trace.SpanExporter {
	// TODO
	return nil
}
