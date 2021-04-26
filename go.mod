module github.com/bloodorangeio/octant-helm

go 1.13

require (
	github.com/Azure/go-autorest/autorest v0.11.18 // indirect
	github.com/elazarl/goproxy/ext v0.0.0-20210110162100-a92cc753f88e // indirect
	github.com/vmware-tanzu/octant v0.19.0
	helm.sh/helm/v3 v3.4.2
	k8s.io/api v0.19.4
	k8s.io/apimachinery v0.19.4
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
	k8s.io/client-go => k8s.io/client-go v0.19.3
)
