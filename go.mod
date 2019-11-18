module github.com/bloodorangeio/octant-helm

go 1.13

require (
	github.com/vmware-tanzu/octant v0.9.2-0.20191116231443-28aa3e91ffa5
	helm.sh/helm/v3 v3.0.0
	k8s.io/api v0.0.0-20191016225839-816a9b7df678
	k8s.io/apimachinery v0.0.0-20191016225534-b1267f8c42b4
)

replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309
