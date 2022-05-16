package main

import (
    "flag"
	"strconv"
    "fmt"
    "os"
    "path/filepath"
//    "time"
    "context"
	"strings"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

func main() {
    var kubeconfig *string
    var ctx, _ = context.WithCancel(context.Background())
    if home := homeDir(); home != "" {
        kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
    } else {
        kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
    }
    flag.Parse()

    // uses the current context in kubeconfig
    config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {
        panic(err.Error())
    }

    // creates the clientset
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err.Error())
    }

    namespaceList, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
    if err != nil {
        panic(err)
    }

    for _, namespace := range namespaceList.Items {
		if strings.Contains(namespace.Name, "kube-system") || strings.Contains(namespace.Name, "monitoring") {
			continue
		} else {	
			serviceList, err := clientset.CoreV1().Services(namespace.Name).List(ctx, metav1.ListOptions{})
			if err != nil {
				fmt.Println(err)
			}

			for _, service := range serviceList.Items {
				for _, p := range service.Spec.Ports  {
					port_int64 := int64(p.Port)
					portStr := strconv.FormatInt(port_int64, 10)
					if strings.Contains(service.Name, "glusterfs") == false {
	                    fmt.Println(service.Name + ".namespace.svc.cluster.local"+ ":" + portStr)
	                    break
					}
				}
			}
		}
	}
}


func homeDir() string {
    if h := os.Getenv("HOME"); h != "" {
        return h
    }
    return os.Getenv("USERPROFILE") // windows
}
