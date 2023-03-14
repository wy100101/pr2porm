package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kingpin"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	yaml "gopkg.in/yaml.v3"
)

var (
	rulesFile    = kingpin.Flag("file.rules", "Prometheus rules file to convert.").Short('f').Required().ExistingFile()
	manifestFile = kingpin.Flag("file.output", "Output file for the dashboard configmap.").Short('o').Default("").String()
	compact      = kingpin.Flag("file.compact", "Output file with compact JSON embedded in ConfigMap.").Short('c').Default("false").Bool()
	rulesName    = kingpin.Flag("rules.name", "Rules file manifest name. (Default: rules file basename)").Short('n').Default("").String()
	k8sNamespace = kingpin.Flag("k8s.namespace", "kubernetes namespace for the configmap.").Short('N').Default("monitoring").String()
	annotations  = kingpin.Flag("metadata.annotations", "Annotations to add to the rules file.").Short('a').StringMap()
	labels       = kingpin.Flag("metadata.labels", "labels to add to the rules file.").Short('l').StringMap()
)

type manifestMetadata struct {
	Name        string            `yaml:"name"`
	Namespace   string            `yaml:"namespace"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type prometheusRulesManifest struct {
	APIVersion string              `yaml:"apiVersion"`
	Kind       string              `yaml:"kind"`
	Metadata   manifestMetadata    `yaml:"metadata"`
	Spec       *rulefmt.RuleGroups `yaml:"spec"`
}

func validateRuleGroups(rgs *rulefmt.RuleGroups) bool {
	for _, rg := range rgs.Groups {
		for _, r := range rg.Rules {
			if _, ok := r.Labels["team"]; !ok {
				return ok
			}
		}
	}
	return true
}

func main() {
	kingpin.Parse()
	var bdfns string
	rgs, err := rulefmt.ParseFile(*rulesFile)
	if err != nil {
		panic(err)
	}
	bdf := filepath.Base(*rulesFile)
	if strings.HasSuffix(bdf, ".yaml") {
		bdfns = strings.TrimSuffix(bdf, ".yaml")
	} else if strings.HasSuffix(bdf, ".yml") {
		bdfns = strings.TrimSuffix(bdf, ".yml")
	}
	if *manifestFile == "" {
		*manifestFile = fmt.Sprintf("%s.yaml", bdfns)
	}
	if *rulesName == "" {
		*rulesName = strings.Replace(bdfns, "_", "-", -1)
	}
	fmt.Printf("%v", rgs)

	mf := prometheusRulesManifest{
		APIVersion: "monitoring.coreos.com/v1",
		Kind:       "PrometheusRule",
		Metadata: manifestMetadata{
			Name:        *rulesName,
			Namespace:   *k8sNamespace,
			Labels:      *labels,
			Annotations: *annotations,
		},
		Spec: rgs,
	}

	op, errs := yaml.Marshal(&mf)
	if errs != nil {
		panic(errs)
	}

	errs = ioutil.WriteFile(*manifestFile, op, 0666)
	if errs != nil {
		panic(fmt.Sprintf("Error: %s could not be written (%s)", *manifestFile, err))
	}
}
