package pr2porm

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/prometheus/prometheus/pkg/rulefmt"
	yaml "gopkg.in/yaml.v3"
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

// ProcessRulesFile(ruleFile, manifestFile, namespace, name, annotations, labels) []error
// Given a prometheus rules file will generate a PrometheusOperator promrule maniest into the manifestFile
func ProcessRulesFile(prf, pormf, ns, n string, ls, as *map[string]string) error {
	var bprfns string
	rgs, errs := rulefmt.ParseFile(prf)
	if len(errs) > 0 {
		return fmt.Errorf("%v", errs)
	}
	bprf := filepath.Base(prf)
	if strings.HasSuffix(bprf, ".yaml") {
		bprfns = strings.TrimSuffix(bprf, ".yaml")
	} else if strings.HasSuffix(bprf, ".yml") {
		bprfns = strings.TrimSuffix(bprf, ".yml")
	}
	if pormf == "" {
		pormf = fmt.Sprintf("%s.promrule.yaml", bprfns)
	}
	if n == "" {
		n = strings.Replace(bprfns, "_", "-", -1)
	}

	mf := prometheusRulesManifest{
		APIVersion: "monitoring.coreos.com/v1",
		Kind:       "PrometheusRule",
		Metadata: manifestMetadata{
			Name:        n,
			Namespace:   ns,
			Labels:      *ls,
			Annotations: *as,
		},
		Spec: rgs,
	}

	op, err := yaml.Marshal(&mf)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(pormf, op, 0666)

	if err != nil {
		return err
	}
	return nil
}
