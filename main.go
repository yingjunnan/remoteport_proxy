package main

import (
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"net"
)

func handleConnection(localConn net.Conn, remoteAddr string) {
	remoteConn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Printf("无法连接到远程地址: %s", err)
		return
	}
	defer remoteConn.Close()

	go func() {
		_, err := io.Copy(remoteConn, localConn)
		if err != nil {
			log.Printf("向远程复制数据时发生错误: %s", err)
		}
	}()

	_, err = io.Copy(localConn, remoteConn)
	if err != nil {
		log.Printf("向本地复制数据时发生错误: %s", err)
	}
}

func main() {
	viper.SetConfigName("config") // 配置文件名为config.yaml
	viper.SetConfigType("yaml")   // 配置文件类型为yaml
	viper.AddConfigPath(".")      // 配置文件路径为当前目录

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("无法读取配置文件: %s", err)
	}

	proxyConfigs := viper.GetStringMap("proxies")
	for localPort, remoteAddr := range proxyConfigs {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%s", localPort))
		if err != nil {
			log.Printf("无法在端口 %s 上启动监听器: %s", localPort, err)
			continue
		}
		defer listener.Close()

		log.Printf("代理已启动在端口 %s，转发到地址 %s", localPort, remoteAddr)

		for {
			localConn, err := listener.Accept()
			if err != nil {
				log.Printf("无法接受连接: %s", err)
				continue
			}

			go handleConnection(localConn, remoteAddr.(string))
		}
	}
}
