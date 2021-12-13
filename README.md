# bssh

##编译
* 平台下执行

        make

##执行

* 配置文件为本地模式

        ./bin/bssh -mode local -conf conf/dev.json

* 配置文件为etcd模式 

        ./bin/bssh -mode etcd -h 127.0.0.1:2379 -path /web/

##配置
* 把本地文件设置到etcd内

        ./bin/etcdTool -set -path /web/ -local conf/dev.json -h 127.0.0.1:2379

* 从etcd 获取配置文件内容

        ./bin/etcdTool -get -path /web/ -h 127.0.0.1:2379

##配置文件格式

```
{
  "password": "1",              //认证的密码
  "user": "machao",             //认证的用户
  "command": [                  //可执行列表，为空则不限制
    "ls"
  ]
}
```