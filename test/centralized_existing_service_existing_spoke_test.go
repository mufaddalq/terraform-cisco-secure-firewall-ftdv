package gwlb_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

func TestExistingServiceExistingSpoke(t *testing.T) {
	fmt.Println("Starting tests...")
	t.Parallel()

	terraDir1 := "./existing_spoke/."
	terraformOptions1 := &terraform.Options{
		TerraformDir: terraDir1,
		VarFiles:     []string{"../../t5.tfvars"},
	}

	terraDir2 := "../examples/centralized_architecture_with_fmc"
	terraformOptions2 := &terraform.Options{
		TerraformDir: terraDir2,
		VarFiles:     []string{"../../t6.tfvars"},
	}

	test_structure.RunTestStage(t, "build_fmc", func() {
		fmt.Println("Setup")

		// Save options for later test stages
		test_structure.SaveTerraformOptions(t, terraDir1, terraformOptions1)

		// Triggers the terraform init and terraform apply commandåç
		terraform.InitAndApply(t, terraformOptions1)

	})

	// Specify the Terraform options
	test_structure.RunTestStage(t, "setup", func() {
		time.Sleep(25 * time.Minute)

		fmt.Println("Setup")

		FMC_EIP := terraform.OutputList(t, terraformOptions1, "fmcv_eip")
		FMC_IP := terraform.OutputList(t, terraformOptions1, "fmcv_ip")

		if terraformOptions2.Vars == nil {
			terraformOptions2.Vars = make(map[string]interface{})
		}

		// Set the FMC output as a variable for the second Terraform configuration
		terraformOptions2.Vars["fmc_ip"] = FMC_IP[0]
		terraformOptions2.Vars["fmc_host"] = FMC_EIP[0]

		// Save options for later test stages
		test_structure.SaveTerraformOptions(t, terraDir2, terraformOptions2)

		// Triggers the terraform init and terraform apply command
		terraform.InitAndApply(t, terraformOptions2)
	})

	// go test -v test/main_test.go -timeout 60m
	// Defer the destruction of resources until the test function is complete
	defer test_structure.RunTestStage(t, "teardown", func() {
		fmt.Println("Waiting for 3 minutes...")
		time.Sleep(3 * time.Minute) // To let the deployment finish before we start destroying.
		terraform.Destroy(t, terraformOptions2)
		defer test_structure.RunTestStage(t, "teardown_fmc", func() {
			fmt.Println("Destroying FMC....")
			terraform.Destroy(t, terraformOptions1)
		})
	})

}