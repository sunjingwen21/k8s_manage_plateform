package service

import (
	"encoding/json"
	"fmt"
	"io"
	"k8s-platform/config"
	"log"
	"net/http"
	"time"

	"k8s.io/client-go/kubernetes/scheme"

	"github.com/gorilla/websocket"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

/**
 * @author 王子龙
 * 时间：2022/10/2 10:39
 */
var Terminal terminal

type terminal struct{}

//定义websocket的handler方法
func (t *terminal) WsHandler(w http.ResponseWriter, r *http.Request) {
	//加载K8s配置
	conf, err := clientcmd.BuildConfigFromFlags("", config.Kubeconfig)
	if err != nil {
		logger.Error("创建k8s配置失败，" + err.Error())
	}
	//解析form入参，获取namespacce、podName、containerName参数
	if err := r.ParseForm(); err != nil {
		return
	}
	namespace := r.Form.Get("namespace")
	podName := r.Form.Get("pod_name")
	containerName := r.Form.Get("container_name")
	logger.Info("exec pod: %s,container: %s,namespace: %s\n", podName, containerName, namespace)
	//new一个TerminalSession类型的pty实例
	pty, err := NewTerminalSession(w, r, nil)
	if err != nil {
		logger.Error("get pty failed: %v\n", err)
		return
	}
	//处理关闭
	defer func() {
		logger.Info("close session.")
		pty.Close()
	}()
	//初始化pod所在的corev1资源组
	//PodExecOptions struct包括Container stdout Command等结构
	//scheme.ParameterCodec 应该是pod的GVK(GroupVersion &Kind)之类的
	//URL长相：
	//https://192.168.1.11:6443/api/v1/namespaces/default/pods/nginx-wf2-778d88d7c-
	//7rmsk/exec?command=%2Fbin%2Fbash&container=nginx-
	//wf2&stderr=true&stdin=true&stdout=true&tty=true
	req := K8s.ClientSet.CoreV1().RESTClient().Post().Resource("pods").
		Name(podName).Namespace(namespace).SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   []string{"/bin/bash"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)
	fmt.Println("req.URL()", req.URL())
	//remotecommand主要实现了http转SPDY添加X-Stream-Protocol-Version相关header并发送请求
	executor, err := remotecommand.NewSPDYExecutor(conf, "POST", req.URL())
	if err != nil {
		return
	}
	//建立链接之后从请求的stream中发送、读取数据
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		TerminalSizeQueue: pty,
		Tty:               pty.tty,
	})
	if err != nil {
		msg := fmt.Sprintf("Exec to pod error! err: %v", err)
		logger.Info(msg)
		//将报错返回出去
		pty.Write([]byte(msg))
		//标记退出stream流
		pty.Done()
	}
}

const END_OF_TRANSMISSION = "\u0004" //终止符
// TerminalMessage is the messaging protocol between ShellController and TerminalSession.
//TerminalMessage定义了终端和容器shell交互内容的格式
//Operation是操作类型
//Data是具体数据内容
//Rows和Cols可以理解为终端的行数和列数，也就是宽、高
type TerminalMessage struct {
	Operation string `json:"operation"`
	Data      string `json:"data"`
	Rows      uint16 `json:"rows"`
	Cols      uint16 `json:"cols"`
}

//初始化一个websocket.Upgrader类型的对象，用于http协议升级为websocket协议
var upgrader = func() websocket.Upgrader {
	upgrader := websocket.Upgrader{}
	upgrader.HandshakeTimeout = time.Second * 2
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	return upgrader
}()

// TerminalSession implements PtyHandler
//定义TerminalSession结构体，实现PtyHandler接口
//wsConn是websocket连接
//sizeChan用来定义终端输入和输出的宽和高
//doneChan用于标记退出终端
type TerminalSession struct {
	wsConn   *websocket.Conn
	sizeChan chan remotecommand.TerminalSize
	doneChan chan struct{}
	tty      bool
}

//该方法用于升级http协议至websocket，并new一个TerminalSession类型的对象返回
func NewTerminalSession(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*TerminalSession, error) {
	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}
	session := &TerminalSession{
		wsConn:   conn,
		sizeChan: make(chan remotecommand.TerminalSize),
		doneChan: make(chan struct{}),
		tty:      true,
	}
	return session, nil
}

// Done done, must call Done() before connection close, or Next() would not exits.
//关闭doneChan，关闭后触发退出终端
func (t *TerminalSession) Done() {
	close(t.doneChan)
}

// Next called in a loop from remotecommand as long as the process is running
//获取web端是否resize，以及是否退出终端
func (t *TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.sizeChan:
		return &size
	case <-t.doneChan:
		return nil
	}
}

//用于读取web端的输入，接收web端输入的指令内容
func (t *TerminalSession) Read(p []byte) (int, error) {
	_, message, err := t.wsConn.ReadMessage()
	if err != nil {
		log.Printf("read message err: %v", err)
		return copy(p, END_OF_TRANSMISSION), err
	}
	var msg TerminalMessage
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		log.Printf("read parse message err: %v", err)
		//return 0,nil
		return copy(p, END_OF_TRANSMISSION), err
	}
	switch msg.Operation {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}
		return 0, nil
	case "ping":
		return 0, nil
	default:
		log.Printf("unknown message type '%s'", msg.Operation)
		//return 0,nil
		return copy(p, END_OF_TRANSMISSION), fmt.Errorf("unknown message type '%s'", msg.Operation)
	}
}

//用于向web端输出，接收web端的指令后，将结果返回出去
func (t *TerminalSession) Write(p []byte) (int, error) {
	msg, err := json.Marshal(TerminalMessage{
		Operation: "stdout",
		Data:      string(p),
	})
	if err != nil {
		log.Printf("write parse message err: %v", err)
		return 0, err
	}
	if err := t.wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Printf("write message err: %v", err)
		return 0, err
	}
	return len(p), nil
}

//用于关闭websocket连接
func (t *TerminalSession) Close() error {
	return t.wsConn.Close()
}

//以下为不知道什么用的方法

//// NewTerminalSessionWs create TerminalSession
//func NewTerminalSessionWs(conn *websocket.Conn) *TerminalSession {
//	return &TerminalSession{
//		wsConn:   conn,
//		tty:      true,
//		sizeChan: make(chan remotecommand.TerminalSize),
//		doneChan: make(chan struct{}),
//	}
//}

// Stdin ...
func (t *TerminalSession) Stdin() io.Reader {
	return t
}

// Stdout ...
func (t *TerminalSession) Stdout() io.Writer {
	return t
}

// Stderr ...
func (t *TerminalSession) Stderr() io.Writer {
	return t
}

// Tty ...
func (t *TerminalSession) Tty() bool {
	return t.tty
}
