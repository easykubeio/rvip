# rvip（Raft Virtual IP）
rvip可用于多台服务器共享一个虚拟IP。IP地址分配给其中一台服务器，当此服务器宕机IP地址将重新分配给集群中的其他服务器。

## 特点
- 不依赖其他服务
- 采用Raft算法

## 用法

环境:

- 三台linux服务器，IP地址: 192.168.0.101, 192.168.0.102, 192.168.0.103)
- 服务器网络接口: eth1
- 虚拟IP（vip）：192.168.0.50
- 可用端口：10000

Server1 (192.168.0.101):

```shell
sudo rvip -id server1 -bind 192.168.0.101:10000 -peers server1=192.168.0.101:10000,server2=192.168.0.102:10000,server3=192.168.0.103:10000 -interface eth1 -vip 192.168.0.50
```

Server2 (192.168.0.102):

```shell
sudo rvip -id server2 -bind 192.168.0.102:10000 -peers server1=192.168.0.101:10000,server2=192.168.0.102:10000,server3=192.168.0.103:10000 -interface eth1 -vip 192.168.0.50
```

Server3 (192.168.0.103):

```shell
sudo rvip -id server3 -bind 192.168.0.102:10000 -peers server1=192.168.0.101:10000,server2=192.168.0.102:10000,server3=192.168.0.103:10000 -interface eth1 -vip 192.168.0.50
```


