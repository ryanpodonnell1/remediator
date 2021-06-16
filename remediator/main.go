package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli/v2"
)

var sess *session.Session
var ec2Creds *ec2.EC2
var bad_cidrs = []string{"0.0.0.0/0"}
var bad_ports = []int64{22, 3389, 3306}

func main() {
	sess, err := session.NewSession()
	checkErrLogFatal(err)
	ec2Creds = ec2.New(sess)

	app := &cli.App{
		Name:  "remediator",
		Usage: "run, detect, evaluate, remediate",
	}
	app.Commands = []*cli.Command{
		{
			Name:  "detect",
			Usage: "parent level command to passively detect/report malformed configurations",
			Action: func(c *cli.Context) error {
				groups := DetectMalformedSecurityGroups(ec2Creds)
				fmt.Println("[WARNING] GROUPS IN VIOLATION:")
				SummarizeSGOutput(groups)
				return nil
			},
		},
		{
			Name:  "remediate",
			Usage: "parent level command to ACTIVELY remediate malformed configurations",
			Action: func(c *cli.Context) error {
				groups := DetectMalformedSecurityGroups(ec2Creds)
				fmt.Println("[WARNING] GROUPS IN VIOLATION:")
				SummarizeSGOutput(groups)
				fmt.Println("[WARNING] GROUPS BEING REMEDIATED")
				RemediateMalformedSecurityGroups(ec2Creds, c, groups)
				return nil
			},
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "dry-run",
					Value: true,
					Usage: "--dry-run=<bool>, set true to run a dry run of autoremediation",
				},
			},
		},
	}

	checkErrLogFatal(app.Run(os.Args))
}

//DetectMalformedSecurityGroups returns an output of non compliant securitygroups
func DetectMalformedSecurityGroups(client *ec2.EC2) []ec2.DescribeSecurityGroupsOutput {

	allGroups := []ec2.DescribeSecurityGroupsOutput{}
	input := ec2.DescribeSecurityGroupsInput{}
	err := client.DescribeSecurityGroupsPages(&input,
		func(page *ec2.DescribeSecurityGroupsOutput, lastPage bool) bool {
			allGroups = append(allGroups, *page)
			return true
		})

	checkErrLogFatal(err)

	return detectMalformedSecurityGroupsIngressEvaluate(allGroups)
}

//takes in a slice of securitygroups and returns ones that are non compliant
func detectMalformedSecurityGroupsIngressEvaluate(sgs []ec2.DescribeSecurityGroupsOutput) []ec2.DescribeSecurityGroupsOutput {
	badGroups := []ec2.DescribeSecurityGroupsOutput{}

	for _, sg := range sgs {
		for _, entries := range sg.SecurityGroups {
			for _, permission := range entries.IpPermissions {
				var from, to int64
				if *permission.IpProtocol == "-1" {
					from = 0
					to = 65535
				}
				if permission.FromPort != nil && permission.ToPort != nil {
					from = *permission.FromPort
					to = *permission.ToPort
				}
				for _, ip := range permission.IpRanges {
					cidr := ip
					cidrBool := isValueInSlice(*cidr.CidrIp, bad_cidrs)
					PortBool := areNumbersInRange(bad_ports, from, to)
					if cidrBool && PortBool {
						output := ec2.DescribeSecurityGroupsOutput{
							SecurityGroups: []*ec2.SecurityGroup{
								{
									GroupId:     entries.GroupId,
									GroupName:   entries.GroupName,
									VpcId:       entries.VpcId,
									Description: entries.Description,
									IpPermissions: []*ec2.IpPermission{
										{
											IpProtocol: permission.IpProtocol,
											FromPort:   permission.FromPort,
											ToPort:     permission.ToPort,
											IpRanges:   []*ec2.IpRange{cidr},
										},
									},
								},
							},
						}
						badGroups = append(badGroups, output)
					}

				}

			}
		}
	}
	return badGroups
}

//isValueInSlice determins if a given string is in a slice of string
func isValueInSlice(value string, xs []string) bool {
	for _, v := range xs {
		if v == value {
			return true
		}
	}
	return false
}

func areNumbersInRange(value []int64, low, high int64) bool {
	for _, v := range value {
		if v >= low && v <= high {
			return true
		}
	}
	return false
}

func SummarizeSGOutput(sgs []ec2.DescribeSecurityGroupsOutput) {
	for _, sg := range sgs {
		for _, entries := range sg.SecurityGroups {
			fmt.Printf("=> sg_id: %v, sg_name: %v, permissions: %#v\n\n", *entries.GroupId, *entries.GroupName, entries.IpPermissions)
		}
	}
}

func checkErrLogFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//RemediateMalformedSecurityGroups ACTIVELY revokes ingress rules for the given DescribeSecurityGroupsOutput provided
func RemediateMalformedSecurityGroups(client *ec2.EC2, c *cli.Context, sgs []ec2.DescribeSecurityGroupsOutput) {
	for _, sg := range sgs {
		for _, entries := range sg.SecurityGroups {
			for _, permissions := range entries.IpPermissions {
				for _, ip := range permissions.IpRanges {
					input := ec2.RevokeSecurityGroupIngressInput{
						DryRun:     aws.Bool(c.Bool("dry-run")),
						CidrIp:     ip.CidrIp,
						FromPort:   permissions.FromPort,
						ToPort:     permissions.ToPort,
						GroupId:    entries.GroupId,
						IpProtocol: permissions.IpProtocol,
					}
					output, err := client.RevokeSecurityGroupIngress(&input)
					if err != nil {
						log.Println("[ERROR] there was an error observed when removing securitygroup rules: ", err)
					} else {
						fmt.Printf("=> SUCCESSFULLY REVOKED INGRESS RULE ON %v:\n AWS RESPONSE: %q\n", *entries.GroupId, output)
					}
				}
			}
		}
	}

}
