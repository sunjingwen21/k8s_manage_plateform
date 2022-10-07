package service

import (
	"sort"
	"strings"
	"time"

	nwv1 "k8s.io/api/networking/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

/**
 * @author 王子龙
 * 时间：2022/9/22 9:50
 */

//dataSelect 用于封装排序、过滤、分页的数据类型
type dataSelector struct {
	GenericDataList []DataCell
	dataSelectQuery *DataSelectQuery //过滤和分页的属性
}

//DataCell接口，用于各种资源list的类型转换，转换后可以使用dataSelector的自定义排序方法
type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

//DataSelectQuery定义过滤和分页的属性，过滤：Name，分页：Limit和Page
//Limit是单页的数据条数
//Page是第几页
type DataSelectQuery struct {
	FilterQuery   *FilterQuery
	PaginateQuery *PaginateQuery
}
type FilterQuery struct {
	Name string
}
type PaginateQuery struct {
	Limit int
	Page  int
}

//实现自定义结构的排序，需要重写Len、Swap、Less方法
//Len方法用于获取数组长度
func (d *dataSelector) Len() int {
	return len(d.GenericDataList)
}

//Swap方法用于数组中的元素在比较大小后的位置交换，可定义升序或降序
func (d *dataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

//Less方法用于定义数组中元素排序的“大小”的比较方式
func (d *dataSelector) Less(i, j int) bool {
	a := d.GenericDataList[i].GetCreation()
	b := d.GenericDataList[j].GetCreation()
	return b.Before(a)
}

//重写以上3个方法使用sort.Sort进行排序
func (d *dataSelector) Sort() *dataSelector {
	sort.Sort(d)
	return d
}

//Filter方法用于过滤元素，比较元素的Name属性，若包含，再返回
func (d *dataSelector) Filter() *dataSelector {
	//若Name的传参为空，则返回所有元素
	if d.dataSelectQuery.FilterQuery.Name == "" {
		return d
	}
	//若Name的传参不为空，则返回元素名中包含Name的所有元素
	filteredList := []DataCell{}
	for _, value := range d.GenericDataList {
		matches := true
		objName := value.GetName()
		//判断字符串s中是否包含子串str
		if !strings.Contains(objName, d.dataSelectQuery.FilterQuery.Name) {
			matches = false
			continue
		}
		if matches {
			filteredList = append(filteredList, value)
		}
	}
	d.GenericDataList = filteredList
	return d
}

//Paginate方法用于数组分页，根据Limit和Page的传参，返回数据
func (d *dataSelector) Paginate() *dataSelector {
	limit := d.dataSelectQuery.PaginateQuery.Limit
	page := d.dataSelectQuery.PaginateQuery.Page
	//验证参数合法，若参数不合法，则返回所有数据
	if limit <= 0 || page <= 0 {
		return d
	}
	//举例：25个元素的数组，limit是10，page是3，startIndex是20，endIndex是30（实际上endIndex是25）
	startIndex := limit * (page - 1)
	endIndex := limit * page
	//处理最后一页，这时候就把endIndex由30改为25了
	if len(d.GenericDataList) < endIndex {
		endIndex = len(d.GenericDataList)
	}
	d.GenericDataList = d.GenericDataList[startIndex:endIndex]
	return d
}

//定义podCell类型，实现GetCreation和getName方法后，可进行类型转换
type podCell corev1.Pod

func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}
func (p podCell) GetName() string {
	return p.Name
}

//实现Deployment的DataCell接口
type deploymentCell appsv1.Deployment //appsv1 "k8s.io/api/apps/v1"
func (d deploymentCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}
func (d deploymentCell) GetName() string {
	return d.Name
}

//实现Service的DataCell接口
type serviceCell corev1.Service //corev1 "k8s.io/api/core/v1"
func (s serviceCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}
func (s serviceCell) GetName() string {
	return s.Name
}

//实现Ingress的DataCell接口
type ingressCell nwv1.Ingress //nwv1 "k8s.io/api/networking/v1"
func (i ingressCell) GetCreation() time.Time {
	return i.CreationTimestamp.Time
}
func (i ingressCell) GetName() string {
	return i.Name
}

type namespaceCell corev1.Namespace

func (n namespaceCell) GetCreation() time.Time {
	return n.CreationTimestamp.Time
}

func (n namespaceCell) GetName() string {
	return n.Name
}

type configMapCell corev1.ConfigMap

func (c configMapCell) GetCreation() time.Time {
	return c.CreationTimestamp.Time
}

func (c configMapCell) GetName() string {
	return c.Name
}

type daemonSetCell appsv1.DaemonSet

func (d daemonSetCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}

func (d daemonSetCell) GetName() string {
	return d.Name
}

type statefulSetCell appsv1.StatefulSet

func (s statefulSetCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s statefulSetCell) GetName() string {
	return s.Name
}

type nodeCell corev1.Node

func (n nodeCell) GetCreation() time.Time {
	return n.CreationTimestamp.Time
}

func (n nodeCell) GetName() string {
	return n.Name
}

type pvCell corev1.PersistentVolume

func (p pvCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p pvCell) GetName() string {
	return p.Name
}

type secretCell corev1.Secret

func (s secretCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s secretCell) GetName() string {
	return s.Name
}

type pvcCell corev1.PersistentVolumeClaim

func (p pvcCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p pvcCell) GetName() string {
	return p.Name
}
