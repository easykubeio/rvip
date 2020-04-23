# rvip（Raft Virtual IP）
rvip can be used to share a Virtual/Floating IP address between many computers. The IP address is assigned to only one computer. If that computer goes down, the same IP address is then reassigned to another computer in the cluster.

## Features

- Self contained binary that can be downloaded from GitHub
- No external dependencies as in sevices
- Uses [Raft Consensus](https://raft.github.io/) for cluster communication

## Usage

Assuming that you have the following setup:

- three Linux servers with the IP addresses: 192.168.0.101, 192.168.0.102, 192.168.0.103)
- the three servers have the network interface eth1
- the Virtual-IP is 192.168.0.50
- the port 10000 for Raft on all three servers is free

### Binary

Download the binary from [release page](https://github.com/easykubeio/rvip/releases).

Then on each server run the following commands

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


