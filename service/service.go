package service

import (
	"context"
	"errors"

	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

/**
 * @author 王子龙
 * 时间：2022/9/27 13:25
 */
var Service deployment

type service struct{}

//定义列表的返回内容，Items是deployment元素列表，Total为deployment元素数量
type ServicesResp struct {
	Items []corev1.Service `json:"items"`
	Total int              `json:"total"`
}

//定义DeploysNp类型，用于返回namespace中deployment的数量
type ServicesNp struct {
	Namespace  string `json:"namespace"`
	ServiceNum int    `json:"Service_num"`
}

func (s *service) toCells(std []corev1.Service) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = serviceCell(std[i])
	}
	return cells
}
func (s *service) fromCells(cells []DataCell) []corev1.Service {
	services := make([]corev1.Service, len(cells))
	for i := range cells {
		services[i] = corev1.Service(cells[i].(serviceCell))
	}
	return services
}

//定义ServiceCreate结构体，用于创建service需要的参数属性的定义
type ServiceCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Type          string            `json:"type"`
	ContainerPort int32             `json:"container_port"`
	Port          int32             `json:"port"`
	NodePort      int32             `json:"node_port"`
	Label         map[string]string `json:"label"`
}

//创建service，接收ServiceCreate对象
func (s *service) CreateService(data *ServiceCreate) (err error) {
	//将data中的数据组装成corev1.Service对象
	service := &corev1.Service{
		//ObjectMeta中定义资源名、命名空间以及标签
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		//Spec中定义类型，端口，选择器
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(data.Type),
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     data.Port,
					Protocol: "TCP",
					TargetPort: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			Selector: data.Label,
		},
	}
	//默认ClusterIP,这里是判断NodePort，添加配置
	if data.NodePort != 0 && data.Type == "NodePort" {
		service.Spec.Ports[0].NodePort = data.NodePort
	}
	//创建Service
	_, err = K8s.ClientSet.CoreV1().Services(data.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建Service失败，" + err.Error()))
		return errors.New("创建Service失败，" + err.Error())
	}
	return nil
}
