/*
Copyright 2015 The Kubernetes Authors.

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

package benchmark

import (
	"testing"
	"time"
)

// BenchmarkScheduling100Nodes0Pods benchmarks the scheduling rate
// when the cluster has 100 nodes and 0 scheduled pods
func BenchmarkScheduling100Nodes0Pods(b *testing.B) {
	benchmarkScheduling(100, 0, b)
}

// BenchmarkScheduling100Nodes1000Pods benchmarks the scheduling rate
// when the cluster has 100 nodes and 1000 scheduled pods
func BenchmarkScheduling100Nodes1000Pods(b *testing.B) {
	benchmarkScheduling(100, 1000, b)
}

// BenchmarkScheduling1000Nodes0Pods benchmarks the scheduling rate
// when the cluster has 1000 nodes and 0 scheduled pods
func BenchmarkScheduling1000Nodes0Pods(b *testing.B) {
	benchmarkScheduling(1000, 0, b)
}

// BenchmarkScheduling1000Nodes1000Pods benchmarks the scheduling rate
// when the cluster has 1000 nodes and 1000 scheduled pods
func BenchmarkScheduling1000Nodes1000Pods(b *testing.B) {
	benchmarkScheduling(1000, 1000, b)
}

// benchmarkScheduling benchmarks scheduling rate with specific number of nodes
// and specific number of pods already scheduled. Since an operation takes relatively
// long time, b.N should be small: 10 - 100.
func benchmarkScheduling(numNodes, numScheduledPods int, b *testing.B) {
	schedulerConfigFactory, finalFunc := mustSetupScheduler()
	defer finalFunc()
	c := schedulerConfigFactory.Client

	makeNodes(c, numNodes)
	makePodsFromRC(c, "rc1", numScheduledPods)
	for {
		scheduled := schedulerConfigFactory.ScheduledPodLister.Indexer.List()
		if len(scheduled) >= numScheduledPods {
			break
		}
		time.Sleep(1 * time.Second)
	}
	// start benchmark
	b.ResetTimer()
	makePodsFromRC(c, "rc2", b.N)
	for {
		// This can potentially affect performance of scheduler, since List() is done under mutex.
		// TODO: Setup watch on apiserver and wait until all pods scheduled.
		scheduled := schedulerConfigFactory.ScheduledPodLister.Indexer.List()
		if len(scheduled) >= numScheduledPods+b.N {
			break
		}
		// Note: This might introduce slight deviation in accuracy of benchmark results.
		// Since the total amount of time is relatively large, it might not be a concern.
		time.Sleep(100 * time.Millisecond)
	}
}
