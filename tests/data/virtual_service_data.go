package data

import (
	"github.com/kiali/kiali/kubernetes"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateEmptyVirtualService(name string, namespace string, hosts []string) kubernetes.IstioObject {
	return (&kubernetes.VirtualService{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			ClusterName: "svc.cluster.local",
		},
		Spec: map[string]interface{}{
			"hosts": hosts,
		},
	}).DeepCopyIstioObject()
}

// TODO Naming etc
func CreateVirtualService() kubernetes.IstioObject {
	return AddRoutesToVirtualService("http", CreateRoute("reviews", "v1", -1),
		AddRoutesToVirtualService("tcp", CreateRoute("reviews", "v1", -1),
			CreateEmptyVirtualService("reviews", "test", []string{"reviews"}),
		),
	)
}

func AddRoutesToVirtualService(routeType string, route map[string]interface{}, vs kubernetes.IstioObject) kubernetes.IstioObject {
	if routeTypeExists, found := vs.GetSpec()[routeType]; found {
		if routeTypeCasted, ok := routeTypeExists.([]interface{}); ok {
			if routeElement, ok := routeTypeCasted[0].(map[string]interface{}); ok {
				if routeValue, found := routeElement["route"]; found {
					if routeOneMoreCast, ok := routeValue.([]interface{}); ok {
						routeOneMoreCast = append(routeOneMoreCast, route)
						routeElement["route"] = routeOneMoreCast
					}
				}
			}
		}
	} else {
		vs.GetSpec()[routeType] = []interface{}{
			map[string]interface{}{
				"route": []interface{}{route},
			},
		}
	}
	return vs
}

func CreateRoute(host string, subset string, weight int64) map[string]interface{} {
	route := make(map[string]interface{})
	route["destination"] = map[string]interface{}{
		"host":   host,
		"subset": subset,
	}
	if weight >= 0 {
		route["weight"] = uint64(weight) // Weight is supposed to be between [0,100] so this is safe
	}
	return route
}

// Example from https://istio.io/docs/reference/config/istio.networking.v1alpha3/#Destination
func CreateVirtualServiceWithServiceEntryTarget() kubernetes.IstioObject {
	vs := CreateEmptyVirtualService("my-wiki-rule", "wikipedia", []string{"wikipedia.org"})
	vs.GetSpec()["http"] = []interface{}{
		map[string]interface{}{
			"timeout": "5s",
			"route": []interface{}{
				map[string]interface{}{
					"destination": map[string]interface{}{
						"host": "wikipedia.org",
					},
				},
			},
		},
	}
	return vs
}

func CreateExternalServiceEntry() kubernetes.IstioObject {
	return (&kubernetes.ServiceEntry{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "external-svc-wikipedia",
			Namespace: "wikipedia",
		},
		Spec: map[string]interface{}{
			"hosts": []interface{}{
				"wikipedia.org",
			},
			"location": "MESH_EXTERNAL",
			"ports": map[string]interface{}{
				"number":   uint64(80),
				"name":     "example-http",
				"protocol": "HTTP",
			},
			"resolution": "DNS",
		},
	}).DeepCopyIstioObject()
}
