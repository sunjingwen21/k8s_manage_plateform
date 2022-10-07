package service

import (
	"context"
	"encoding/json"
	"errors"

	nwv1 "k8s.io/api/networking/v1"

	"github.com/wonderivan/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/**
 * @author 王子龙
 * 时间：2022/9/27 13:25
 */
var Ingress ingress

type ingress struct{}

//定义列表的返回内容，Items是ingress元素列表，Total为ingress元素数量
type IngressesResp struct {
	Items []nwv1.Ingress `json:"items"`
	Total int            `json:"total"`
}

func (s *ingress) toCells(std []nwv1.Ingress) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = ingressCell(std[i])
	}
	return cells
}
func (s *ingress) fromCells(cells []DataCell) []nwv1.Ingress {
	ingresses := make([]nwv1.Ingress, len(cells))
	for i := range cells {
		ingresses[i] = nwv1.Ingress(cells[i].(ingressCell))
	}
	return ingresses
}

//定义ingress的path结构体
type HttpPath struct {
	Path        string        `json:"path"`
	PathType    nwv1.PathType `json:"path_type"`
	ServiceName string        `json:"service_name"`
	ServicePort int32         `json:"service_port"`
}

//定义IngressCreate结构体，用于创建ingress需要的参数属性的定义
type IngressCreate struct {
	Name      string                 `json:"name"`
	Namespace string                 `json:"namespace"`
	Label     map[string]string      `json:"label"`
	Hosts     map[string][]*HttpPath `json:"hosts"`
}

//获取ingress列表，支持过滤、排序、分页
func (i *ingress) GetIngresses(filterName, namespace string, limit, page int) (ingressesResp *IngressesResp, err error) {
	//获取ingressList类型的ingress列表
	ingressList, err := K8s.ClientSet.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取Ingress列表失败，" + err.Error()))
		return nil, errors.New("获取Ingress列表失败，" + err.Error())
	}
	//将ingressList中的ingress列表（Items），放进dataselector对象中，进行排序
	selectableData := &dataSelector{
		GenericDataList: i.toCells(ingressList.Items),
		dataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{
				Name: filterName,
			},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	data := filtered.Sort().Paginate()
	//将[]DataCell类型的ingress列表转为v1.ingress列表
	ingresses := i.fromCells(data.GenericDataList)
	return &IngressesResp{
		Items: ingresses,
		Total: total,
	}, nil
}

//获取ingress详情
func (i *ingress) GetIngressDetail(ingressName, namespace string) (ingress *nwv1.Ingress, err error) {
	ingress, err = K8s.ClientSet.NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Ingress详情失败，" + err.Error()))
		return nil, errors.New("获取Ingress详情失败，" + err.Error())
	}
	return ingress, nil
}

//创建ingress，接收IngressCreate对象
func (i *ingress) CreateIngress(data *IngressCreate) (err error) {
	//声明nwv1.IngressRule和nwv1.HTTPIngressPath变量，后面组装数据用到
	var ingressRules []nwv1.IngressRule
	var httpIngressPATHs []nwv1.HTTPIngressPath
	//将data中的数据组装成nwv1.Ingress对象
	ingress := &nwv1.Ingress{
		//ObjectMeta中定义资源名、命名空间以及标签
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		Status: nwv1.IngressStatus{},
	}
	//第一层for循环是将host组装成nwv1.IngressRule类型性的对象
	//一个host对应一个ingressrule,每个ingressrule中包含一个host和多个path
	for key, value := range data.Hosts {
		ir := nwv1.IngressRule{
			Host: key,
			//这里现将nwv1.HTTPIngressRuleValue类型中的Paths置为空，后面组装好数据再赋值
			IngressRuleValue: nwv1.IngressRuleValue{
				HTTP: &nwv1.HTTPIngressRuleValue{Paths: nil},
			},
		}
		//第二层for循环是将path组装成nwv1.HTTPIngressPath类型的对象
		for _, HttpPath := range value {
			hip := nwv1.HTTPIngressPath{
				Path:     HttpPath.Path,
				PathType: &HttpPath.PathType,
				Backend: nwv1.IngressBackend{
					Service: &nwv1.IngressServiceBackend{
						Name: HttpPath.ServiceName,
						Port: nwv1.ServiceBackendPort{
							Number: HttpPath.ServicePort,
						},
					},
				},
			}
			//将每个hip对象组装成数组
			httpIngressPATHs = append(httpIngressPATHs, hip)
		}
		//给Paths赋值，前面置为空了
		ir.IngressRuleValue.HTTP.Paths = httpIngressPATHs
		//将每个ir对象组装成数组，这个ir对象就是IngressRule，每个元素是一个host和多个path
		ingressRules = append(ingressRules, ir)
	}
	//将ingressRules对象加入到ingress的规则中
	ingress.Spec.Rules = ingressRules
	//创建ingress
	_, err = K8s.ClientSet.NetworkingV1().Ingresses(data.Namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建Ingress失败，" + err.Error()))
		return errors.New("创建Ingress失败，" + err.Error())
	}
	return nil
}

//删除ingress
func (i *ingress) DeleteIngress(ingressName, namespace string) (err error) {
	err = K8s.ClientSet.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除Ingress失败，" + err.Error()))
		return errors.New("删除Ingress失败，" + err.Error())
	}
	return nil
}

//更新ingress
func (i *ingress) UpdateIngress(namespace, content string) (err error) {
	var ingress = &nwv1.Ingress{}
	err = json.Unmarshal([]byte(content), ingress)
	if err != nil {
		logger.Error(errors.New("反序列化失败，" + err.Error()))
		return errors.New("反序列化失败，" + err.Error())
	}
	_, err = K8s.ClientSet.NetworkingV1().Ingresses(namespace).Update(context.TODO(), ingress, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新ingress失败，" + err.Error()))
		return errors.New("更新ingress失败，" + err.Error())
	}
	return nil
}
