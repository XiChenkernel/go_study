package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"rankingList/rankingList/cache"
	"sort"
	"strconv"
)

func main() {
	cache.Redis() // 连接redis
	server := gin.Default()

	server.GET("show/:id", ShowViewCount) // 将增加播放量接口注册在 localhost:3000/show 地址

	server.GET("rank", GetRank) // 将排行榜注册在 localhost:3000/rank 地址

	server.Run(":9000") // 运行在本地3000端口
}

type Share struct {
	Id        string
	ViewCount int64
}

// 获取播放量函数
func (share *Share) GetViewCount() (num int64) {
	countStr, _ := cache.RedisClient.Get(context.Background(), cache.ShareKey(share.Id)).Result()
	if countStr == "" {
		return 0
	}
	num, _ = strconv.ParseInt(countStr, 10, 64)
	return
}

// AddViewCount 增加播放量函数
func (share *Share) AddViewCount() {
	// 增加播放量
	cache.RedisClient.Incr(context.Background(), cache.ShareKey(share.Id))
	// 增加在排行榜中的播放量
	cache.RedisClient.ZIncrBy(context.Background(), cache.DailyRankKey, 1, share.Id)
}

func ShowViewCount(ctx *gin.Context) {
	id := ctx.Param("id")
	share := Share{
		Id: id,
	}
	// 增加播放量
	share.AddViewCount()
	// 获取Redis数据
	share.ViewCount = share.GetViewCount()
	ctx.JSON(200, share)
}

func GetRank(ctx *gin.Context) {
	shares := make([]Share, 0, 16)

	// 从Redis中获取排行榜
	shareRank, err := cache.RedisClient.ZRevRange(context.Background(), cache.DailyRankKey, 0, 9).Result()
	if err != nil {
		ctx.JSON(200, err.Error())
		return
	}

	// 获取排行榜内对应排名的播放量
	if len(shareRank) > 0 {
		for _, shareId := range shareRank {
			share := Share{
				Id: shareId,
			}
			share.ViewCount = share.GetViewCount()
			shares = append(shares, share)
		}
	}

	// 填充空的排行榜排名至十个
	emptyShare := Share{
		Id:        "虚位以待",
		ViewCount: 0,
	}
	for len(shares) < 10 {
		shares = append(shares, emptyShare)
	}

	// 由于获取排行榜时有可能排行榜的Zset发生变动，需要按照确定的播放数重新排名一次
	sort.Slice(shares, func(i, j int) bool {
		return shares[i].ViewCount > shares[j].ViewCount
	})

	ctx.JSON(200, shares)
}
