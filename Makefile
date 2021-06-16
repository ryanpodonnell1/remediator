DEFAULT_AWS_PROFILE:="proton"
VPC_ID:=vpc-0748e15cfac36c97f

.PHONY: build
tf_init: 
	export AWS_PROFILE=$(DEFAULT_AWS_PROFILE) && \
	cd terraform/ && \
	terraform init 

tf_apply: tf_init
	export AWS_PROFILE=$(DEFAULT_AWS_PROFILE) && \
	cd terraform/ && \
	terraform apply -var vpc_id=$(VPC_ID)

tf_destroy: tf_init
	export AWS_PROFILE=$(DEFAULT_AWS_PROFILE) && \
	cd terraform/ && \
	terraform destroy -var vpc_id=$(VPC_ID)

clean:
	rm bin/*
	
build: clean
	cd remediator && GOMODULE111=on go build -o ../bin/remediator

detect: build
	export AWS_PROFILE=$(DEFAULT_AWS_PROFILE) && \
	cd bin && \
	./remediator detect

remediate_dry: build
	export AWS_PROFILE=$(DEFAULT_AWS_PROFILE) && \
	cd bin && \
	./remediator remediate --dry-run=true

remediate_active: build
	export AWS_PROFILE=$(DEFAULT_AWS_PROFILE) && \
	cd bin && \
	./remediator remediate --dry-run=false


