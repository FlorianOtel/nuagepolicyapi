package main

import (
	"flag"
	"fmt"
	"github.com/FlorianOtel/nuagepolicyapi/implementer"
	"github.com/FlorianOtel/nuagepolicyapi/policies"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	flagVSDCredentials = "vsd-credentials"
	flagAddPolicy      = "add-policy"
	flagDeletePolicy   = "delete-policy"
	flagPolicyFile     = "policy-file"
	flagPolicyID       = "policy-id"
	flagEnterprise     = "enterprise"
	flagDomain         = "domain"
)

func main() {
	var vsdCredentialsYAML string
	flag.StringVar(&vsdCredentialsYAML, flagVSDCredentials, "", "YAML file with VSD credentials")

	var addPolicy bool
	flag.BoolVar(&addPolicy, flagAddPolicy, false, "Add Nuage policy")

	var policyFile string
	flag.StringVar(&policyFile, flagPolicyFile, "", "Policy YAML")

	var delPolicy bool
	flag.BoolVar(&delPolicy, flagDeletePolicy, false, "Delete Nuage policy")

	var policyID string
	flag.StringVar(&policyID, flagPolicyID, "", "Policy ID")

	var enterprise string
	flag.StringVar(&enterprise, flagEnterprise, "", "Enterprise")

	var domain string
	flag.StringVar(&domain, flagDomain, "", "Domain")

	flag.Parse()

	if !addPolicy && !delPolicy {
		fmt.Printf("Invalid policy action\n")
		os.Exit(1)
	}

	if delPolicy {
		if policyID == "" {
			fmt.Printf("Missing policy ID\n")
			os.Exit(1)
		}

		if enterprise == "" {
			fmt.Printf("Enterprise missing\n")
			os.Exit(1)
		}

		if domain == "" {
			fmt.Printf("Domain missing\n")
			os.Exit(1)
		}
	}

	credFile, err := filepath.Abs(vsdCredentialsYAML)
	if err != nil {
		fmt.Printf("Unable to the absolute path for the vsd credential file\n")
		os.Exit(1)
	}

	credData, err := ioutil.ReadFile(credFile)
	if err != nil {
		fmt.Printf("Problem reading the VSD credentials\n")
		os.Exit(1)
	}

	var vsdCredentials implementer.VSDCredentials
	err = yaml.Unmarshal(credData, &vsdCredentials)
	if err != nil {
		fmt.Printf("Problem unmarshalling the VSD credentials\n")
		os.Exit(1)
	}

	var policyImplementer implementer.PolicyImplementer
	if err := policyImplementer.Init(&vsdCredentials); err != nil {
		fmt.Printf("Unable to connect to VSD\n")
		os.Exit(1)
	}

	if addPolicy {
		pfile, err := filepath.Abs(policyFile)
		if err != nil {
			fmt.Printf("Unable to the absolute path for the policy file\n")
			os.Exit(1)
		}

		policyData, err := ioutil.ReadFile(pfile)
		if err != nil {
			fmt.Printf("Problem reading the policy file\n")
			os.Exit(1)
		}

		nuagePolicy, err := policies.LoadPolicyFromYAML(string(policyData))
		if err != nil {
			fmt.Printf("Problem loading the nuage policy %+v\n", err)
			os.Exit(1)
		}

		err = policyImplementer.ImplementPolicy(nuagePolicy)
		if err != nil {
			fmt.Printf("Problem implementing the nuage policy %+v\n", err)
			os.Exit(1)
		}
	}

	if delPolicy {
		err := policyImplementer.DeletePolicy(policyID, enterprise, domain)
		if err != nil {
			fmt.Printf("Problem deleting the policy\n")
			os.Exit(1)
		}
	}
}
