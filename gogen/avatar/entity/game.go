package entity

import "time"

// GameType GameType
type GameType string

// GameDetail detail info
type GameDetail struct {
	GameDetailUpdateTime  time.Time  `json:"gameDetailUpdateTime" bson:"gameDetailUpdateTime"`
	Name                  string     `json:"name" bson:"name"`
	IconID                string     `json:"iconID" bson:"iconID"`
	IconURL               string     `json:"iconURL" bson:"-"`
	CoverID               string     `json:"coverID" bson:"coverID"`
	CoverURL              string     `json:"coverURL" bson:"-"`
	ThumbnailIDs          []string   `json:"thumbnailIDs" bson:"thumbnailIDs"`
	ThumbnailURLs         []string   `json:"thumbnailURLs" bson:"-"`
	Description           string     `json:"description" bson:"description"`
	Tags                  []string   `json:"tags" bson:"tags"`
	AllowSinglePlayerMode bool       `json:"allowSinglePlayerMode" bson:"allowSinglePlayerMode"` // 该游戏允许单人模式
	MinPlayers            int        `json:"minPlayers" bson:"minPlayers"`
	MaxPlayers            int        `json:"maxPlayers" bson:"maxPlayers"`
	GameTypes             []GameType `json:"gameTypes" bson:"gameTypes"` // 游戏类型
	// TODO 旧的字段@ 在创造者中心一期要求封面跟视频都可以传递多个
	VideoURL         string   `json:"videoURL" bson:"videoURL"`                 // 这些游戏详情数据审核通过后会覆盖线上游戏
	VideoStoreID     string   `json:"videoStoreID" bson:"videoStoreID"`         // 这些游戏详情数据审核通过后会覆盖线上游戏
	VideoStoreIDList []string `json:"videoStoreIDList" bson:"videoStoreIDList"` // 这些游戏详情数据审核通过后会覆盖线上游戏
	// TODO 旧的字段@ 在创造者中心一期要求封面跟视频都可以传递多个
	VideoCoverURL         string   `json:"videoCoverURL" bson:"videoCoverURL"`                 // 这些游戏详情数据审核通过后会覆盖线上游戏
	VideoCoverStoreID     string   `json:"videoCoverStoreID" bson:"videoCoverStoreID"`         // 这些游戏详情数据审核通过后会覆盖线上游戏
	VideoCoverStoreIDList []string `json:"videoCoverStoreIDList" bson:"videoCoverStoreIDList"` // 这些游戏详情数据审核通过后会覆盖线上游戏
	DeveloperID           string   `json:"developerID" bson:"developerID"`                     // 这些游戏详情数据审核通过后会覆盖线上游戏
	DeveloperSpeaking     string   `json:"developerSpeaking" bson:"developerSpeaking"`         // 这些游戏详情数据审核通过后会覆盖线上游戏
	UpdateLog             string   `json:"updateLog" bson:"updateLog"`                         // 更新日志
}
