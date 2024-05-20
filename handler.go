package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func handler(ctx context.Context) error {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("[ERROR]: %s", err)
		return err
	}
	svc_ec2 := ec2.NewFromConfig(cfg)

	/**
	* Get all the regions in the account and then we
	* can loop over all the regions to scan more
	 */
	regions, err := svc_ec2.DescribeRegions(context.TODO(), &ec2.DescribeRegionsInput{})
	if err != nil {
		log.Printf("[ERROR]: %s", err)
		return err
	}

	// Loop
	var security_groups []types.SecurityGroup
	for _, region := range regions.Regions {
		cfg.Region = *region.RegionName

		/**
		* We are now looping over all the regions available
		* in this account and now we need to loop over all
		* all the security groups in the region
		 */
		paginator := ec2.NewDescribeSecurityGroupsPaginator(svc_ec2, &ec2.DescribeSecurityGroupsInput{})
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				log.Printf("[ERROR]: %s", err)
			}
			fmt.Println(page.SecurityGroups)
			security_groups = append(security_groups, page.SecurityGroups...)
		}
	}

	fmt.Println(security_groups)

	return nil
}

func main() {
	lambda.Start(handler)
}
