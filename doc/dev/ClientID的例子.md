
# client_id可能的例子


```
=== RUN   TestClientId_GenerateGatewayClientId
    generateClientId_test.go:15: Gen :Ipv4 ClientID AAAAAAAAAAExOTIuMTY4LjAuMQ==
    generateClientId_test.go:17: Parse :&{ClientGatewayNum:1 ClientGatewayAddr:192.168.0.1}
    generateClientId_test.go:27: Gen :Ipv6 ClientID AAAAAAAAAAEyMDAxOjBkYjg6ODVhMzowMDAwOjAwMDA6OGEyZTowMzcwOjczMzQ=
    generateClientId_test.go:29: Parse :&{ClientGatewayNum:1 ClientGatewayAddr:2001:0db8:85a3:0000:0000:8a2e:0370:7334}
    generateClientId_test.go:39: Gen :Ipv6 ClientID AAAAAAAAAAEyMDAxOmRiODo4NWEzOjo4YTJlOjM3MDo3MzM0 [compressed
    generateClientId_test.go:41: Parse :&{ClientGatewayNum:1 ClientGatewayAddr:2001:db8:85a3::8a2e:370:7334}
    generateClientId_test.go:51: Gen :domain(www.example.com) ClientID AAAAAAAAAAEyMDAxOjBkYjg6ODVhMzowMDAwOjAwMDA6OGEyZTowMzcwOjczMzQ=
    generateClientId_test.go:53: Parse :&{ClientGatewayNum:1 ClientGatewayAddr:www.example.com}
--- PASS: TestClientId_GenerateGatewayClientId (0.00s)
PASS
```