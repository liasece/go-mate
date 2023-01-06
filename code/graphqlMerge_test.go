package code

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func fixTestingBlockIDAndRegIndexAndBlockParser(want *Block, got *Block) {
	got.ID = want.ID
	got.RegOriginStrings = want.RegOriginStrings
	got.RegOriginIndexes = want.RegOriginIndexes
	got.SubOriginIndex = want.SubOriginIndex
	got.BlockParser = want.BlockParser
	for i, v := range got.SubList {
		for j, vv := range v {
			if len(want.SubList) > i && len(want.SubList[i]) > j {
				fixTestingBlockIDAndRegIndexAndBlockParser(want.SubList[i][j], vv)
			}
		}
	}
}

func getMergeTestJsonString(old *Block, income *Block, got *Block, expect *Block) string {
	value := []map[string]interface{}{
		{"old": old},
		{"income": income},
		{"got": got},
		{"expect": expect},
	}
	jsStr, err := json.MarshalIndent(value, "", "\t")
	if err != nil {
		return ""
	}
	return string(jsStr)
}

var graphqlContent1 string = `
type GameEntry implements Node {
  id: ID!
  name: String!
  channelID: String!
  createAt: Timestamp!
  updateAt: Timestamp!
  """
  可见性 0: 只有本人可见 1: 公开 2: 只有所属的队伍可见。
  经过与讨论，在 dev channel 中，可见性有存在的必要；但是在 prod channel 中，可见性其实表现为是否上架，如果所有人可见就是已上架，如果只有本人可见就是未上架。
  所以这个可见性在 workshop 游戏中是其本身含义，在线上游戏中表示为这个游戏是否上架。
  """
  visibilityType: Int!
  ownerID: String!
  groupID: String!
  oldGameID: String!
  oldGameVersion: String!
  detailID: String!
  indexWeight: Int!

  detail: GameDetail @goField(forceResolver: true)
  game: Game @goField(forceResolver: true)
}

input GameEntryUpdater {
  justDelete: Boolean
}

extend type Query {
  gameEntry(id: ID!): GameEntry! @HasPermission(auth: { prefixAny: [GAME, PLAYER] })
  gameEntries(filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!): GameEntryConnection!
    @HasPermission(auth: { prefixAny: [GAME, PLAYER] })
  searchGameEntry(filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!): GameEntryConnection!
    @HasPermission(auth: { prefixAny: [GAME, PLAYER] })
}
`

func TestGraphqlBlockFromString(t *testing.T) {
	c := NewGraphqlCodeBlockParser()
	type args struct {
		p       *Block
		content string
	}
	tests := []struct {
		name string
		args args
		want *Block
	}{
		{
			name: "test1",
			args: args{
				content: graphqlContent1,
			},
			want: &Block{
				OriginString:    graphqlContent1,
				SubOriginString: []string{graphqlContent1},
				Type:            BlockType{SubMergeType: []bool{true}, SubsSeparator: "\n"},
				SubList: [][]*Block{
					{
						{
							Key:             "GameEntry",
							Type:            GraphqlBlockTypeType,
							OriginString:    "type GameEntry implements Node {\n  id: ID!\n  name: String!\n  channelID: String!\n  createAt: Timestamp!\n  updateAt: Timestamp!\n  \"\"\"\n  可见性 0: 只有本人可见 1: 公开 2: 只有所属的队伍可见。\n  经过与讨论，在 dev channel 中，可见性有存在的必要；但是在 prod channel 中，可见性其实表现为是否上架，如果所有人可见就是已上架，如果只有本人可见就是未上架。\n  所以这个可见性在 workshop 游戏中是其本身含义，在线上游戏中表示为这个游戏是否上架。\n  \"\"\"\n  visibilityType: Int!\n  ownerID: String!\n  groupID: String!\n  oldGameID: String!\n  oldGameVersion: String!\n  detailID: String!\n  indexWeight: Int!\n\n  detail: GameDetail @goField(forceResolver: true)\n  game: Game @goField(forceResolver: true)\n}\n",
							SubOriginString: []string{"  id: ID!\n  name: String!\n  channelID: String!\n  createAt: Timestamp!\n  updateAt: Timestamp!\n  \"\"\"\n  可见性 0: 只有本人可见 1: 公开 2: 只有所属的队伍可见。\n  经过与讨论，在 dev channel 中，可见性有存在的必要；但是在 prod channel 中，可见性其实表现为是否上架，如果所有人可见就是已上架，如果只有本人可见就是未上架。\n  所以这个可见性在 workshop 游戏中是其本身含义，在线上游戏中表示为这个游戏是否上架。\n  \"\"\"\n  visibilityType: Int!\n  ownerID: String!\n  groupID: String!\n  oldGameID: String!\n  oldGameVersion: String!\n  detailID: String!\n  indexWeight: Int!\n\n  detail: GameDetail @goField(forceResolver: true)\n  game: Game @goField(forceResolver: true)\n"},
							SubList: [][]*Block{
								{
									{
										Key:             "id",
										Type:            GraphqlBlockTypeTypeField,
										OriginString:    "  id: ID!\n",
										SubOriginString: []string{"", "", ""},
										SubList:         [][]*Block{nil, nil, nil},
									},
									{
										Key:             "name",
										Type:            GraphqlBlockTypeTypeField,
										OriginString:    "  name: String!\n",
										SubOriginString: []string{"", "", ""},
										SubList:         [][]*Block{nil, nil, nil},
									},
									{
										Key:             "channelID",
										Type:            GraphqlBlockTypeTypeField,
										OriginString:    "  channelID: String!\n",
										SubOriginString: []string{"", "", ""},
										SubList:         [][]*Block{nil, nil, nil},
									},
									{
										Key:             "createAt",
										Type:            GraphqlBlockTypeTypeField,
										OriginString:    "  createAt: Timestamp!\n",
										SubOriginString: []string{"", "", ""},
										SubList:         [][]*Block{nil, nil, nil},
									},
									{
										Key:             "updateAt",
										Type:            GraphqlBlockTypeTypeField,
										OriginString:    "  updateAt: Timestamp!\n",
										SubOriginString: []string{"", "", ""},
										SubList:         [][]*Block{nil, nil, nil},
									},
									{
										Key:             "visibilityType",
										Type:            GraphqlBlockTypeTypeField,
										OriginString:    "  \"\"\"\n  可见性 0: 只有本人可见 1: 公开 2: 只有所属的队伍可见。\n  经过与讨论，在 dev channel 中，可见性有存在的必要；但是在 prod channel 中，可见性其实表现为是否上架，如果所有人可见就是已上架，如果只有本人可见就是未上架。\n  所以这个可见性在 workshop 游戏中是其本身含义，在线上游戏中表示为这个游戏是否上架。\n  \"\"\"\n  visibilityType: Int!\n",
										SubOriginString: []string{"  \"\"\"\n  可见性 0: 只有本人可见 1: 公开 2: 只有所属的队伍可见。\n  经过与讨论，在 dev channel 中，可见性有存在的必要；但是在 prod channel 中，可见性其实表现为是否上架，如果所有人可见就是已上架，如果只有本人可见就是未上架。\n  所以这个可见性在 workshop 游戏中是其本身含义，在线上游戏中表示为这个游戏是否上架。\n  \"\"\"\n", "", ""},
										SubList: [][]*Block{
											{
												{
													Key:             "可见性 0: 只有本人可见 1: 公开 2: 只有所属的队伍可见。",
													Type:            GraphqlBlockExplain,
													OriginString:    "  \"\"\"\n  可见性 0: 只有本人可见 1: 公开 2: 只有所属的队伍可见。\n  经过与讨论，在 dev channel 中，可见性有存在的必要；但是在 prod channel 中，可见性其实表现为是否上架，如果所有人可见就是已上架，如果只有本人可见就是未上架。\n  所以这个可见性在 workshop 游戏中是其本身含义，在线上游戏中表示为这个游戏是否上架。\n  \"\"\"",
													SubOriginString: []string{},
													SubList:         [][]*Block{},
												},
											},
											nil,
											nil,
										},
									},
									{
										Key:             "ownerID",
										Type:            GraphqlBlockTypeTypeField,
										OriginString:    "  ownerID: String!\n",
										SubOriginString: []string{"", "", ""},
										SubList:         [][]*Block{nil, nil, nil},
									},
									{
										Key:             "groupID",
										Type:            GraphqlBlockTypeTypeField,
										OriginString:    "  groupID: String!\n",
										SubOriginString: []string{"", "", ""},
										SubList:         [][]*Block{nil, nil, nil},
									},
									{
										Key:             "oldGameID",
										Type:            GraphqlBlockTypeTypeField,
										OriginString:    "  oldGameID: String!\n",
										SubOriginString: []string{"", "", ""},
										SubList:         [][]*Block{nil, nil, nil},
									},
									{
										Key:             "oldGameVersion",
										Type:            GraphqlBlockTypeTypeField,
										OriginString:    "  oldGameVersion: String!\n",
										SubOriginString: []string{"", "", ""},
										SubList:         [][]*Block{nil, nil, nil},
									},
									{
										Key:             "detailID",
										Type:            GraphqlBlockTypeTypeField,
										OriginString:    "  detailID: String!\n",
										SubOriginString: []string{"", "", ""},
										SubList:         [][]*Block{nil, nil, nil},
									},
									{
										Key:             "indexWeight",
										Type:            GraphqlBlockTypeTypeField,
										OriginString:    "  indexWeight: Int!\n",
										SubOriginString: []string{"", "", ""},
										SubList:         [][]*Block{nil, nil, nil},
									},
									{
										Key:             "detail",
										Type:            GraphqlBlockTypeTypeField,
										OriginString:    "  detail: GameDetail @goField(forceResolver: true)\n",
										SubOriginString: []string{"", "", ""},
										SubList:         [][]*Block{nil, nil, nil},
									},
									{
										Key:             "game",
										Type:            GraphqlBlockTypeTypeField,
										OriginString:    "  game: Game @goField(forceResolver: true)\n",
										SubOriginString: []string{"", "", ""},
										SubList:         [][]*Block{nil, nil, nil},
									},
								},
							},
						},
						{
							Key:             "GameEntryUpdater",
							Type:            GraphqlBlockTypeInput,
							OriginString:    "input GameEntryUpdater {\n  justDelete: Boolean\n}\n",
							SubOriginString: []string{"  justDelete: Boolean\n"},
							SubList: [][]*Block{
								{
									{
										Key:             "justDelete",
										Type:            GraphqlBlockTypeInputField,
										OriginString:    "  justDelete: Boolean\n",
										SubOriginString: []string{"", ""},
										SubList:         [][]*Block{nil, nil},
									},
								},
							},
						},
						{
							Key:             "Query",
							Type:            GraphqlBlockTypeType,
							OriginString:    "extend type Query {\n  gameEntry(id: ID!): GameEntry! @HasPermission(auth: { prefixAny: [GAME, PLAYER] })\n  gameEntries(filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!): GameEntryConnection!\n    @HasPermission(auth: { prefixAny: [GAME, PLAYER] })\n  searchGameEntry(filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!): GameEntryConnection!\n    @HasPermission(auth: { prefixAny: [GAME, PLAYER] })\n}\n",
							SubOriginString: []string{"  gameEntry(id: ID!): GameEntry! @HasPermission(auth: { prefixAny: [GAME, PLAYER] })\n  gameEntries(filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!): GameEntryConnection!\n    @HasPermission(auth: { prefixAny: [GAME, PLAYER] })\n  searchGameEntry(filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!): GameEntryConnection!\n    @HasPermission(auth: { prefixAny: [GAME, PLAYER] })\n"},
							SubList: [][]*Block{
								{
									{
										Key:             "gameEntry",
										Type:            GraphqlBlockTypeTypeField,
										OriginString:    "  gameEntry(id: ID!): GameEntry! @HasPermission(auth: { prefixAny: [GAME, PLAYER] })\n",
										SubOriginString: []string{"", "", "id: ID!"},
										SubList: [][]*Block{
											nil,
											nil,
											{
												{
													Key:             "id",
													Type:            GraphqlBlockTypeTypeFieldArg,
													OriginString:    "id: ID!",
													SubOriginString: []string{"", ""},
													SubList:         [][]*Block{nil, nil},
												},
											},
										},
									},
									{
										Key:             "gameEntries",
										Type:            GraphqlBlockTypeTypeField,
										OriginString:    "  gameEntries(filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!): GameEntryConnection!\n    @HasPermission(auth: { prefixAny: [GAME, PLAYER] })\n",
										SubOriginString: []string{"", "", "filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!"},
										SubList: [][]*Block{nil, nil,
											{
												{
													Key:             "filter",
													Type:            GraphqlBlockTypeTypeFieldArg,
													OriginString:    "filter: GameEntryFilter!",
													SubOriginString: []string{"", ""},
													SubList:         [][]*Block{nil, nil},
												},
												{
													Key:             "sorts",
													Type:            GraphqlBlockTypeTypeFieldArg,
													OriginString:    " sorts: [GameEntrySorter!]",
													SubOriginString: []string{"", ""},
													SubList:         [][]*Block{nil, nil},
												},
												{
													Key:             "offset",
													Type:            GraphqlBlockTypeTypeFieldArg,
													OriginString:    " offset: Int!",
													SubOriginString: []string{"", ""},
													SubList:         [][]*Block{nil, nil},
												},
												{
													Key:             "limit",
													Type:            GraphqlBlockTypeTypeFieldArg,
													OriginString:    " limit: Int!",
													SubOriginString: []string{"", ""},
													SubList:         [][]*Block{nil, nil},
												},
											},
										},
									},
									{
										Key:             "searchGameEntry",
										Type:            GraphqlBlockTypeTypeField,
										OriginString:    "  searchGameEntry(filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!): GameEntryConnection!\n    @HasPermission(auth: { prefixAny: [GAME, PLAYER] })\n",
										SubOriginString: []string{"", "", "filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!"},
										SubList: [][]*Block{nil, nil,
											{
												{
													Key:             "filter",
													Type:            GraphqlBlockTypeTypeFieldArg,
													OriginString:    "filter: GameEntryFilter!",
													SubOriginString: []string{"", ""},
													SubList:         [][]*Block{nil, nil},
												},
												{
													Key:             "sorts",
													Type:            GraphqlBlockTypeTypeFieldArg,
													OriginString:    " sorts: [GameEntrySorter!]",
													SubOriginString: []string{"", ""},
													SubList:         [][]*Block{nil, nil},
												},
												{
													Key:             "offset",
													Type:            GraphqlBlockTypeTypeFieldArg,
													OriginString:    " offset: Int!",
													SubOriginString: []string{"", ""},
													SubList:         [][]*Block{nil, nil},
												},
												{
													Key:             "limit",
													Type:            GraphqlBlockTypeTypeFieldArg,
													OriginString:    " limit: Int!",
													SubOriginString: []string{"", ""},
													SubList:         [][]*Block{nil, nil},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := c.Parse(tt.args.content)
			fillWantParent(nil, tt.want)
			fixTestingBlockIDAndRegIndexAndBlockParser(tt.want, got)
			if !assert.Equal(t, tt.want, got) {
				js, err := json.MarshalIndent(got, "", "\t")
				if err != nil {
					t.Errorf("json.MarshalIndent error: %v", err)
				}
				fmt.Println("Got:\n" + string(js))
			}
		})
	}
}

func TestGraphqlBlock_Merge(t *testing.T) {
	c := NewGraphqlCodeBlockParser()
	type args struct {
		income *Block
	}

	tests := []struct {
		name string
		b    *Block
		args args
		want *Block
	}{
		{
			name: "merge",
			b:    c.Parse(graphqlContent1),
			args: args{
				income: c.Parse(`
type GameEntryNew {
  test: Int!
}

type GameEntry {
  test: Int!
  test2: Int!
  Game: Game @goField(forceResolver: false)
}

extend type Query {
  gameEntries(
    filter: GameEntryFilter!
    sorts: [GameEntrySorter!]
    skip: Int!
    limit: Int!
    test: Int!
  ): GameEntryConnection!
    @HasPermission(auth: { prefixAny: [GAME, PLAYER] })
  searchGameEntry(filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!): GameEntryConnection!
    @HasPermission(auth: { prefixAny: [GAME, PLAYER] })
  gameEntry(filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!): GameEntry! @HasPermission(auth: { prefixAny: [GAME, PLAYER] })
}
`),
			},
			want: c.Parse(`
type GameEntry implements Node {
  id: ID!
  name: String!
  channelID: String!
  createAt: Timestamp!
  updateAt: Timestamp!
  """
  可见性 0: 只有本人可见 1: 公开 2: 只有所属的队伍可见。
  经过与讨论，在 dev channel 中，可见性有存在的必要；但是在 prod channel 中，可见性其实表现为是否上架，如果所有人可见就是已上架，如果只有本人可见就是未上架。
  所以这个可见性在 workshop 游戏中是其本身含义，在线上游戏中表示为这个游戏是否上架。
  """
  visibilityType: Int!
  ownerID: String!
  groupID: String!
  oldGameID: String!
  oldGameVersion: String!
  detailID: String!
  indexWeight: Int!

  detail: GameDetail @goField(forceResolver: true)
  game: Game @goField(forceResolver: true)
  test: Int!
  test2: Int!
}

input GameEntryUpdater {
  justDelete: Boolean
}

extend type Query {
  gameEntry(id: ID!): GameEntry! @HasPermission(auth: { prefixAny: [GAME, PLAYER] })
  gameEntries(filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!): GameEntryConnection!
    @HasPermission(auth: { prefixAny: [GAME, PLAYER] })
  searchGameEntry(filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!): GameEntryConnection!
    @HasPermission(auth: { prefixAny: [GAME, PLAYER] })
}

type GameEntryNew {
  test: Int!
}
`),
		},
		{
			name: "merge",
			b: c.Parse(`
extend type Query {
  gameDetailByID(id: ID!): GameDetail! @HasPermission(auth: { any: [GAME_DETAIL] })
  """
  游戏详情
  """
  gameDetail(gameID: String!, version: String!): DisplayGameAuditInfo!
    @HasPermission(auth: { prefixAny: [GAME, PLAYER, ASSET, OFFICIAL] })
    @sunset
  gameDetails(filter: GameDetailFilter!, sorts: [GameDetailSorter!], offset: Int!, limit: Int!): GameDetailConnection!
    @HasPermission(auth: { any: [GAME_DETAIL] })
}

input GameDetailUpdater {
	colorPanelKindIn: [Int!] @logConstraint(format: "max=100")
	colorPanelKindNin: [Int!] @logConstraint(format: "max=100")
	colorPanelKindGt: Int
}
`),
			args: args{
				income: c.Parse(`
extend type Query {
  """
  获取游戏详情页列表
  """
  gameDetails(filter: GameDetailFilter!, sorts: [GameDetailSorter!], offset: Int!, limit: Int!): GameDetailConnection!
  gameDetail(id: ID!): GameDetail!
  newField: String!
}

input GameDetailUpdater {
	colorModifyMetaNin: [String!] @logConstraint(format: "max=100")
}
`),
			},
			want: c.Parse(`
extend type Query {
  gameDetailByID(id: ID!): GameDetail! @HasPermission(auth: { any: [GAME_DETAIL] })
  """
  游戏详情
  """
  gameDetail(gameID: String!, version: String!): DisplayGameAuditInfo!
    @HasPermission(auth: { prefixAny: [GAME, PLAYER, ASSET, OFFICIAL] })
    @sunset
  """
  获取游戏详情页列表
  """
  gameDetails(filter: GameDetailFilter!, sorts: [GameDetailSorter!], offset: Int!, limit: Int!): GameDetailConnection!
    @HasPermission(auth: { any: [GAME_DETAIL] })
  newField: String!
}

input GameDetailUpdater {
	colorPanelKindIn: [Int!] @logConstraint(format: "max=100")
	colorPanelKindNin: [Int!] @logConstraint(format: "max=100")
	colorPanelKindGt: Int
	colorModifyMetaNin: [String!] @logConstraint(format: "max=100")
}
`),
		},
		{
			name: "merge",
			b: c.Parse(`
extend type Query {
  """
  该分类下的服装可修改颜色的元数据, 在 ColorModifiable 为 true 时，不可为空
  """
  colorModifyMeta: String!
}
`),
			args: args{
				income: c.Parse(`
extend type Query {
  """
  该分类下的服装可修改颜色的元数据, 在 ColorModifiable 为 true 时，不可为空
  """
  colorModifyMeta: String!
}
`),
			},
			want: c.Parse(`
extend type Query {
  """
  该分类下的服装可修改颜色的元数据, 在 ColorModifiable 为 true 时，不可为空
  """
  colorModifyMeta: String!
}
`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldB := tt.b.Clone()
			got := tt.b.Merge(0, tt.args.income)
			fixTestingBlockIDAndRegIndexAndBlockParser(tt.want, got)
			if !assert.Equal(t, tt.want, got) {
				fmt.Println("getMergeTestJsonString:\n```" + getMergeTestJsonString(oldB, tt.args.income, got, tt.want) + "```")
			}
		})
	}
}
