# Panther 🐆

![Panther](https://repository-images.githubusercontent.com/364797946/05d71b00-aebc-11eb-82e8-e506c20c8390)

[![Go Report Card](https://goreportcard.com/badge/valerauko/panther)](https://goreportcard.com/report/valerauko/panther)

A small service to use as a [HTTP provider](https://go-acme.github.io/lego/dns/httpreq/) for lego (or stuff that uses lego like Traefik) on Civo.

## Usage

1. you can use the [manifests](https://github.com/valerauko/panther/tree/main/manifests) to release Panther into your cluster
1. set up your lego to use Panther's endpoint as a HTTP provider for the DNS-01 challenge

The provided manifests follow the `main` branch.

### Create secret

If you use the provided manifests, you'll need a Secret called `civo-api-secret` in Panther's namespace as well. If you install it to `kube-system`, that could be created like this:

```
$ read TOKEN
# <input your token from https://www.civo.com/account/security>

$ read REGION
# <input your region>

$ cat <<EOF | kubectl -n kube-system apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: civo-api-secret
data:
  TOKEN: $(echo -n $TOKEN | base64)
  REGION: $(echo -n $REGION | base64)
EOF
# secret/civo-api-secret created
```

## Credit

The cover picture of the repo is by <a href="https://unsplash.com/@soberanes?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText">Uriel Soberanes</a> on <a href="https://unsplash.com/s/photos/panther?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText">Unsplash</a>.
