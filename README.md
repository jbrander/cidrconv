# cidrconv

Weekend project written in late 2022 as Golang learning exericse.

CLI utility that generates the arugments expected by Terraform's [cidrsubnet](https://developer.hashicorp.com/terraform/language/functions/cidrsubnet) function to return the desired subnet CIDR given a specified network CIDR block.

```
$ cidrconv -s 192.168.2.0/24 -n 192.168.0.0/16
cidrsubnet("192.168.0.0/16", 8, 2)
```

Requires Go version 1.19 or higher. Build using:

```
git clone https://github.com/jbrander/cidrconv
cd cidrconv
go build
```
