# Panther 🐆

[![Go Report Card](https://goreportcard.com/badge/valerauko/panther)](https://goreportcard.com/report/valerauko/panther)

A small service to use as a [HTTP provider](https://go-acme.github.io/lego/dns/httpreq/) for lego (or stuff that uses lego like Traefik) on Civo.

## Usage

1. release into your cluster by some means and set the `CIVO_API_TOKEN` and `CIVO_API_REGION` environment variables
2. set up your lego to use Panther's endpoint as a HTTP provider for the DNS-01 challenge
