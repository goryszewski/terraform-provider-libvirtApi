update:
	go get -u github.com/goryszewski/libvirtApi-client

plan:
	cd ./example/ ; terraform plan

apply:
	cd ./example/ ; TF_LOG=INFO  terraform apply   -auto-approve

state:
	cd ./example/ ; TF_LOG=INFO  terraform show

run: 
	go run main.go

install:
	go install .