package main

import (
	"context"
	"fmt"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	//get our kubeconfig file
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")

	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		panic(err)
	}

	list, err := clientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, record := range list.Items {
		fmt.Println(record.Name)
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-pod",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx-container",
					Image: "nginx",
				},
			},
		},
	}

	//create pod client
	podClient := clientSet.CoreV1().Pods("default")
	result, err := podClient.Create(context.Background(), pod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println(result.Name)

	//create deployment
	replicas := int32(2)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "web",
							Image: "nginx",
						},
					},
				},
			},
		},
	}

	//create deployment client
	deploymentClient := clientSet.AppsV1().Deployments("default")
	deployResult, err := deploymentClient.Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	//todo update deployment
	//deploymentClient.Update()
	//todo delete deployment
	//deploymentClient.Delete()
	fmt.Println(deployResult.Name)

}
