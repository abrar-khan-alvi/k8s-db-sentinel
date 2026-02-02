package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var config *rest.Config
	var err error

	// 1. Check if running locally (looking for file)
	userHome, _ := os.UserHomeDir()
	kubeconfigPath := filepath.Join(userHome, ".kube", "config")

	// We use 'flag' to allow overrides, which fixes the "unused import" error
	kubeconfig := flag.String("kubeconfig", kubeconfigPath, "(optional) absolute path to the kubeconfig file")
	flag.Parse()

	if _, err := os.Stat(*kubeconfig); err == nil {
		fmt.Println("[INIT] Running locally. Using .kube/config")
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	} else {
		fmt.Println("[INIT] Running inside cluster. Using ServiceAccount.")
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		panic(err.Error())
	}

	// 2. Create the Client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("--- DB Sentinel Started. Monitoring 'my-postgres' ---")

	// 3. The Infinite Loop
	for {
		_, err := clientset.CoreV1().Pods("default").Get(context.TODO(), "my-postgres", metav1.GetOptions{})

		if err != nil {
			if errors.IsNotFound(err) {
				fmt.Println("[ALERT] Database is MISSING! Initiating Recovery... 🚑")
				createPod(clientset)
			} else {
				fmt.Printf("[ERROR] Unexpected error: %s\n", err.Error())
			}
		} else {
			fmt.Println("[OK] Database is healthy. ✅")
		}
		time.Sleep(5 * time.Second)
	}
}

func createPod(clientset *kubernetes.Clientset) {
	newPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "my-postgres"},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "postgres",
					Image: "postgres:13",
					Env: []corev1.EnvVar{
						{Name: "POSTGRES_PASSWORD", Value: "secret"},
					},
				},
			},
		},
	}
	fmt.Println("[HEALER] Creating new Pod...")
	_, err := clientset.CoreV1().Pods("default").Create(context.TODO(), newPod, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("[FAIL] Failed to heal: %s\n", err.Error())
	} else {
		fmt.Println("[SUCCESS] New database pod created! Recovery complete. 🚀")
	}
}
