# 环境
java17及以上(需配置环境变量)
golang1.21
# 使用方式
将soot-analysis-1.0-SNAPSHOT.jar文件和go文件放在同一目录下
运行命令：go run main.go (目标文件夹)
比如：go run main.go C:\Users\li_sh\Desktop\msexample\example\0ruoyiExample
# 样例要求
要求每个文件夹下放一个jar包及相应的一个k8s部署文件
比如
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
# 输出
在output文件夹下生成network=policy的json文件
