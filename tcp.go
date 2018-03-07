package main

import (
	"fmt"
)

/**
 *  tcp 启动 三次 握手  关闭四次
 *  慢启动 + 拥塞窗口
 *  门限控制 < 拥塞窗口 进入拥塞避免
 *  在发送数据后 检测 ack 确认超时 或者 收到重复ack 确认 操作门限值 和 拥塞窗口值 来限流
 */

func main() {

}
