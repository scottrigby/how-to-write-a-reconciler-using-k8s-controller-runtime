package controllers

import (
	"fmt"
	"strings"
	"testing"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/fluxcd/pkg/apis/meta"
	"github.com/fluxcd/pkg/runtime/conditions"
	"github.com/fluxcd/pkg/runtime/patch"
	. "github.com/onsi/gomega"
	talksv1 "github.com/scottrigby/how-to-write-a-reconciler-using-k8s-controller-runtime/cfp/api/v1"
)

func Test_Proposal_Reconcile(t *testing.T) {
	g := NewGomegaWithT(t)

	testCases := []struct {
		name             string
		title            string
		abstract         string
		proposalType     string
		final            bool
		assertConditions []metav1.Condition
		beforeFunc       func(obj *talksv1.Proposal)
		assertFunc       func(obj *talksv1.Proposal, speaker *talksv1.Speaker, assertConditions []metav1.Condition)
	}{
		{
			name:         "test create proposal reconciliation without existing speaker",
			title:        "this is a test proposal",
			abstract:     "this is a test abstract",
			proposalType: "lightning",
			final:        true,
			assertConditions: []metav1.Condition{
				*conditions.TrueCondition(meta.ReconcilingCondition, meta.ProgressingReason, "Reconciling a new generation of the object <generation>"),
				*conditions.FalseCondition(meta.ReadyCondition, meta.FailedReason, "unable to get speaker <namespacedName>: <group> \"<name>\" not found"),
				*conditions.TrueCondition(talksv1.FetchFailedCondition, talksv1.FetchFailedReason, "unable to get speaker <namespacedName>: <group> \"<name>\" not found"),
			},
			beforeFunc: func(obj *talksv1.Proposal) {
				key := client.ObjectKey{Name: obj.Name, Namespace: obj.Namespace}
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

				patch, err := patch.NewHelper(obj, testEnv)
				g.Expect(err).ToNot(HaveOccurred())
				// Set wrong speaker name
				obj.Spec.SpeakerRef = &talksv1.SpeakerRef{
					Name: "unkown",
				}
				g.Expect(patch.Patch(ctx, obj)).To(Succeed())
			},
			assertFunc: func(obj *talksv1.Proposal, _ *talksv1.Speaker, assertConditions []metav1.Condition) {
				key := client.ObjectKey{Name: obj.Name, Namespace: obj.Namespace}
				// Wait for finalizer to be set
				g.Eventually(func() bool {
					if err := testEnv.Get(ctx, key, obj); err != nil {
						return false
					}
					return len(obj.Finalizers) > 0
				}, timeout).Should(BeTrue())

				// Wait for proposal to not be Ready
				g.Eventually(func() bool {
					if err := testEnv.Get(ctx, key, obj); err != nil {
						return false
					}
					if !conditions.IsFalse(obj, meta.ReadyCondition) {
						return false
					}
					return true
				}, timeout).Should(BeTrue())

				for k := range assertConditions {
					assertConditions[k].Message = strings.ReplaceAll(assertConditions[k].Message, "<name>", obj.Spec.SpeakerRef.Name)
					assertConditions[k].Message = strings.ReplaceAll(assertConditions[k].Message, "<generation>", fmt.Sprint(obj.Generation))
					assertConditions[k].Message = strings.ReplaceAll(assertConditions[k].Message, "<namespacedName>", fmt.Sprintf("%s/%s", obj.Namespace, obj.Spec.SpeakerRef.Name))
					assertConditions[k].Message = strings.ReplaceAll(assertConditions[k].Message, "<group>", fmt.Sprintf("Speaker.%s", talksv1.GroupVersion.Group))
				}
				g.Expect(obj.Status.Conditions).To(conditions.MatchConditions(assertConditions))

			},
		},
		{
			name:         "test create proposal reconciliation with existing speaker",
			title:        "this is a test proposal",
			abstract:     "this is a test abstract",
			proposalType: "lightning",
			final:        true,
			assertConditions: []metav1.Condition{
				*conditions.TrueCondition(meta.ReadyCondition, meta.SucceededReason, "reconciled '<name>' successfully"),
			},
			assertFunc: func(obj *talksv1.Proposal, _ *talksv1.Speaker, assertConditions []metav1.Condition) {
				key := client.ObjectKey{Name: obj.Name, Namespace: obj.Namespace}
				// Wait for Proposal to be Ready
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
			name:         "test update proposal reconciliation from draft to final",
			title:        "this is a test proposal",
			abstract:     "this is a test abstract",
			proposalType: "lightning",
			final:        false,
			assertConditions: []metav1.Condition{
				*conditions.TrueCondition(meta.ReadyCondition, meta.SucceededReason, "reconciled '<name>' successfully"),
			},
			assertFunc: func(obj *talksv1.Proposal, _ *talksv1.Speaker, assertConditions []metav1.Condition) {
				key := client.ObjectKey{Name: obj.Name, Namespace: obj.Namespace}
				// Wait for Proposal to be Ready
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

				g.Expect(obj.Status.Submission).To(Equal("draft"))

				patch, err := patch.NewHelper(obj, testEnv)
				g.Expect(err).ToNot(HaveOccurred())
				// Set the status to final
				obj.Spec.Final = true
				g.Expect(patch.Patch(ctx, obj)).To(Succeed())

				// Wait for Proposal to be Ready
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
						obj.Status.Submission == talksv1.ProposalStateFinal
				}, timeout).Should(BeTrue())

				for k := range assertConditions {
					assertConditions[k].Message = strings.ReplaceAll(assertConditions[k].Message, "<name>", obj.Name)
				}
				g.Expect(obj.Status.Conditions).To(conditions.MatchConditions(assertConditions))

			},
		},
		{
			name:         "test delete proposal reconciliation",
			title:        "this is a test proposal",
			abstract:     "this is a test abstract",
			proposalType: "lightning",
			final:        true,
			assertFunc: func(obj *talksv1.Proposal, _ *talksv1.Speaker, assertConditions []metav1.Condition) {
				key := client.ObjectKey{Name: obj.Name, Namespace: obj.Namespace}
				// Wait for Proposal to be Ready
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

				// Delete Proposal
				g.Expect(testEnv.Delete(ctx, obj)).To(Succeed())

				// Wait for Proposal to be deleted
				g.Eventually(func() bool {
					if err := testEnv.Get(ctx, key, obj); err != nil {
						return apierrors.IsNotFound(err)
					}
					return false
				}, timeout).Should(BeTrue())
			},
		},
		{
			name:         "test delete Speaker, expect a proposal status change",
			title:        "this is a test proposal",
			abstract:     "this is a test abstract",
			proposalType: "lightning",
			final:        true,
			assertConditions: []metav1.Condition{
				*conditions.FalseCondition(meta.ReadyCondition, meta.FailedReason, "unable to get speaker <namespacedName>: <group> \"<name>\" not found"),
				*conditions.TrueCondition(talksv1.FetchFailedCondition, talksv1.FetchFailedReason, "unable to get speaker <namespacedName>: <group> \"<name>\" not found"),
			},
			assertFunc: func(obj *talksv1.Proposal, speaker *talksv1.Speaker, assertConditions []metav1.Condition) {
				speakerKey := client.ObjectKey{Name: speaker.Name, Namespace: speaker.Namespace}
				// Wait for Speaker to be Ready
				g.Eventually(func() bool {
					if err := testEnv.Get(ctx, speakerKey, speaker); err != nil {
						return false
					}
					if !conditions.IsReady(speaker) {
						return false
					}
					readyCondition := conditions.Get(speaker, meta.ReadyCondition)
					return speaker.Generation == readyCondition.ObservedGeneration &&
						speaker.Generation == speaker.Status.ObservedGeneration
				}, timeout).Should(BeTrue())

				key := client.ObjectKey{Name: obj.Name, Namespace: obj.Namespace}
				// Wait for Proposal to be Ready
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
				g.Expect(testEnv.Delete(ctx, speaker)).To(Succeed())

				// Wait for <speakerto be deleted
				g.Eventually(func() bool {
					if err := testEnv.Get(ctx, speakerKey, obj); err != nil {
						return apierrors.IsNotFound(err)
					}
					return false
				}, timeout).Should(BeTrue())

				// Wait for proposal to not be Ready
				g.Eventually(func() bool {
					if err := testEnv.Get(ctx, key, obj); err != nil {
						return false
					}
					if !conditions.IsFalse(obj, meta.ReadyCondition) {
						return false
					}
					readyCondition := conditions.Get(speaker, meta.ReadyCondition)
					return obj.Generation == readyCondition.ObservedGeneration &&
						obj.Generation == obj.Status.ObservedGeneration
				}, timeout).Should(BeTrue())
				for k := range assertConditions {
					assertConditions[k].Message = strings.ReplaceAll(assertConditions[k].Message, "<name>", obj.Spec.SpeakerRef.Name)
					assertConditions[k].Message = strings.ReplaceAll(assertConditions[k].Message, "<namespacedName>", fmt.Sprintf("%s/%s", obj.Namespace, obj.Spec.SpeakerRef.Name))
					assertConditions[k].Message = strings.ReplaceAll(assertConditions[k].Message, "<group>", fmt.Sprintf("Speaker.%s", talksv1.GroupVersion.Group))
				}
				g.Expect(obj.Status.Conditions).To(conditions.MatchConditions(assertConditions))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ns, err := testEnv.CreateNamespace(ctx, "proposal-ns")
			g.Expect(err).NotTo(HaveOccurred())
			defer func() {
				g.Expect(testEnv.Delete(ctx, ns)).To(Succeed())
			}()

			speaker := &talksv1.Speaker{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "speaker",
					Namespace: ns.Name,
				},
				Spec: talksv1.SpeakerSpec{
					Name:  "test",
					Bio:   "test speaker",
					Email: "test.speaker@gmail.com",
				},
			}

			err = testEnv.CreateAndWait(ctx, speaker)
			g.Expect(err).NotTo(HaveOccurred())

			key := client.ObjectKey{Name: speaker.Name, Namespace: speaker.Namespace}
			// Wait for Speaker to be Ready
			g.Eventually(func() bool {
				if err := testEnv.Get(ctx, key, speaker); err != nil {
					return false
				}
				if !conditions.IsReady(speaker) {
					return false
				}
				readyCondition := conditions.Get(speaker, meta.ReadyCondition)
				return speaker.Generation == readyCondition.ObservedGeneration &&
					speaker.Generation == speaker.Status.ObservedGeneration
			}, timeout).Should(BeTrue())

			obj := &talksv1.Proposal{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "proposal",
					Namespace: ns.Name,
				},
				Spec: talksv1.ProposalSpec{
					Title:    tc.title,
					Abstract: tc.abstract,
					Type:     tc.proposalType,
					Final:    tc.final,
					SpeakerRef: &talksv1.SpeakerRef{
						Name: speaker.Name,
					},
				},
			}

			err = testEnv.CreateAndWait(ctx, obj)
			g.Expect(err).NotTo(HaveOccurred())

			if tc.beforeFunc != nil {
				tc.beforeFunc(obj)
			}

			tc.assertFunc(obj, speaker, tc.assertConditions)
		})
	}
}
