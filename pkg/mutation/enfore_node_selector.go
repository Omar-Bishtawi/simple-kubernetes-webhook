package mutation

import (
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

type enforeNodeSelector struct {
	Logger logrus.FieldLogger
}

var _ podMutator = (*enforeNodeSelector)(nil)

func (se enforeNodeSelector) Name() string {
	return "enfore_node_selector"
}

func (se enforeNodeSelector) Mutate(pod *corev1.Pod) (*corev1.Pod, error) {
	se.Logger = se.Logger.WithField("mutation", se.Name())
	mpod := pod.DeepCopy()

	// build out env var slice
	nodeAffinity := corev1.NodeAffinity{
		RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
			NodeSelectorTerms: []corev1.NodeSelectorTerm{
				{
					MatchExpressions: []corev1.NodeSelectorRequirement{
						{
							Key:      "workload",
							Operator: corev1.NodeSelectorOpIn,
							Values:   []string{"normal-arm-karpenter"},
						},
					},
				},
			},
		},
	}

	namespacesToEnforce := []string{"test"}

	if isValidNamespace(mpod.Namespace, namespacesToEnforce) {
		if mpod.Spec.Affinity == nil {
			mpod.Spec.Affinity = &corev1.Affinity{}
		}
		mpod.Spec.Affinity.NodeAffinity = &nodeAffinity
		se.Logger.Debugf("pod node affinity enforced %s", nodeAffinity)
	}

	return mpod, nil
}

func isValidNamespace(ns string, namespaces []string) bool {
	for _, n := range namespaces {
		if ns == n {
			return true
		}
	}
	return false
}
