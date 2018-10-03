package cluster

import (
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type namespaceClient struct {
	dynamicClient dynamic.Interface
}

var _ NamespaceInterface = (*namespaceClient)(nil)

func newNamespaceClient(dynamicClient dynamic.Interface) *namespaceClient {
	return &namespaceClient{
		dynamicClient: dynamicClient,
	}
}

func (n *namespaceClient) Names() ([]string, error) {
	namespaces, err := Namespaces(n.dynamicClient)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, namespace := range namespaces {
		names = append(names, namespace.GetName())
	}

	return names, nil
}

// Namespaces returns available namespaces.
func Namespaces(dc dynamic.Interface) ([]corev1.Namespace, error) {
	res := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "namespaces",
	}

	nri := dc.Resource(res)

	list, err := nri.List(metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "list namespaces")
	}

	var nsList corev1.NamespaceList
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(list.UnstructuredContent(), &nsList)
	if err != nil {
		return nil, errors.Wrap(err, "convert object to namespace list")
	}

	return nsList.Items, nil
}