package controllers

import (
	"strings"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/fluxcd/pkg/apis/meta"
	"github.com/fluxcd/pkg/runtime/conditions"
	"github.com/fluxcd/pkg/runtime/patch"
	. "github.com/onsi/gomega"
	talksv1 "github.com/scottrigby/how-to-write-a-reconciler-using-k8s-controller-runtime/cfp/api/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

func Test_Speaker_Reconcile(t *testing.T) {
	g := NewGomegaWithT(t)

	testCases := []struct {
		name             string
		Speaker          string
		Bio              string
		Email            string
		assertConditions []metav1.Condition
		assertFunc       func(obj *talksv1.Speaker, assertConditions []metav1.Condition)
	}{
		{
			name:    "test create speaker reconciliation",
			Speaker: "Luke Skywalker",
			Bio:     "First speaker bio",
			Email:   "first@protonmail.com",
			assertConditions: []metav1.Condition{
				*conditions.TrueCondition(meta.ReadyCondition, meta.SucceededReason, "reconciled '<name>' successfully"),
			},
			assertFunc: func(obj *talksv1.Speaker, assertConditions []metav1.Condition) {
				key := client.ObjectKey{Name: obj.Name, Namespace: obj.Namespace}
				// Wait for finalizer to be set
				g.Eventually(func() bool {
					if err := testEnv.Get(ctx, key, obj); err != nil {
						return false
					}
					return len(obj.Finalizers) > 0
				}, timeout).Should(BeTrue())

				// Wait for Speaker to be Ready
				g.Eventually(func() bool {
					if err := testEnv.Get(ctx, key, obj); err != nil {
						return false
					}
					if !conditions.IsReady(obj) {
						return false
					}
					readyCondition := conditions.Get(obj, meta.ReadyCondition)
					return obj.Generation == readyCondition.ObservedGeneration &&
						obj.Generation == obj.Status.ObservedGeneration
				}, timeout).Should(BeTrue())

				for k := range assertConditions {
					assertConditions[k].Message = strings.ReplaceAll(assertConditions[k].Message, "<name>", obj.Name)
				}
				g.Expect(obj.Status.Conditions).To(conditions.MatchConditions(assertConditions))

			},
		},
		{
			name:    "test update speaker reconciliation",
			Speaker: "Jesse Pinkman",
			Bio:     "Dealer of meth",
			Email:   "blue.meth@protonmail.com",
			assertConditions: []metav1.Condition{
				*conditions.TrueCondition(meta.ReadyCondition, meta.SucceededReason, "reconciled '<name>' successfully"),
			},
			assertFunc: func(obj *talksv1.Speaker, assertConditions []metav1.Condition) {
				key := client.ObjectKey{Name: obj.Name, Namespace: obj.Namespace}
				// Wait for finalizer to be set
				g.Eventually(func() bool {
					if err := testEnv.Get(ctx, key, obj); err != nil {
						return false
					}
					return len(obj.Finalizers) > 0
				}, timeout).Should(BeTrue())

				// Wait for Speaker to be Ready
				g.Eventually(func() bool {
					if err := testEnv.Get(ctx, key, obj); err != nil {
						return false
					}
					if !conditions.IsReady(obj) {
						return false
					}
					readyCondition := conditions.Get(obj, meta.ReadyCondition)
					return obj.Generation == readyCondition.ObservedGeneration &&
						obj.Generation == obj.Status.ObservedGeneration
				}, timeout).Should(BeTrue())

				// Update Speaker
				previousGeneration := obj.Generation
				patch, err := patch.NewHelper(obj, testEnv)
				g.Expect(err).ToNot(HaveOccurred())
				obj.Spec.Bio = "Updated speaker bio"
				g.Expect(patch.Patch(ctx, obj)).To(Succeed())

				// Wait for Speaker to be Ready for this new generation
				g.Eventually(func() bool {
					if err := testEnv.Get(ctx, key, obj); err != nil {
						return false
					}
					if !conditions.IsReady(obj) {
						return false
					}
					readyCondition := conditions.Get(obj, meta.ReadyCondition)
					return obj.Generation == readyCondition.ObservedGeneration &&
						obj.Generation == obj.Status.ObservedGeneration &&
						obj.Generation == previousGeneration+1
				}, timeout).Should(BeTrue())

				for k := range assertConditions {
					assertConditions[k].Message = strings.ReplaceAll(assertConditions[k].Message, "<name>", obj.Name)
				}
				g.Expect(obj.Status.Conditions).To(conditions.MatchConditions(assertConditions))

			},
		},
		{
			name:    "test delete speaker reconciliation",
			Speaker: "Jesse Pinkman",
			Bio:     "Dealer of meth",
			Email:   "blue.meth@protonmail.com",
			assertFunc: func(obj *talksv1.Speaker, assertConditions []metav1.Condition) {
				key := client.ObjectKey{Name: obj.Name, Namespace: obj.Namespace}
				// Wait for finalizer to be set
				g.Eventually(func() bool {
					if err := testEnv.Get(ctx, key, obj); err != nil {
						return false
					}
					return len(obj.Finalizers) > 0
				}, timeout).Should(BeTrue())

				// Wait for Speaker to be Ready
				g.Eventually(func() bool {
					if err := testEnv.Get(ctx, key, obj); err != nil {
						return false
					}
					if !conditions.IsReady(obj) {
						return false
					}
					readyCondition := conditions.Get(obj, meta.ReadyCondition)
					return obj.Generation == readyCondition.ObservedGeneration &&
						obj.Generation == obj.Status.ObservedGeneration
				}, timeout).Should(BeTrue())

				// Delete Speaker
				g.Expect(testEnv.Delete(ctx, obj)).To(Succeed())

				// Wait for Speaker to be deleted
				g.Eventually(func() bool {
					if err := testEnv.Get(ctx, key, obj); err != nil {
						return apierrors.IsNotFound(err)
					}
					return false
				}, timeout).Should(BeTrue())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ns, err := testEnv.CreateNamespace(ctx, "speaker-ns")
			g.Expect(err).NotTo(HaveOccurred())
			defer func() {
				g.Expect(testEnv.Delete(ctx, ns)).To(Succeed())
			}()

			obj := &talksv1.Speaker{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "speaker",
					Namespace: ns.Name,
				},
				Spec: talksv1.SpeakerSpec{
					Name:  tc.Speaker,
					Bio:   tc.Bio,
					Email: tc.Email,
				},
			}

			err = testEnv.CreateAndWait(ctx, obj)
			g.Expect(err).NotTo(HaveOccurred())

			tc.assertFunc(obj, tc.assertConditions)
		})
	}
}
