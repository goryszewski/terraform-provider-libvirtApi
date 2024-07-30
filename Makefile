update:
	go get -u github.com/goryszewski/libvirtApi-client

plan:
	cd ./example/ ; terraform plan

apply:
	cd ./example/ ; TF_LOG=INFO  terraform apply   -auto-approve

run: 
	go run main.go

install:
	go install .