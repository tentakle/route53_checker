package main

import (
	"crypto/tls"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"log"
)

func main() {
	sess := session.Must(session.NewSession())

	svc := route53.New(sess)

	params := &route53.ListHostedZonesInput{}

	pageNum := 0
	err := svc.ListHostedZonesPages(params,
		func(page *route53.ListHostedZonesOutput, lastPage bool) bool {
			pageNum++

			for _, zone := range page.HostedZones {
				fmt.Printf("Hosted Zone: %v\n\n", *zone.Name)

				params := &route53.ListResourceRecordSetsInput{
					HostedZoneId: zone.Id,
				}

				pageNum := 0
				err := svc.ListResourceRecordSetsPages(params,
					func(page *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
						pageNum++

						for _, record := range page.ResourceRecordSets {
							if *record.Type == "A" || *record.Type == "AAAA" || *record.Type == "CNAME" {
								domain := *record.Name
								conn, err := tls.Dial("tcp", domain+":443", nil)
								if err != nil {
									fmt.Printf("Type: %v\nName: %v\nError: %v\n\n", *record.Type, *record.Name, err)
								} else {
									err = conn.VerifyHostname(domain)
									if err != nil {
										fmt.Printf("Type: %v\nName: %v\nError: %v\n\n", *record.Type, *record.Name, err)
										conn.Close()
									}
								}
							}
						}
						return *page.IsTruncated
					})
				if err != nil {
					log.Println(err)
				}
			}
			return *page.IsTruncated
		})
	if err != nil {
		log.Println(err)
	}
}
