package networkmanager_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/networkmanager"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/PixarV/terraform-provider-ritt/internal/acctest"
	"github.com/PixarV/terraform-provider-ritt/internal/conns"
	tfec2 "github.com/PixarV/terraform-provider-ritt/internal/service/ec2"
	tfnetworkmanager "github.com/PixarV/terraform-provider-ritt/internal/service/networkmanager"
	"github.com/PixarV/terraform-provider-ritt/internal/tfresource"
)

func TestAccNetworkManagerTransitGatewayRegistration_serial(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"basic":                     testAccNetworkManagerTransitGatewayRegistration_basic,
		"disappears":                testAccNetworkManagerTransitGatewayRegistration_disappears,
		"disappears_TransitGateway": testAccNetworkManagerTransitGatewayRegistration_disappears_TransitGateway,
		"crossRegion":               testAccNetworkManagerTransitGatewayRegistration_crossRegion,
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccNetworkManagerTransitGatewayRegistration_basic(t *testing.T) {
	resourceName := "aws_networkmanager_transit_gateway_registration.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, networkmanager.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTransitGatewayRegistrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayRegistrationConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayRegistrationExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetworkManagerTransitGatewayRegistration_disappears(t *testing.T) {
	resourceName := "aws_networkmanager_transit_gateway_registration.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, networkmanager.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTransitGatewayRegistrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayRegistrationConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayRegistrationExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tfnetworkmanager.ResourceTransitGatewayRegistration(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccNetworkManagerTransitGatewayRegistration_disappears_TransitGateway(t *testing.T) {
	resourceName := "aws_networkmanager_transit_gateway_registration.test"
	transitGatewayResourceName := "aws_ec2_transit_gateway.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, networkmanager.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTransitGatewayRegistrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayRegistrationConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayRegistrationExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tfec2.ResourceTransitGateway(), transitGatewayResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccNetworkManagerTransitGatewayRegistration_crossRegion(t *testing.T) {
	resourceName := "aws_networkmanager_transit_gateway_registration.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	var providers []*schema.Provider

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t); acctest.PreCheckMultipleRegion(t, 2) },
		ErrorCheck:        acctest.ErrorCheck(t, networkmanager.EndpointsID),
		ProviderFactories: acctest.FactoriesAlternate(&providers),
		CheckDestroy:      testAccCheckTransitGatewayRegistrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayRegistrationCrossRegionConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayRegistrationExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTransitGatewayRegistrationDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).NetworkManagerConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_networkmanager_transit_gateway_registration" {
			continue
		}

		globalNetworkID, transitGatewayARN, err := tfnetworkmanager.TransitGatewayRegistrationParseResourceID(rs.Primary.ID)

		if err != nil {
			return err
		}

		_, err = tfnetworkmanager.FindTransitGatewayRegistrationByTwoPartKey(context.TODO(), conn, globalNetworkID, transitGatewayARN)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("Network Manager Transit Gateway Registration %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckTransitGatewayRegistrationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Network Manager Transit Gateway Registration ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).NetworkManagerConn

		globalNetworkID, transitGatewayARN, err := tfnetworkmanager.TransitGatewayRegistrationParseResourceID(rs.Primary.ID)

		if err != nil {
			return err
		}

		_, err = tfnetworkmanager.FindTransitGatewayRegistrationByTwoPartKey(context.TODO(), conn, globalNetworkID, transitGatewayARN)

		if err != nil {
			return err
		}

		return nil
	}
}
func testAccTransitGatewayRegistrationConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_networkmanager_global_network" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_networkmanager_transit_gateway_registration" "test" {
  global_network_id   = aws_networkmanager_global_network.test.id
  transit_gateway_arn = aws_ec2_transit_gateway.test.arn
}
`, rName)
}

func testAccTransitGatewayRegistrationCrossRegionConfig(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAlternateRegionProvider(), fmt.Sprintf(`
resource "aws_networkmanager_global_network" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway" "test" {
  provider = "awsalternate"

  tags = {
    Name = %[1]q
  }
}

resource "aws_networkmanager_transit_gateway_registration" "test" {
  global_network_id   = aws_networkmanager_global_network.test.id
  transit_gateway_arn = aws_ec2_transit_gateway.test.arn
}
`, rName))
}
