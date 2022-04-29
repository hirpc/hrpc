package uniqueid

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

// 因为snowFlake目的是解决分布式下生成唯一id 所以ID中是包含集群和节点编号在内的
const (
	workerBits uint8 = 10 // 每台机器(节点)的ID位数 10位最大可以有2^10=1024个节点
	numberBits uint8 = 12 // 表示每个集群下的每个节点，1毫秒内可生成的id序号的二进制位数 即每毫秒可生成 2^12-1=4096个唯一ID
	// 这里求最大值使用了位运算，-1 的二进制表示为 1 的补码，感兴趣的同学可以自己算算试试 -1 ^ (-1 << nodeBits) 这里是不是等于 1023
	nodeMax     int64 = -1 ^ (-1 << workerBits) // 节点ID的最大值，用于防止溢出
	numberMax   int64 = -1 ^ (-1 << numberBits) // 同上，用来表示生成id序号的最大值
	timeShift   uint8 = workerBits + numberBits // 时间戳向左的偏移量
	workerShift uint8 = numberBits              // 节点ID向左的偏移量
	// 41位字节作为时间戳数值的话 大约68年就会用完
	// 假如你2010年1月1日开始开发系统 如果不减去2010年1月1日的时间戳 那么白白浪费40年的时间戳啊！
	// 这个一旦定义且开始生成ID后千万不要改了 不然可能会生成相同的ID
	epoch int64 = 1597472019000 // 这个是我在写epoch这个变量时的时间戳(毫秒)
)

// serviceNode should be used for global
var (
	serviceNode *Node
	nodeID      int64
	myIP        string
)

// Node 定义一个Node工作节点所需要的基本参数
type Node struct {
	mu        sync.Mutex // 添加互斥锁 确保并发安全
	timestamp int64      // 记录时间戳
	nodeID    int64      // 该节点的ID
	number    int64      // 当前毫秒已经生成的id序列号(从0开始累加) 1毫秒内最多生成4096个ID
	uniqueID  int64
}

type Option struct {
	NodeID int64
}

func init() {
	var err error
	nodeID, err = createNodeID()
	if err != nil {
		panic(err)
	}
	serviceNode = &Node{nodeID: nodeID}
}

// NewNode 实例化一个工作节点
// 优先级是: 传入参数 > 已存在的NodeID > createNodeID
func New(options ...Option) *Node {
	options = append(options, Option{NodeID: nodeID})

	// Check the nodeID is valid or not
	// 此时options
	for _, v := range options {
		if err := v.validate(); err == nil {
			// 一旦有校验通过的，则赋值 return
			serviceNode = &Node{nodeID: nodeID}
			return serviceNode
		}
	}
	return serviceNode
}

func (n *Node) Load() error {
	return nil
}

func (n *Node) Name() string {
	return "hrpc-uniqueid"
}

func (n *Node) DependsOn() []string {
	return nil
}

func (o Option) validate() error {
	if o.NodeID < 0 || o.NodeID > nodeMax {
		return errors.New("Node ID excess of quantity")
	}
	return nil
}

// newNode will return an unique ID for different uses
func newNode() *Node {
	// 获取生成时的时间戳
	now := time.Now().UnixNano() / 1e6 // 纳秒转毫秒
	if serviceNode.timestamp == now {
		serviceNode.number++
		// 这里要判断，当前工作节点是否在1毫秒内已经生成numberMax个ID
		if serviceNode.number > numberMax {
			// 如果当前工作节点在1毫秒内生成的ID已经超过上限 需要等待1毫秒再继续生成
			for now <= serviceNode.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		// 如果当前时间与工作节点上一次生成ID的时间不一致 则需要重置工作节点生成ID的序号
		serviceNode.number = 0
		serviceNode.timestamp = now // 将机器上一次生成ID的时间更新为当前时间
	}
	// 第一段 now - epoch 为该算法目前已经奔跑了xxx毫秒
	// 如果在程序跑了一段时间修改了epoch这个值 可能会导致生成相同的ID
	serviceNode.uniqueID = int64((now-epoch)<<timeShift | (serviceNode.nodeID << workerShift) | (serviceNode.number))
	return serviceNode
}

func String() string {
	serviceNode.mu.Lock()
	defer serviceNode.mu.Unlock()

	return strconv.FormatInt(newNode().uniqueID, 10)
}

func Number() int64 {
	serviceNode.mu.Lock()
	defer serviceNode.mu.Unlock()

	return newNode().uniqueID
}

func NodeID() int64 {
	return nodeID
}

func IP() string {
	return myIP
}
