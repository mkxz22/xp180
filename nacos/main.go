package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// 定义 Nacos 配置结构体
type NacosConfig struct {
	Scheme      string
	Host        string
	Port        uint64
	NamespaceId string
	DataId      string
	Group       string
}

// 初始化 Nacos 配置
func GetNacosConfig() NacosConfig {
	return NacosConfig{
		Scheme:      "14.103.149.202",
		Host:        os.Getenv("NacosHost"),
		Port:        8848, // 假设默认端口
		NamespaceId: os.Getenv("NacosNamespaceId"),
		DataId:      os.Getenv("NacosDataId"),
		Group:       os.Getenv("NacosGroup"),
	}
}

func main() {
	NConfig := GetNacosConfig()
	fmt.Printf("获取环境信息,%s,%T,%v\n", NConfig.Scheme, NConfig.Scheme, NConfig)
	// 1. 配置Nacos服务器参数
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      "14.103.149.202", // Nacos服务器IP
			Port:        8848,             // Nacos端口
			Scheme:      "http",
			ContextPath: "/nacos", // NaCos 服务的上下文路径
		},
	}

	// 2. 客户端配置（含鉴权） 如果我们配置鉴权的话
	clientConfig := constant.ClientConfig{
		NamespaceId:         "bcdc858d-5542-453c-bb9b-88829c378dc8", // 命名空间ID（默认public）
		TimeoutMs:           5000,                                   // 请求超时时间
		NotLoadCacheAtStart: true,                                   // 启动时不读取本地缓存
		//Username:        "nacos",                              // Nacos账号
		//Password:        "nacos",                              // Nacos密码
		LogDir:   "./tmp/nacos/log",   // 日志目录
		CacheDir: "./tmp/nacos/cache", // 缓存目录
	}

	// 3. 创建配置客户端
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		})
	if err != nil {
		fmt.Printf("创建 Nacos 配置客户端失败: %v\n", err)
		return
	}
	fmt.Printf("成功创建 Nacos 配置客户端: %v\n", client)

	// 4. 获取配置

	content, err := client.GetConfig(vo.ConfigParam{
		DataId: "statstudy",
		Group:  "zg4",
	})
	if err != nil {
		log.Fatal("获取配置失败: ", err)
	}
	fmt.Println("=== 原始配置内容 ===")
	fmt.Println(content)

	//监听配置变更
	err = client.ListenConfig(vo.ConfigParam{
		DataId: "statstudy",
		Group:  "zg4",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("配置发生变更")
			fmt.Println(data)
		},
	})
	if err != nil {
		log.Fatal("监听配置失败", err)
		//保持程序运行以测变变更中

	}
	time.Sleep(10 * time.Minute)
}
