module github.com/bloodorangeio/octant-helm

go 1.13

require (
	github.com/elazarl/goproxy/ext v0.0.0-20210110162100-a92cc753f88e // indirect
	github.com/vmware-tanzu/octant v0.18.0
	helm.sh/helm/v3 v3.0.0
	k8s.io/api v0.19.3
	k8s.io/apimachinery v0.19.3
)

replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309
