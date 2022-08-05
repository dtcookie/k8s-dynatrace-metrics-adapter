/*
Copyright 2021 Reinhard Pilz.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/component-base/logs"
	"k8s.io/klog/v2"

	basecmd "github.com/kubernetes-sigs/custom-metrics-apiserver/pkg/cmd"
	"github.com/kubernetes-sigs/custom-metrics-apiserver/pkg/provider"
)

type MetricsAdapter struct {
	basecmd.AdapterBase

	// Message is printed on successful startup
	Message string
}

func main() {
	if os.Getenv("KUBERNETES_SERVICE_HOST") == "" {
		localMain()
	} else {
		k8sMain()
	}
}

func k8sMain() {
	logs.InitLogs()
	defer logs.FlushLogs()

	cmd := &MetricsAdapter{}

	cmd.Flags().StringVar(&cmd.Message, "msg", "starting adapter...", "startup message")
	cmd.Flags().AddGoFlagSet(flag.CommandLine) // make sure we get the klog flags
	cmd.Flags().Parse(os.Args)

	cmd.WithExternalMetrics(NewProvider())

	klog.Infof(cmd.Message)
	if err := cmd.Run(wait.NeverStop); err != nil {
		klog.Fatalf("unable to run custom metrics  adapter: %v", err)
	}
}

func ListAllExternalMetrics(prov provider.ExternalMetricsProvider) {
	for {
		metricInfos := prov.ListAllExternalMetrics()
		klog.Infof("%d metrics available", len(metricInfos))
		time.Sleep(30 * time.Second)
	}
}

func localMain() {
	prov := NewProvider()

	go ListAllExternalMetrics(prov)
	for {
		selector := labels.NewSelector()

		valueList, err := prov.GetExternalMetric("default", selector,
			provider.ExternalMetricInfo{Metric: "dsfm:synthetic.browser.engine_utilization:max:filter(eq(\"dt.entity.synthetic_location\", \"SYNTHETIC_LOCATION-7BA305221EFA8DBF\")):merge(\"location.name\", \"host.name\", \"dt.active_gate.working_mode\", \"dt.active_gate.id\"):last"})
		if err != nil {
			klog.Error("failed to query external metric", err.Error())
		}
		if valueList != nil {
			klog.Infof("%d metric values available", len(valueList.Items))
		} else {
			klog.Infof("%d metric values available", 0)
		}
		time.Sleep(15 * time.Second)
	}

	// selector := labels.NewSelector()
	// valueList, err := prov.GetExternalMetric("default", selector, provider.ExternalMetricInfo{Metric: "com.dynatrace.builtin:pgi.jvm.garbagecollectioncount_"})
	// if err != nil {
	// 	panic(err)
	// }
	// if valueList != nil {
	// 	for idx, item := range valueList.Items {
	// 		data, _ := json.Marshal(item)
	// 		fmt.Println(idx, string(data))
	// 	}
	// }

}
