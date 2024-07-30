# 环境
java17及以上(需配置环境变量)
golang1.22
CGO_ENABLED=1
gcc (Ubuntu 13.2.0-23ubuntu4) 13.2.0
# 使用方式
将soot-analysis-1.0-SNAPSHOT.jar文件和go文件放在同一目录下
运行命令：go run main.go policy (目标文件夹)
比如：
    go run main.go policy C:\Users\li_sh\Desktop\msexample\example\0ruoyiExample
# 样例要求
    - pods
        - ruoyi-auth
            - ruoyi-auth.jar
            - deployment.yaml
        - ruoyi-system
            - ruo-system-api
                - ruo-system-api.jar
                - deployment.yaml
            - ruo-system-file
                - ruo-system-file.jar
                - ruoyi-deployment.yaml              
如果有注册中心的额外的配置配置文件 可以放在另一个文件夹中，并使用参数-c,--configS
    go run main.go policy C:\Users\li_sh\Desktop\msexample\example\0ruoyiExample -c C:\Users\li_sh\Desktop\msexample\example\0ruoyiExample\nacos_config_export\DEFAULT_GROUP
若有中间件须在根目录下svc.json文件存放serviceName到labels的键值对
# 输出
可以直接使用的网络策略yaml文件
