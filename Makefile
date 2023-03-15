
benchmark:
	docker run -ti --rm --ulimit nofile=65535:65535 --network=host alpine/bombardier --method=POST  --fasthttp    -c 400 -n 1000000   -H "Content-Type: application/json" --body="[{\"AmountRequest\":200000,\"PartnerCode\":\"TEST\",\"PartnerIdentification\":1234,\"OrderID\":12132131,\"TypeRequest\":\"payment\"},{\"AmountRequest\":200000,\"PartnerCode\":\"TEST\",\"PartnerIdentification\":1235,\"OrderID\":12132132,\"TypeRequest\":\"payment\"},{\"AmountRequest\":200000,\"PartnerCode\":\"TEST\",\"PartnerIdentification\":1236,\"OrderID\":12132134,\"TypeRequest\":\"payment\"},{\"AmountRequest\":200000,\"PartnerCode\":\"TEST\",\"PartnerIdentification\":1240,\"OrderID\":12132140,\"TypeRequest\":\"payment\"},{\"AmountRequest\":200000,\"PartnerCode\":\"TEST\",\"PartnerIdentification\":1241,\"OrderID\":12132141,\"TypeRequest\":\"payment\"},{\"AmountRequest\":200000,\"PartnerCode\":\"TEST1\",\"PartnerIdentification\":1242,\"OrderID\":12132142,\"TypeRequest\":\"payment\"},{\"AmountRequest\":200000,\"PartnerCode\":\"TEST1\",\"PartnerIdentification\":1243,\"OrderID\":12132143,\"TypeRequest\":\"payment\"},{\"AmountRequest\":200000,\"PartnerCode\":\"TEST1\",\"PartnerIdentification\":1244,\"OrderID\":12132144,\"TypeRequest\":\"payment\"},{\"AmountRequest\":200000,\"PartnerCode\":\"TEST1\",\"PartnerIdentification\":1245,\"OrderID\":12132145,\"TypeRequest\":\"payment\"},{\"AmountRequest\":200000,\"PartnerCode\":\"TEST1\",\"PartnerIdentification\":1246,\"OrderID\":12132146,\"TypeRequest\":\"payment\"}]"  -l http://localhost:8080/balance/request-balance

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
	docker build   -t balance-service:latest -f=devops/Dockerfile .

build-docker-images-no-cache:
	docker build   -t balance-service:latest --no-cache=true -f=devops/Dockerfile .

run-docker-container:
	docker run -p 8080:8080  --env-file=devops/.env balance-service:latest

run-main-go:
	go run main.go  --env-file=devops/.env

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

