# 环境
java17及以上(需配置环境变量)
golang1.22
# 使用方式
将soot-analysis-1.0-SNAPSHOT.jar文件和go文件放在同一目录下
运行命令：go run main.go (目标文件夹)
比如：
    go run main.go C:\Users\li_sh\Desktop\msexample\example\0ruoyiExample
# 样例要求
要求每个文件夹下放一个jar包及相应的一个k8s部署文件
比如：
    - example
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
    go run main.go C:\Users\li_sh\Desktop\msexample\example\0ruoyiExample -c C:\Users\li_sh\Desktop\msexample\example\0ruoyiExample\nacos_config_export\DEFAULT_GROUP
# 输出
在output文件夹下生成networkpolicy的json文件（供450的策略生成器使用，可以忽视）
和可以直接使用的网络策略yaml文件
