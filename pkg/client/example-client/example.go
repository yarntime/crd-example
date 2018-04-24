package example_client

import (
	"github.com/yarntime/crd-example/pkg/types"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

type ExampleGetter interface {
	Examples(namespace string) ExampleInterface
}

type ExampleInterface interface {
	Create(*types.Example) (*types.Example, error)
	Update(*types.Example) (*types.Example, error)
	UpdateStatus(*types.Example) (*types.Example, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*types.Example, error)
	List(opts meta_v1.ListOptions) (*types.ExampleList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
}

type examples struct {
	client rest.Interface
	ns     string
}

func newExamples(c *ExampleClient, namespace string) *examples {
	return &examples{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

func (c *examples) Create(monitoredTarget *types.Example) (result *types.Example, err error) {
	result = &types.Example{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource(ResourceKind).
		Body(monitoredTarget).
		Do().
		Into(result)
	return
}

func (c *examples) Update(monitoredTarget *types.Example) (result *types.Example, err error) {
	result = &types.Example{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource(ResourceKind).
		Name(monitoredTarget.Name).
		Body(monitoredTarget).
		Do().
		Into(result)
	return
}

func (c *examples) UpdateStatus(monitoredTarget *types.Example) (result *types.Example, err error) {
	result = &types.Example{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource(ResourceKind).
		Name(monitoredTarget.Name).
		SubResource("status").
		Body(monitoredTarget).
		Do().
		Into(result)
	return
}

func (c *examples) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource(ResourceKind).
		Name(name).
		Body(options).
		Do().
		Error()
}

func (c *examples) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource(ResourceKind).
		VersionedParams(&listOptions, ParameterCodec).
		Body(options).
		Do().
		Error()
}

func (c *examples) Get(name string, options meta_v1.GetOptions) (result *types.Example, err error) {
	result = &types.Example{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource(ResourceKind).
		Name(name).
		VersionedParams(&options, ParameterCodec).
		Do().
		Into(result)
	return
}

func (c *examples) List(opts meta_v1.ListOptions) (result *types.ExampleList, err error) {
	result = &types.ExampleList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource(ResourceKind).
		VersionedParams(&opts, ParameterCodec).
		Do().
		Into(result)
	return
}

func (c *examples) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource(ResourceKind).
		VersionedParams(&opts, ParameterCodec).
		Watch()
}
