package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// terraform wrapper for external resource
// due to a bug on []strings, make(map[string]string) is the only format usable with external resource,
// so structure may look strange
func main() {
	var input map[string]string
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decode stdin: %v\n", err)
		os.Exit(1)
	}

	nsRaw, ok := input["ns"]
	if !ok || nsRaw == "" {
		fmt.Fprintln(os.Stderr, "Missing 'ns' input")
		os.Exit(1)
	}

	namespaces := strings.Split(nsRaw, ",")
	kubeconfigPath := os.Getenv("KUBECONFIG")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading kubeconfig: %v\n", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating clientset: %v\n", err)
		os.Exit(1)
	}

	serviceNames := make(map[string]string)

	for _, ns := range namespaces {
		svcs, err := clientset.CoreV1().Services(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing services in namespace %s: %v\n", ns, err)
			os.Exit(1)
		}
		for _, svc := range svcs.Items {
			serviceNames[fmt.Sprintf("%s.%s.svc", svc.Name, svc.Namespace)] = "ok"
		}
	}

	out, _ := json.Marshal(serviceNames)
	fmt.Println(string(out))

}
