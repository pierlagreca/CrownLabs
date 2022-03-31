// Copyright 2020-2022 Politecnico di Torino
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

package tenant_controller

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	tnOpinternalErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "tenant_operator_internal_errors",
		Help: "The number of errors occurred internally during the reconcile of the tenant operator",
	},
		[]string{"controller", "reason"},
	)

	tnTokenConsumed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tenant_cpu_token_consumed",
		Help: "The number of cpu token consumed by the tenant through the creation of new instances",
	},
		[]string{"name"},
	)
)

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(tnOpinternalErrors, tnTokenConsumed)
}
