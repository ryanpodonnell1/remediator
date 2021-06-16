DEFAULT_AWS_PROFILE:="proton" #FOR TF_APPLY AND REMEDIATOR AUTH
VPC_ID:="vpc-0748e15cfac36c97f" #FOR TF_APPLY

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

#### REMEDIATOR
clean:
	rm bin/* || echo "no file"

build: clean
	cd remediator && GOMODULE111=on go build -o ../bin/remediator

detect: build #PASSIVE DETECTON
	export AWS_PROFILE=$(DEFAULT_AWS_PROFILE) && \
	cd bin && \
	./remediator detect

remediate_dry: build #DRY RUN ACTIVE REMEDIATON
	export AWS_PROFILE=$(DEFAULT_AWS_PROFILE) && \
	cd bin && \
	./remediator remediate --dry-run=true
 
remediate_active: build #ACTIVE REMEDIATION
	export AWS_PROFILE=$(DEFAULT_AWS_PROFILE) && \
	cd bin && \
	./remediator remediate --dry-run=false


