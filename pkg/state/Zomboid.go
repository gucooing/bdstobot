package state

import (
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"
)

// MotdPM信息
type MotdPMInfo struct {
	Status    string `json:"status"`     //服务器状态
	Host      string `json:"host"`       //服务器Host
	Name      string `json:"name"`       //Motd信息
	Version   string `json:"version"`    //支持的游戏版本
	Online    int    `json:"online"`     //在线人数
	Max       int    `json:"max"`        //最大在线人数
	Pvp       bool   `json:"pvp"`        //是否开启pvp
	Open      bool   `json:"open"`       //是否需要密码
	LevelName string `json:"level_name"` //存档名字
	GameMod   string `json:"gamemode"`   //游戏mod
	Ping      int64  `json:"ping"`       //连接延迟
}

func MotdPm(Host string) (*MotdPMInfo, error) {
	if Host == "" {
		MotdInfo := &MotdPMInfo{
			Status: "offline",
		}
		return MotdInfo, nil
	}
	// 创建第一次udp连接
	socket, err := net.Dial("udp", Host)
	if err != nil {
		MotdInfo := &MotdPMInfo{
			Status: "offline",
		}
		return MotdInfo, err
	}
	defer socket.Close()
	// 第一次发送数据
	time11 := time.Now().UnixNano() / 1e6 //记录发送时间
	senddata11, _ := hex.DecodeString("ffffffff54536f7572636520456e67696e6520517565727900")
	_, err = socket.Write(senddata11)
	if err != nil {
		MotdInfo := &MotdPMInfo{
			Status: "offline",
		}
		return MotdInfo, err
	}
	// 第一次接收数据
	UDPdata11 := make([]byte, 9)
	socket.SetReadDeadline(time.Now().Add(5 * time.Second)) //设置读取五秒超时
	_, err = socket.Read(UDPdata11)
	if err != nil {
		MotdInfo := &MotdPMInfo{
			Status: "offline",
		}
		return MotdInfo, err
	}
	// 二次发送数据，将接收的数据按要求进行处理
	hexStr1 := hex.EncodeToString(UDPdata11)
	hexStr2 := "ffffffff54536f7572636520456e67696e6520517565727900"
	hexStrCombined := hexStr2 + hexStr1[len(hexStr1)-8:]
	// 转换回 16 进制字符串
	senddata12, _ := hex.DecodeString(hexStrCombined)
	_, err = socket.Write(senddata12)
	if err != nil {
		MotdInfo := &MotdPMInfo{
			Status: "offline",
		}
		return MotdInfo, err
	}
	// 第二次接收数据
	UDPdata12 := make([]byte, 4048)
	socket.SetReadDeadline(time.Now().Add(5 * time.Second)) // 设置读取五秒超时
	_, err = socket.Read(UDPdata12)
	if err != nil {
		MotdInfo := &MotdPMInfo{
			Status: "offline",
		}
		return MotdInfo, err
	}
	time12 := time.Now().UnixNano() / 1e6 //记录接收时间

	//第二次udp连接
	socket2, err := net.Dial("udp", Host)
	if err != nil {
		MotdInfo := &MotdPMInfo{
			Status: "offline",
		}
		return MotdInfo, err
	}
	defer socket2.Close()
	// 第一次发送数据
	time21 := time.Now().UnixNano() / 1e6 //记录发送时间
	senddata21, _ := hex.DecodeString("ffffffff5600000000")
	_, err = socket2.Write(senddata21)
	if err != nil {
		MotdInfo := &MotdPMInfo{
			Status: "offline",
		}
		return MotdInfo, err
	}
	// 第一次接收数据
	UDPdata21 := make([]byte, 9)
	socket2.SetReadDeadline(time.Now().Add(5 * time.Second)) //设置读取五秒超时
	_, err = socket2.Read(UDPdata21)
	if err != nil {
		MotdInfo := &MotdPMInfo{
			Status: "offline",
		}
		return MotdInfo, err
	}
	hexBytes := []byte(hex.EncodeToString(UDPdata21))
	hexBytes[8] = '5'
	hexBytes[9] = '6'
	// 第二次发送数据
	senddata22, _ := hex.DecodeString(string(hexBytes))
	_, err = socket2.Write(senddata22)
	if err != nil {
		MotdInfo := &MotdPMInfo{
			Status: "offline",
		}
		return MotdInfo, err
	}
	// 第二次接收数据
	UDPdata22 := make([]byte, 4048)
	socket2.SetReadDeadline(time.Now().Add(5 * time.Second)) //设置读取五秒超时
	_, err = socket2.Read(UDPdata22)
	if err != nil {
		MotdInfo := &MotdPMInfo{
			Status: "offline",
		}
		return MotdInfo, err
	}
	time22 := time.Now().UnixNano() / 1e6 //记录接收时间
	// 数据处理 处理第一次udp连接的最后数据
	searchBytes := []byte{80, 114, 111, 106, 101, 99, 116, 32, 90, 111, 109, 98, 111, 105, 100}
	offset := 3
	pos := search(UDPdata12, searchBytes)
	Online := UDPdata12[pos+len(searchBytes)+offset]
	Max := UDPdata12[pos+len(searchBytes)+offset+1]
	MotdData := strings.Split(string(UDPdata12[6:]), string(hexsr("00")))
	// 处理第二次udp连接的最后数据
	MotdData2 := strings.Split(string(UDPdata22), string(hexsr("00")))
	pvp := bools(MotdData2[12])
	open := bools(MotdData2[8])
	// motd数据解析
	if err == nil {
		MotdInfo := &MotdPMInfo{
			Status:    "online",
			Host:      Host,
			Name:      MotdData[0],
			Version:   MotdData2[14],
			Online:    int(Online),
			Max:       int(Max),
			Pvp:       pvp,
			Open:      open,
			LevelName: MotdData[1],
			GameMod:   MotdData2[6],
			Ping:      ((time12 - time11) + (time22 - time21)) / 4,
		}
		return MotdInfo, nil
	}
	MotdInfo := &MotdPMInfo{
		Status: "offline",
	}
	return MotdInfo, err
}

func hexsr(data string) []byte {
	Data, err := hex.DecodeString(data)
	if err != nil {
		fmt.Println("转换失败:", err)
	}
	return Data
}

func search(data []byte, searchBytes []byte) int {
	for i := 0; i < len(data); i++ {
		if data[i] == searchBytes[0] && i+len(searchBytes) <= len(data) {
			match := true
			for j := 1; j < len(searchBytes); j++ {
				if data[i+j] != searchBytes[j] {
					match = false
					break
				}
			}
			if match {
				return i
			}
		}
	}
	return -1
}

func bools(data string) bool {
	if data == "0" {
		return true
	} else {
		return false
	}
}
