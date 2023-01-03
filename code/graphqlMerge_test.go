package code

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
				SubOriginString: graphqlContent1,
				Type:            BlockType{MergeAble: true, SubsSeparator: "\n"},
				SubList: []*Block{
					{
						Key:             "GameEntry",
						Type:            GraphqlBlockTypeType,
						OriginString:    "type GameEntry implements Node {\n  id: ID!\n  name: String!\n  channelID: String!\n  createAt: Timestamp!\n  updateAt: Timestamp!\n  \"\"\"\n  可见性 0: 只有本人可见 1: 公开 2: 只有所属的队伍可见。\n  经过与讨论，在 dev channel 中，可见性有存在的必要；但是在 prod channel 中，可见性其实表现为是否上架，如果所有人可见就是已上架，如果只有本人可见就是未上架。\n  所以这个可见性在 workshop 游戏中是其本身含义，在线上游戏中表示为这个游戏是否上架。\n  \"\"\"\n  visibilityType: Int!\n  ownerID: String!\n  groupID: String!\n  oldGameID: String!\n  oldGameVersion: String!\n  detailID: String!\n  indexWeight: Int!\n\n  detail: GameDetail @goField(forceResolver: true)\n  game: Game @goField(forceResolver: true)\n}",
						SubOriginString: "id: ID!\n  name: String!\n  channelID: String!\n  createAt: Timestamp!\n  updateAt: Timestamp!\n  \"\"\"\n  可见性 0: 只有本人可见 1: 公开 2: 只有所属的队伍可见。\n  经过与讨论，在 dev channel 中，可见性有存在的必要；但是在 prod channel 中，可见性其实表现为是否上架，如果所有人可见就是已上架，如果只有本人可见就是未上架。\n  所以这个可见性在 workshop 游戏中是其本身含义，在线上游戏中表示为这个游戏是否上架。\n  \"\"\"\n  visibilityType: Int!\n  ownerID: String!\n  groupID: String!\n  oldGameID: String!\n  oldGameVersion: String!\n  detailID: String!\n  indexWeight: Int!\n\n  detail: GameDetail @goField(forceResolver: true)\n  game: Game @goField(forceResolver: true)",
						SubList: []*Block{
							{
								Key:             "id",
								Type:            GraphqlBlockTypeTypeField,
								OriginString:    "id: ID!",
								SubOriginString: "",
							},
							{
								Key:             "name",
								Type:            GraphqlBlockTypeTypeField,
								OriginString:    "name: String!",
								SubOriginString: "",
							},
							{
								Key:             "channelID",
								Type:            GraphqlBlockTypeTypeField,
								OriginString:    "channelID: String!",
								SubOriginString: "",
							},
							{
								Key:             "createAt",
								Type:            GraphqlBlockTypeTypeField,
								OriginString:    "createAt: Timestamp!",
								SubOriginString: "",
							},
							{
								Key:             "updateAt",
								Type:            GraphqlBlockTypeTypeField,
								OriginString:    "updateAt: Timestamp!",
								SubOriginString: "",
							},
							{
								Key:             "可见性 0: 只有本人可见 1: 公开 2: 只有所属的队伍可见。\n",
								Type:            GraphqlBlockExplain,
								OriginString:    "  \"\"\"\n  可见性 0: 只有本人可见 1: 公开 2: 只有所属的队伍可见。\n  经过与讨论，在 dev channel 中，可见性有存在的必要；但是在 prod channel 中，可见性其实表现为是否上架，如果所有人可见就是已上架，如果只有本人可见就是未上架。\n  所以这个可见性在 workshop 游戏中是其本身含义，在线上游戏中表示为这个游戏是否上架。\n  \"\"\"\n",
								SubOriginString: "",
								SubList:         nil,
							},
							{
								Key:             "visibilityType",
								Type:            GraphqlBlockTypeTypeField,
								OriginString:    "visibilityType: Int!",
								SubOriginString: "",
							},
							{
								Key:             "ownerID",
								Type:            GraphqlBlockTypeTypeField,
								OriginString:    "ownerID: String!",
								SubOriginString: "",
							},
							{
								Key:             "groupID",
								Type:            GraphqlBlockTypeTypeField,
								OriginString:    "groupID: String!",
								SubOriginString: "",
							},
							{
								Key:             "oldGameID",
								Type:            GraphqlBlockTypeTypeField,
								OriginString:    "oldGameID: String!",
								SubOriginString: "",
							},
							{
								Key:             "oldGameVersion",
								Type:            GraphqlBlockTypeTypeField,
								OriginString:    "oldGameVersion: String!",
								SubOriginString: "",
							},
							{
								Key:             "detailID",
								Type:            GraphqlBlockTypeTypeField,
								OriginString:    "detailID: String!",
								SubOriginString: "",
							},
							{
								Key:             "indexWeight",
								Type:            GraphqlBlockTypeTypeField,
								OriginString:    "indexWeight: Int!",
								SubOriginString: "",
							},
							{
								Key:             "detail",
								Type:            GraphqlBlockTypeTypeField,
								OriginString:    "detail: GameDetail @goField(forceResolver: true)",
								SubOriginString: "",
							},
							{
								Key:             "game",
								Type:            GraphqlBlockTypeTypeField,
								OriginString:    "game: Game @goField(forceResolver: true)",
								SubOriginString: "",
							},
						},
					},
					{
						Key:             "GameEntryUpdater",
						Type:            GraphqlBlockTypeInput,
						OriginString:    "input GameEntryUpdater {\n  justDelete: Boolean\n}\n",
						SubOriginString: "justDelete: Boolean",
						SubList: []*Block{
							{
								Key:             "justDelete",
								Type:            GraphqlBlockTypeInputField,
								OriginString:    "justDelete: Boolean",
								SubOriginString: "",
							},
						},
					},
					{
						Key:             "Query",
						Type:            GraphqlBlockTypeType,
						OriginString:    "extend type Query {\n  gameEntry(id: ID!): GameEntry! @HasPermission(auth: { prefixAny: [GAME, PLAYER] })\n  gameEntries(filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!): GameEntryConnection!\n    @HasPermission(auth: { prefixAny: [GAME, PLAYER] })\n  searchGameEntry(filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!): GameEntryConnection!\n    @HasPermission(auth: { prefixAny: [GAME, PLAYER] })\n}",
						SubOriginString: "gameEntry(id: ID!): GameEntry! @HasPermission(auth: { prefixAny: [GAME, PLAYER] })\n  gameEntries(filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!): GameEntryConnection!\n    @HasPermission(auth: { prefixAny: [GAME, PLAYER] })\n  searchGameEntry(filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!): GameEntryConnection!\n    @HasPermission(auth: { prefixAny: [GAME, PLAYER] })",
						SubList: []*Block{
							{
								Key:             "gameEntry",
								Type:            GraphqlBlockTypeTypeField,
								OriginString:    "gameEntry(id: ID!): GameEntry! @HasPermission(auth: { prefixAny: [GAME, PLAYER] })",
								SubOriginString: "id: ID!",
								SubList: []*Block{
									{
										Key:             "id",
										Type:            GraphqlBlockTypeTypeFieldArg,
										OriginString:    "id: ID!",
										SubOriginString: "",
									},
								},
							},
							{
								Key:             "gameEntries",
								Type:            GraphqlBlockTypeTypeField,
								OriginString:    "gameEntries(filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!): GameEntryConnection!\n    @HasPermission(auth: { prefixAny: [GAME, PLAYER] })",
								SubOriginString: "filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!",
								SubList: []*Block{
									{
										Key:             "filter",
										Type:            GraphqlBlockTypeTypeFieldArg,
										OriginString:    "filter: GameEntryFilter!",
										SubOriginString: "",
									},
									{
										Key:             "sorts",
										Type:            GraphqlBlockTypeTypeFieldArg,
										OriginString:    "sorts: [GameEntrySorter!]",
										SubOriginString: "",
									},
									{
										Key:             "offset",
										Type:            GraphqlBlockTypeTypeFieldArg,
										OriginString:    "offset: Int!",
										SubOriginString: "",
									},
									{
										Key:             "limit",
										Type:            GraphqlBlockTypeTypeFieldArg,
										OriginString:    "limit: Int!",
										SubOriginString: "",
									},
								},
							},
							{
								Key:             "searchGameEntry",
								Type:            GraphqlBlockTypeTypeField,
								OriginString:    "searchGameEntry(filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!): GameEntryConnection!\n    @HasPermission(auth: { prefixAny: [GAME, PLAYER] })",
								SubOriginString: "filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!",
								SubList: []*Block{
									{
										Key:             "filter",
										Type:            GraphqlBlockTypeTypeFieldArg,
										OriginString:    "filter: GameEntryFilter!",
										SubOriginString: "",
									},
									{
										Key:             "sorts",
										Type:            GraphqlBlockTypeTypeFieldArg,
										OriginString:    "sorts: [GameEntrySorter!]",
										SubOriginString: "",
									},
									{
										Key:             "offset",
										Type:            GraphqlBlockTypeTypeFieldArg,
										OriginString:    "offset: Int!",
										SubOriginString: "",
									},
									{
										Key:             "limit",
										Type:            GraphqlBlockTypeTypeFieldArg,
										OriginString:    "limit: Int!",
										SubOriginString: "",
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
	incomeStr := `
type GameEntryNew {
  test: Int!
}

type GameEntry {
  test: Int!
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
}`
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
				income: c.Parse(incomeStr),
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
}`),
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
`),
			args: args{
				income: c.Parse(`
extend type Query {
gameDetails(filter: GameDetailFilter!, sorts: [GameDetailSorter!], offset: Int!, limit: Int!): GameDetailConnection!
gameDetail(id: ID!): GameDetail!
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
  gameDetails(filter: GameDetailFilter!, sorts: [GameDetailSorter!], offset: Int!, limit: Int!): GameDetailConnection!
    @HasPermission(auth: { any: [GAME_DETAIL] })
}
`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsOld, err := json.MarshalIndent(tt.b, "", "\t")
			if err != nil {
				t.Errorf("json.MarshalIndent error: %v", err)
			}
			incomeJs, err := json.MarshalIndent(tt.args.income, "", "\t")
			if err != nil {
				t.Errorf("json.MarshalIndent error: %v", err)
			}
			got := tt.b.Merge(tt.args.income)

			if !assert.Equal(t, tt.want, got) {
				fmt.Println("Old:\n```" + string(jsOld) + "```")
				fmt.Println("Income:\n```" + string(incomeJs) + "```")
				js, err := json.MarshalIndent(got, "", "\t")
				if err != nil {
					t.Errorf("json.MarshalIndent error: %v", err)
				}
				fmt.Println("Got:\n```" + string(js) + "```")
				fmt.Println("Got origin:\n```" + got.OriginString + "```")
			}
		})
	}
}