package compare

import (
	"fmt"
	utils "github.com/RHsyseng/operator-utils/pkg/resource/test"
	oappsv1 "github.com/openshift/api/apps/v1"
	obuildv1 "github.com/openshift/api/build/v1"
	routev1 "github.com/openshift/api/route/v1"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"reflect"
	"testing"
)

func TestCompareRoutes(t *testing.T) {
	routes := utils.GetRoutes(2)
	routes[0].Status = routev1.RouteStatus{
		Ingress: []routev1.RouteIngress{
			{
				Host: "localhost",
			},
		},
	}
	routes[1].Name = routes[0].Name

	assert.False(t, reflect.DeepEqual(routes[0], routes[1]), "Inconsequential differences between two routes should make equality test fail")
	assert.True(t, deepEquals(&routes[0], &routes[1]), "Expected resources to be deemed equal")
	assert.True(t, equalRoutes(&routes[0], &routes[1]), "Expected resources to be deemed equal based on route comparator")
}

func TestCompareServices(t *testing.T) {
	services := utils.GetServices(2)
	services[0].Status = corev1.ServiceStatus{
		LoadBalancer: corev1.LoadBalancerStatus{
			Ingress: []corev1.LoadBalancerIngress{
				{
					IP:       "127.0.0.1",
					Hostname: "localhost",
				},
			},
		},
	}
	services[1].Name = services[0].Name

	assert.False(t, reflect.DeepEqual(services[0], services[1]), "Inconsequential differences between two services should make equality test fail")
	assert.True(t, deepEquals(&services[0], &services[1]), "Expected resources to be deemed equal")
	assert.True(t, equalServices(&services[0], &services[1]), "Expected resources to be deemed equal based on service comparator")
}

func TestCompareDeploymentConfigs(t *testing.T) {
	dcs := utils.GetDeploymentConfigs(2)
	dcs[1].Name = dcs[0].Name
	dcs[1].Status = oappsv1.DeploymentConfigStatus{
		ReadyReplicas: 1,
	}

	assert.False(t, reflect.DeepEqual(dcs[0], dcs[1]), "Inconsequential differences between two DCs should make equality test fail")
	assert.True(t, deepEquals(&dcs[0], &dcs[1]), "Expected resources to be deemed equal")
	assert.True(t, equalDeploymentConfigs(&dcs[0], &dcs[1]), "Expected resources to be deemed equal based on DC comparator")
}

func TestCompareEmptyAnnotations(t *testing.T) {
	routes := utils.GetRoutes(2)
	routes[1].Name = routes[0].Name
	routes[0].Annotations = make(map[string]string)
	routes[0].Annotations["openshift.io/host.generated"] = "true"
	routes[1].Annotations = nil
	assert.True(t, equalRoutes(&routes[0], &routes[1]), "Routes should be considered equal")
}

func TestCompareDeploymentConfigLastTriggeredImage(t *testing.T) {
	dcs := utils.GetDeploymentConfigs(2)
	dcs[1].Name = dcs[0].Name
	dcs[0].Spec.Triggers = []oappsv1.DeploymentTriggerPolicy{
		{
			ImageChangeParams: &oappsv1.DeploymentTriggerImageChangeParams{
				Automatic:          false,
				ContainerNames:     nil,
				From:               corev1.ObjectReference{},
				LastTriggeredImage: "some generated value",
			},
		},
	}
	dcs[1].Spec.Triggers = []oappsv1.DeploymentTriggerPolicy{
		{
			ImageChangeParams: &oappsv1.DeploymentTriggerImageChangeParams{
				Automatic:      false,
				ContainerNames: nil,
				From:           corev1.ObjectReference{},
			},
		},
	}
	assert.True(t, equalDeploymentConfigs(&dcs[0], &dcs[1]), "Expected resources to be deemed equal based on DC comparator")
}

func TestCompareDeploymentConfigImageChange(t *testing.T) {
	dcs := utils.GetDeploymentConfigs(2)
	dcs[1].Name = dcs[0].Name
	dcs[0].Spec.Triggers = []oappsv1.DeploymentTriggerPolicy{
		{
			ImageChangeParams: &oappsv1.DeploymentTriggerImageChangeParams{
				Automatic: false,
				ContainerNames: []string{
					"container1",
					"container2",
				},
				From: corev1.ObjectReference{
					Kind:      "ImageStreamTag",
					Namespace: "namespace",
					Name:      "image",
				},
			},
		},
	}
	dcs[0].Spec.Template.Spec.Containers = []corev1.Container{
		{
			Name:  "container1",
			Image: "some generated value",
		},
	}
	dcs[1].Spec.Triggers = []oappsv1.DeploymentTriggerPolicy{
		{
			ImageChangeParams: &oappsv1.DeploymentTriggerImageChangeParams{
				Automatic: false,
				ContainerNames: []string{
					"container1",
					"container2",
				},
				From: corev1.ObjectReference{
					Kind:      "ImageStreamTag",
					Namespace: "namespace",
					Name:      "image",
				},
			},
		},
	}
	dcs[1].Spec.Template.Spec.Containers = []corev1.Container{
		{
			Name:  "container1",
			Image: "image",
		},
	}
	assert.True(t, equalDeploymentConfigs(&dcs[0], &dcs[1]), "Expected resources to be deemed equal based on DC comparator")
}

func TestCompareBuildConfigWebHooks(t *testing.T) {
	bcs := utils.GetBuildConfigs(2)
	bcs[1].Name = bcs[0].Name
	bcs[0].Spec.RunPolicy = obuildv1.BuildRunPolicySerial
	bcs[0].Spec.Triggers = []obuildv1.BuildTriggerPolicy{
		{
			GitLabWebHook: &obuildv1.WebHookTrigger{
				AllowEnv:        false,
				SecretReference: &obuildv1.SecretLocalReference{Name: "dafsaf"},
			},
		},
	}
	bcs[1].Spec.Triggers = []obuildv1.BuildTriggerPolicy{
		{
			GitLabWebHook: &obuildv1.WebHookTrigger{
				AllowEnv:        false,
				SecretReference: &obuildv1.SecretLocalReference{Name: "eqwrer"},
			},
		},
	}
	assert.True(t, equalBuildConfigs(&bcs[0], &bcs[1]), "Expected resources to be deemed equal based on BC comparator")
}

func TestCompareDeployments(t *testing.T) {
	deployments := utils.GetDeployments(2)
	deployments[1].Name = deployments[0].Name
	deployments[1].Status = appsv1.DeploymentStatus{
		ReadyReplicas: 1,
	}

	assert.False(t, reflect.DeepEqual(deployments[0], deployments[1]), "Inconsequential differences between two Deployments should make equality test fail")
	assert.True(t, deepEquals(&deployments[0], &deployments[1]), "Expected resources to be deemed equal")
	assert.True(t, equalDeployment(&deployments[0], &deployments[1]), "Expected resources to be deemed equal based on Deployment comparator")
}

func TestCompareDeploymentLastTriggeredImage(t *testing.T) {
	imageTriggersFormat := "[{\"from\":{\"kind\":\"ImageStreamTag\",\"name\":\"%s\"},\"fieldPath\":\"spec.template.spec.containers[?(@.name==\\\"%s\\\")].image\"}]"
	deployments := utils.GetDeployments(2)
	deployments[1].Name = deployments[0].Name
	deployments[0].Annotations = map[string]string{
		imageTriggersAnnotation: fmt.Sprintf(imageTriggersFormat, "my-image", "my-container"),
	}
	deployments[1].Annotations = map[string]string{
		imageTriggersAnnotation: fmt.Sprintf(imageTriggersFormat, "my-image", "my-container"),
	}
	deployments[0].Spec.Template.Spec.Containers = []corev1.Container{
		{Name: "my-container", Image: "some generated value"},
	}
	deployments[1].Spec.Template.Spec.Containers = []corev1.Container{
		{Name: "my-container", Image: "quay.io/namespace/image:tag"},
	}
	assert.True(t, equalDeployment(&deployments[0], &deployments[1]), "Expected resources to be deemed equal based on deployment comparator")
}

func TestCompareDeploymentGenerateValue(t *testing.T) {
	deployments := utils.GetDeployments(2)
	deployments[1].Name = deployments[0].Name
	deployments[0].Spec.Template.Spec.DNSPolicy = corev1.DNSClusterFirst

	assert.True(t, equalDeployment(&deployments[0], &deployments[1]), "Expected resources to be deemed equal based on deployment comparator")
}
