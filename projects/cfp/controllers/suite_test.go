/*
Copyright 2022.

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

package controllers

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/fluxcd/pkg/runtime/testenv"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	"k8s.io/client-go/kubernetes/scheme"

	//+kubebuilder:scaffold:imports

	talksv1 "github.com/scottrigby/how-to-write-a-reconciler-using-k8s-controller-runtime/projects/cfp/api/v1"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

const timeout = 10 * time.Second

var testEnv *testenv.Environment

func TestMain(m *testing.M) {
	utilruntime.Must(talksv1.AddToScheme(scheme.Scheme))

	testEnv = testenv.New(testenv.WithCRDPath(filepath.Join("..", "config", "crd", "bases")))

	if err := (&SpeakerReconciler{
		Client: testEnv.Client,
	}).SetupWithManager(testEnv.Manager); err != nil {
		panic(err)
	}
}
