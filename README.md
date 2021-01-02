# gogrpccertAndgateway
1.根目录下的main.go实现了grpc 的双向认证
2.server 和client 实现了grpc gateway，如果grpc启用https 那么gateway 暴露出来的就有些问题【gateway 一般是前段用的 不建议用https】
