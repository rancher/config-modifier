/*
Copyright 2020 Rancher Labs, Inc.

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

// Code generated by main. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	v1 "github.com/rancher/config-modifier/pkg/apis/some.api.group/v1"
	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type FooHandler func(string, *v1.Foo) (*v1.Foo, error)

type FooController interface {
	generic.ControllerMeta
	FooClient

	OnChange(ctx context.Context, name string, sync FooHandler)
	OnRemove(ctx context.Context, name string, sync FooHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() FooCache
}

type FooClient interface {
	Create(*v1.Foo) (*v1.Foo, error)
	Update(*v1.Foo) (*v1.Foo, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.Foo, error)
	List(namespace string, opts metav1.ListOptions) (*v1.FooList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Foo, err error)
}

type FooCache interface {
	Get(namespace, name string) (*v1.Foo, error)
	List(namespace string, selector labels.Selector) ([]*v1.Foo, error)

	AddIndexer(indexName string, indexer FooIndexer)
	GetByIndex(indexName, key string) ([]*v1.Foo, error)
}

type FooIndexer func(obj *v1.Foo) ([]string, error)

type fooController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewFooController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) FooController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &fooController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromFooHandlerToHandler(sync FooHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.Foo
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.Foo))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *fooController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.Foo))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateFooDeepCopyOnChange(client FooClient, obj *v1.Foo, handler func(obj *v1.Foo) (*v1.Foo, error)) (*v1.Foo, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *fooController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *fooController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *fooController) OnChange(ctx context.Context, name string, sync FooHandler) {
	c.AddGenericHandler(ctx, name, FromFooHandlerToHandler(sync))
}

func (c *fooController) OnRemove(ctx context.Context, name string, sync FooHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromFooHandlerToHandler(sync)))
}

func (c *fooController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *fooController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *fooController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *fooController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *fooController) Cache() FooCache {
	return &fooCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *fooController) Create(obj *v1.Foo) (*v1.Foo, error) {
	result := &v1.Foo{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *fooController) Update(obj *v1.Foo) (*v1.Foo, error) {
	result := &v1.Foo{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *fooController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *fooController) Get(namespace, name string, options metav1.GetOptions) (*v1.Foo, error) {
	result := &v1.Foo{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *fooController) List(namespace string, opts metav1.ListOptions) (*v1.FooList, error) {
	result := &v1.FooList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *fooController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *fooController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v1.Foo, error) {
	result := &v1.Foo{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type fooCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *fooCache) Get(namespace, name string) (*v1.Foo, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1.Foo), nil
}

func (c *fooCache) List(namespace string, selector labels.Selector) (ret []*v1.Foo, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Foo))
	})

	return ret, err
}

func (c *fooCache) AddIndexer(indexName string, indexer FooIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.Foo))
		},
	}))
}

func (c *fooCache) GetByIndex(indexName, key string) (result []*v1.Foo, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1.Foo, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1.Foo))
	}
	return result, nil
}