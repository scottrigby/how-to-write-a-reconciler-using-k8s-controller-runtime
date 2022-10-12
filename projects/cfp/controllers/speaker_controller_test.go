package controllers

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/onsi/gomega"
	talksv1 "github.com/scottrigby/how-to-write-a-reconciler-using-k8s-controller-runtime/projects/cfp/api/v1"
)

func Test_Speaker_Reconcile(t *testing.T) {
	g := NewGomegaWithT(t)

	testCases := []struct {
		name    string
		Speaker string
		Bio     string
	}{
		{
			name:    "test that is set",
			Speaker: "Speaker one",
			Bio:     "First speaker bio",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			obj := &talksv1.Speaker{
				ObjectMeta: metav1.ObjectMeta{
					Name:      tc.name,
					Namespace: "default",
				},
				Spec: talksv1.SpeakerSpec{
					Name: tc.Speaker,
					Bio:  tc.Bio,
				},
			}

			ctx := context.Background()

			err := testEnv.CreateAndWait(ctx, obj)
			g.Expect(err).NotTo(HaveOccurred())

			key := client.ObjectKey{Name: obj.Name, Namespace: obj.Namespace}

			// Wait for finalizer to be set
			g.Eventually(func() bool {
				if err := testEnv.Get(ctx, key, obj); err != nil {
					return false
				}
				return len(obj.Finalizers) > 0
			}, timeout).Should(BeTrue())

		})
	}
}
