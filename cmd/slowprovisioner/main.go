package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/gman0/slow-provisioner/pkg/slow"
	"github.com/golang/glog"
	"github.com/kubernetes-incubator/external-storage/lib/controller"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig      = flag.String("kubeconfig", "", "Path to a kube config. Only required if out-of-cluster.")
	provisionerName = flag.String("provisioner", "slow-provisioner", "Name of the provisioner. The provisioner will only provision volumes for claims that request a StorageClass with a provisioner field set equal to this name.")
	nodePluginName  = flag.String("nodeplugin", "csi-slowplugin", "Name of the CSI node plugin")

	defaultDelay = flag.String("defaultdelay", "", "default delay applied to all provisioner operations")
	delay        = flag.String("delay", "", "delay settings for individual provisioner operations, e.g. Provision=20..50,Delete=inf")
)

func main() {
	flag.Parse()
	flag.Set("logtostderr", "true")

	var delayOptions slow.DelayOptions

	if *defaultDelay != "" {
		if dr, err := parseDelay(*defaultDelay); err != nil {
			glog.Fatalf("invalid format in default delay: %v", err)
		} else {
			setDefaultDelay(&delayOptions, dr)
		}
	}

	if *delay != "" {
		if err := parseDelayArgs(&delayOptions, *delay); err != nil {
			glog.Fatalf("invalid delay: %v", err)
		}
	}

	// Create an InClusterConfig and use it to create a client for the controller
	// to use to communicate with Kubernetes
	config, err := buildConfig(*kubeconfig)
	if err != nil {
		glog.Fatalf("failed to create config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Fatalf("failed to create client: %v", err)
	}

	// The controller needs to know what the server version is because out-of-tree
	// provisioners aren't officially supported until 1.5
	serverVersion, err := clientset.Discovery().ServerVersion()
	if err != nil {
		glog.Fatalf("error getting server version: %v", err)
	}

	provisioner := controller.NewProvisionController(
		clientset,
		*provisionerName,
		slow.NewProvisioner(clientset, *nodePluginName, &delayOptions),
		serverVersion.GitVersion,
	)

	glog.Infof("Running slow provisioner with delays: %s", delayOptions.ToString())

	provisioner.Run(wait.NeverStop)
}

func setDefaultDelay(o *slow.DelayOptions, dr slow.DelayRange) {
	o.Provision = dr
	o.Delete = dr
}

func parseDelayArgs(o *slow.DelayOptions, args string) error {
	parts := strings.Split(args, ",")
	for _, p := range parts {
		eq := strings.Index(p, "=")
		if eq <= 0 {
			return fmt.Errorf("invalid format: '%s' is not in format of 'CALL=delay'", p)
		}

		op := p[:eq]
		val := p[eq+1:]

		var dr *slow.DelayRange

		switch op {
		case "Provision":
			dr = &o.Provision
		case "Delete":
			dr = &o.Delete
		default:
			return fmt.Errorf("no delay option for operation '%s'", op)
		}

		ret, err := parseDelay(val)
		if err != nil {
			return fmt.Errorf("failed to parse delay value for '%s': %v", p, err)
		}

		*dr = ret
	}

	return nil
}

func parseDelay(delay string) (slow.DelayRange, error) {
	if delay == "inf" {
		return slow.DelayRange{-1, -1}, nil
	}

	dots := strings.Index(delay, "..")
	if dots == -1 {
		if d, err := strconv.Atoi(delay); err != nil {
			return slow.DelayRange{}, err
		} else {
			return slow.DelayRange{d, d}, nil
		}
	} else {
		min, err := strconv.Atoi(delay[:dots])
		if err != nil {
			return slow.DelayRange{}, err
		}

		max, err := strconv.Atoi(delay[dots+2:])
		if err != nil {
			return slow.DelayRange{}, err
		}

		if min >= max {
			return slow.DelayRange{}, fmt.Errorf("%d ≮ %d", min, max)
		}

		return slow.DelayRange{Min: min, Max: max}, nil
	}
}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}
