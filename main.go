package main

import (
	"github.com/alecthomas/kingpin"
	pr2porm "github.com/wy100101/pr2porm/pkg"
)

var (
	rulesFile    = kingpin.Flag("file.rules", "Prometheus rules file to convert.").Short('f').Required().ExistingFile()
	manifestFile = kingpin.Flag("file.output", "Output file for the dashboard configmap.").Short('o').Default("").String()
	rulesName    = kingpin.Flag("rules.name", "Rules file manifest name. (Default: rules file basename)").Short('n').Default("").String()
	k8sNamespace = kingpin.Flag("k8s.namespace", "kubernetes namespace for the configmap.").Short('N').Default("monitoring").String()
	annotations  = kingpin.Flag("metadata.annotations", "Annotations to add to the rules file.").Short('a').StringMap()
	labels       = kingpin.Flag("metadata.labels", "labels to add to the rules file.").Short('l').StringMap()
)

func main() {
	kingpin.Parse()
	errs := pr2porm.ProcessRulesFile(*rulesFile, *manifestFile, *k8sNamespace, *rulesName, labels, annotations)
	if errs != nil {
		panic(errs)
	}
}
