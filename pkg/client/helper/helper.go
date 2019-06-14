package helper

import (
	"flag"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"strings"
)

func GetKubeConfigPath() string {
	// For multiple calls
	if fl := flag.Lookup("kubeconfig"); fl != nil {
		return fl.Value.String()
	}
	// Set kubeconfig value
	var kubeconfig *string
	var home string
	if home = os.Getenv("HOME"); home == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	return *kubeconfig
}

func GetContexts() []string {
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: GetKubeConfigPath()},
		&clientcmd.ConfigOverrides{}).RawConfig()

	if err != nil {
		panic(err.Error())
	}

	ctxs := make([]string, 0, len(config.Contexts))
	for k, _ := range config.Contexts { // Currently ignoring mappings
		// Hardcode ignore test clusters
		if !strings.Contains(k, "test") { // REVIEW: Remove this when possible
			ctxs = append(ctxs, k)
		}
	}

	return ctxs
}