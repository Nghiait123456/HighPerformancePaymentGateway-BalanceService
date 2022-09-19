benchmarker:
	docker run -ti --rm --ulimit nofile=65535:65535 --network=host alpine/bombardier --method=POST  --fasthttp    -c 400 -d 600s -t 1s -n 1 -H "Content-Type: application/json" --body="[{\"AmountRequest\":200000,\"PartnerCode\":\"TEST\",\"PartnerIdentification\":1234,\"OrderID\":12132131,\"TypeRequest\":\"payment\"},{\"AmountRequest\":200000,\"PartnerCode\":\"TEST1\",\"PartnerIdentification\":1234,\"OrderID\":12132131,\"TypeRequest\":\"payment\"}]"  -l http://localhost:3000/request-balance

test:
	go test -v ./...

