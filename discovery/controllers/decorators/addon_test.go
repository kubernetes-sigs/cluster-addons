package decorators

import (
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	k8sscheme "k8s.io/client-go/kubernetes/scheme"

	discoveryv1alpha1 "sigs.k8s.io/cluster-addons/discovery/api/v1alpha1"
)

func TestAddonNames(t *testing.T) {
	type args struct {
		labels map[string]string
	}
	type results struct {
		names []types.NamespacedName
	}

	tests := []struct {
		description string
		args        args
		results     results
	}{
		{
			description: "SingleAddon",
			args: args{
				labels: map[string]string{
					ComponentLabelKeyPrefix + "lobster": "",
				},
			},
			results: results{
				names: []types.NamespacedName{
					{Name: "lobster"},
				},
			},
		},
		{
			description: "MultipleAddons",
			args: args{
				labels: map[string]string{
					ComponentLabelKeyPrefix + "lobster": "",
					ComponentLabelKeyPrefix + "cod":     "",
				},
			},
			results: results{
				names: []types.NamespacedName{
					{Name: "lobster"},
					{Name: "cod"},
				},
			},
		},
		{
			description: "NoAddons",
			args: args{
				labels: map[string]string{
					"robot": "whirs_and_clicks",
				},
			},
			results: results{
				names: nil,
			},
		},
		{
			description: "NoLabels",
			args: args{
				labels: nil,
			},
			results: results{
				names: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			require.ElementsMatch(t, tt.results.names, AddonNames(tt.args.labels))
		})
	}
}

func TestAddComponents(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, k8sscheme.AddToScheme(scheme))

	type fields struct {
		addon *discoveryv1alpha1.Addon
	}
	type args struct {
		components []runtime.Object
	}
	type results struct {
		addon *discoveryv1alpha1.Addon
		err   error
	}

	tests := []struct {
		description string
		fields      fields
		args        args
		results     results
	}{
		{
			description: "Empty/ComponentsAdded",
			fields: fields{
				addon: func() *discoveryv1alpha1.Addon {
					addon := &discoveryv1alpha1.Addon{}
					addon.SetName("puffin")

					return addon
				}(),
			},
			args: args{
				components: []runtime.Object{
					func() runtime.Object {
						namespace := &corev1.Namespace{}
						namespace.SetName("atlantic")
						namespace.SetLabels(map[string]string{
							ComponentLabelKeyPrefix + "puffin": "",
						})

						return namespace
					}(),
					func() runtime.Object {
						pod := &corev1.Pod{}
						pod.SetNamespace("atlantic")
						pod.SetName("puffin")
						pod.SetLabels(map[string]string{
							ComponentLabelKeyPrefix + "puffin": "",
						})
						pod.Status.Conditions = []corev1.PodCondition{
							{
								Type:   corev1.PodReady,
								Status: corev1.ConditionTrue,
							},
						}

						return pod
					}(),
				},
			},
			results: results{
				addon: func() *discoveryv1alpha1.Addon {
					addon := &discoveryv1alpha1.Addon{}
					addon.SetName("puffin")
					addon.Status.Components = &discoveryv1alpha1.Components{
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									Key:      ComponentLabelKeyPrefix + addon.GetName(),
									Operator: metav1.LabelSelectorOpExists,
								},
							},
						},
					}
					addon.Status.Components.Refs = []discoveryv1alpha1.RichReference{
						{
							ObjectReference: &corev1.ObjectReference{
								APIVersion: "v1",
								Kind:       "Namespace",
								Name:       "atlantic",
							},
						},
						{
							ObjectReference: &corev1.ObjectReference{
								APIVersion: "v1",
								Kind:       "Pod",
								Namespace:  "atlantic",
								Name:       "puffin",
							},
							Conditions: []discoveryv1alpha1.Condition{
								{
									Type:   discoveryv1alpha1.ConditionType(corev1.PodReady),
									Status: corev1.ConditionTrue,
								},
							},
						},
					}

					return addon
				}(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			addon := &Addon{
				Addon:  tt.fields.addon,
				scheme: scheme,
			}
			err := addon.AddComponents(tt.args.components...)
			require.Equal(t, tt.results.err, err)
			require.Equal(t, tt.results.addon, addon.Addon)
		})
	}
}
