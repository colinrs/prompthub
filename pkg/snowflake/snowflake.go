package snowflake

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	snowflakeExternal "github.com/bwmarrin/snowflake"
)

type Snowflake interface {
	GetNodeID() int64
	NextID(ctx context.Context) (int64, error)
	NextIDs(ctx context.Context, n int64) ([]int64, error)
}

type snowflake struct {
	nodeID int64
	idGen  *snowflakeExternal.Node
}

func NewSnowflake(nodeID int64) Snowflake {
	idGen, err := snowflakeExternal.NewNode(nodeID)
	logx.Must(err)
	return &snowflake{
		nodeID: nodeID,
		idGen:  idGen,
	}
}

func (s *snowflake) NextID(ctx context.Context) (int64, error) {
	return s.idGen.Generate().Int64(), nil
}

func (s *snowflake) NextIDs(ctx context.Context, n int64) ([]int64, error) {
	ids := make([]int64, 0, n)
	for i := 0; i < int(n); i++ {
		ids = append(ids, s.idGen.Generate().Int64()) //nolint:s.idGen.Generate()
	}
	return ids, nil
}

func (s *snowflake) GetNodeID() int64 {
	return s.nodeID
}

type SnowflakeInfo struct {
	NodeID    int
	Epoch     int64
	Time      string
	Timestamp int64
	Sequence  int
}

// ParseSnowflakeID 解析 Snowflake ID 的详细信息
func ParseSnowflakeID(id int64) *SnowflakeInfo {
	// Snowflake 算法的起始时间戳
	const epoch = 1288834974657 // 推特的默认纪元时间戳
	// 提取各个部分
	// Snowflake ID 结构：1位符号位 + 41位时间戳 + 10位机器ID + 12位序列号
	timePart := (id >> 22) & 0x1FFFFFFFFFF
	nodePart := (id >> 12) & 0x3FF
	sequencePart := id & 0xFFF
	// 计算时间
	timestamp := timePart + epoch
	timeObj := time.Unix(timestamp/1000, (timestamp%1000)*1000000)
	return &SnowflakeInfo{
		NodeID:    int(nodePart),
		Epoch:     epoch,
		Timestamp: timestamp,
		Time:      timeObj.Format("2006-01-02 15:04:05.000"),
		Sequence:  int(sequencePart),
	}
}
