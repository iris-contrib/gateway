module github.com/iris-contrib/gateway

go 1.15

require (
	github.com/apex/gateway/v2 v2.0.0-20200703123654-59bba3473042
	github.com/aws/aws-lambda-go v1.19.1
	github.com/kataras/iris/v12 v12.1.9-0.20200823145529-ef5685bf7eeb
)

replace github.com/apex/gateway/v2 v2.0.0-20200703123654-59bba3473042 => github.com/kataras/gateway/v2 v2.0.0-20200823133619-5f644b75fcd5
