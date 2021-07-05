module github.com/featheredtoast/ddns-route53/cli

go 1.14

require (
	github.com/alecthomas/kong v0.2.17
	github.com/featheredtoast/ddns-route53/internal/aws v0.0.0-00010101000000-000000000000
	github.com/featheredtoast/ddns-route53/internal/iplookup v0.0.0-00010101000000-000000000000
)

replace github.com/featheredtoast/ddns-route53/internal/aws => ../internal/aws

replace github.com/featheredtoast/ddns-route53/internal/iplookup => ../internal/iplookup
