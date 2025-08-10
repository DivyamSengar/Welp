package services_test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func TestFrontendHealth(t *testing.T) {
	testPodHealth("frontend", t)
}

func TestDetailHealth(t *testing.T) {
	testPodHealth("detail", t)
}

func TestReservationHealth(t *testing.T) {
	testPodHealth("reservation", t)
}

func TestReviewHealth(t *testing.T) {
	testPodHealth("review", t)
}

func testPodHealth(podNamePrefix string, t *testing.T) {
	// Step 1: setup kubernetes client
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		t.Fatalf("Error building kubeconfig: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Error creating Kubernetes client: %v", err)
	}

	// Step 2: retrieve the pod
	// podName should be the output name from `k get po`, i.e. frontend-7668fcd98f-69flj
	namespace := "default"
	// pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		// t.Fatalf("Error retrieving pods: %v", err)
		t.Fatalf("Error retrieving pods: %v", err)
	}
	var pod *corev1.Pod
	hasPod := false
	for _, currPod := range pods.Items {
		if strings.HasPrefix(currPod.Name, podNamePrefix) {
			pod = &currPod
			hasPod = true
			break
		}
	}
	if !hasPod {
		t.Fatalf("Pod with prefix %s does not exist", podNamePrefix)
	}

	// Step 3: check status is 'ready' for pod
	for _, containerStatus := range pod.Status.ContainerStatuses {
		containerName := containerStatus.Name
		if !containerStatus.Ready {
			t.Errorf("Container %s is not ready", containerName)
		}
		if containerStatus.State.Running == nil {
			t.Errorf("Container %s is not running", containerName)
		}
	}
}
