
benchmarker:
	docker run -ti --rm --ulimit nofile=65535:65535 --network=host alpine/bombardier --method=POST  --fasthttp    -c 400 -d 600s -t 1s -n 1 -H "Content-Type: application/json" --body="[{\"AmountRequest\":200000,\"PartnerCode\":\"TEST\",\"PartnerIdentification\":1234,\"OrderID\":12132131,\"TypeRequest\":\"payment\"},{\"AmountRequest\":200000,\"PartnerCode\":\"TEST1\",\"PartnerIdentification\":1234,\"OrderID\":12132131,\"TypeRequest\":\"payment\"}]"  -l http://localhost:3000/request-balance

test:
	go test -v ./...

generate-wire-di:
	go install github.com/google/wire/cmd/wire@latest
	cd balance/
	wire

cache-authen-git:
	git config --global credential.helper cache
	// apter pull, push one times, after success
build-docker-images:
	docker build   -t balance-service:latest  .

build-docker-images-no-cache:
	docker build   -t balance-service:latest --no-cache=true .

run-docker-container:
	docker run -p 8080:8080  --env-file=devops/.env --name payment-balance-service   balance-service:latest


update-kube-config:
	aws eks update-kubeconfig --name payment-gateway-balance-calculator --region ap-southeast-1 --role-arn arn:aws:iam::387867911189:user/nghiaIT
total-kube-config:
	ls -l ~/.kube

show-kube-config:
	cat ~/.kube/config

set-default-kube-config:
	kubectl config use-context arn:aws:eks:ap-southeast-1:387867911189:cluster/payment-gateway-balance-calculator

apply-aws-auth-eks:
	kubectl apply -f eks/aws-auth.yaml

cluster-info:
	kubectl cluster-info

switch-aws-profile:
	 export AWS_PROFILE=...

