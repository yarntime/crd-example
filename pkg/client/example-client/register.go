package example_client

import (
	"reflect"

	"github.com/yarntime/crd-example/pkg/types"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"time"
)

const (
	GroupName    = "yarntime.io"
	ResourceKind = "example"
	GroupVersion = "v1"
	FullCRDName  = ResourceKind + "." + GroupName
)

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: GroupVersion}

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	SchemeBuilder      runtime.SchemeBuilder
	LocalSchemeBuilder = &SchemeBuilder
)

func init() {
	SchemeBuilder.Register(addKnownTypes)
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&types.Example{},
		&types.ExampleList{},
	)

	meta_v1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

func RegisterExample(config *rest.Config) error {
	clientset, err := apiextcs.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	crd := &apiextv1beta1.CustomResourceDefinition{
		ObjectMeta: meta_v1.ObjectMeta{Name: FullCRDName},
		Spec: apiextv1beta1.CustomResourceDefinitionSpec{
			Group:   GroupName,
			Version: GroupVersion,
			Scope:   apiextv1beta1.NamespaceScoped,
			Names: apiextv1beta1.CustomResourceDefinitionNames{
				Plural:     ResourceKind,
				Singular:   "example",
				ShortNames: []string{"example"},
				Kind:       reflect.TypeOf(types.Example{}).Name(),
			},
		},
	}

	_, err = clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if err != nil && apierrors.IsAlreadyExists(err) {
		return nil
	}

	retryCount := 5
	for retryCount > 0 {
		_, err = clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Get(FullCRDName, meta_v1.GetOptions{})
		if err == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
		retryCount--
	}

	return err
}
