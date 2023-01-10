package code

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var protoContent1 string = `

message UpdateGameDetailByGameIDRequest {
	string game_id = 1;
	GameDetailUpdater updater = 2;
	optional string op_user_id = 3;
}
message UpdateGameDetailByGameIDResponse {
}

message BatchGetAssetTokenConfigResponse {
	repeated AssetTokenConfig info_list = 1 [json_name = "nodes"];
}


service GamedevService {
	option (base.soptions) = {ds_rpc: true, lua_export: true};
	rpc FullSetGameHotfixData(FullSetGameHotfixDataRequest) returns (FullSetGameHotfixDataResponse);
	rpc GetGameHotfixData(GetGameHotfixDataRequest) returns (GetGameHotfixDataResponse) {
		option (base.moptions) = {ds_rpc: true, lua_export: true};
	};
	rpc GetGameHotfixDataView(GetGameHotfixDataViewRequest) returns (GetGameHotfixDataViewResponse);

	// game channel
	rpc ListChannelGameEntry(ListChannelGameEntryRequest) returns(ListChannelGameEntryResponse);
	rpc InsertGameProject(InsertGameProjectRequest) returns (InsertGameProjectResponse);
}
`

func TestProtoBlockFromString(t *testing.T) {
	c := NewProtoBufCodeBlockParser()
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
				content: protoContent1,
			},
			want: &Block{
				OriginString:    protoContent1,
				SubOriginString: []string{protoContent1},
				Type:            BlockType{SubMergeType: []*MergeConfig{{Append: true, ReplaceBlockType: nil}}, SubsSeparator: "\n"},
				SubList: [][]*Block{
					{
						{
							Key:             "UpdateGameDetailByGameIDRequest",
							Type:            ProtoBlockTypeMessage,
							OriginString:    "message UpdateGameDetailByGameIDRequest {\n\tstring game_id = 1;\n\tGameDetailUpdater updater = 2;\n\toptional string op_user_id = 3;\n}\n",
							SubOriginString: []string{"	string game_id = 1;\n\tGameDetailUpdater updater = 2;\n\toptional string op_user_id = 3;\n"},
							SubList: [][]*Block{
								{
									{
										Key:             "game_id",
										Type:            ProtoBlockTypeMessageField,
										OriginString:    "	string game_id = 1;\n",
										SubOriginString: []string{""},
										SubList:         [][]*Block{nil},
									},
									{
										Key:             "updater",
										Type:            ProtoBlockTypeMessageField,
										OriginString:    "	GameDetailUpdater updater = 2;\n",
										SubOriginString: []string{""},
										SubList:         [][]*Block{nil},
									},
									{
										Key:             "op_user_id",
										Type:            ProtoBlockTypeMessageField,
										OriginString:    "	optional string op_user_id = 3;\n",
										SubOriginString: []string{""},
										SubList:         [][]*Block{nil},
									},
								},
							},
						},
						{
							Key:             "UpdateGameDetailByGameIDResponse",
							Type:            ProtoBlockTypeMessage,
							OriginString:    "message UpdateGameDetailByGameIDResponse {\n}\n",
							SubOriginString: []string{""},
							SubList:         [][]*Block{nil},
						},
						{
							Key:             "BatchGetAssetTokenConfigResponse",
							Type:            ProtoBlockTypeMessage,
							OriginString:    "message BatchGetAssetTokenConfigResponse {\n\trepeated AssetTokenConfig info_list = 1 [json_name = \"nodes\"];\n}\n",
							SubOriginString: []string{"	repeated AssetTokenConfig info_list = 1 [json_name = \"nodes\"];\n"},
							SubList: [][]*Block{
								{
									{
										Key:             "info_list",
										Type:            ProtoBlockTypeMessageField,
										OriginString:    "	repeated AssetTokenConfig info_list = 1 [json_name = \"nodes\"];\n",
										SubOriginString: []string{"json_name = \"nodes\""},
										SubList: [][]*Block{
											{
												{
													Key:             "json_name",
													Type:            ProtoBlockTypeOptionItem,
													OriginString:    "json_name = \"nodes\"",
													SubOriginString: []string{},
													SubList:         [][]*Block{},
												},
											},
										},
									},
								},
							},
						},
						{
							Key:             "GamedevService",
							Type:            ProtoBlockTypeService,
							OriginString:    "service GamedevService {\n\toption (base.soptions) = {ds_rpc: true, lua_export: true};\n\trpc FullSetGameHotfixData(FullSetGameHotfixDataRequest) returns (FullSetGameHotfixDataResponse);\n\trpc GetGameHotfixData(GetGameHotfixDataRequest) returns (GetGameHotfixDataResponse) {\n\t\toption (base.moptions) = {ds_rpc: true, lua_export: true};\n\t};\n\trpc GetGameHotfixDataView(GetGameHotfixDataViewRequest) returns (GetGameHotfixDataViewResponse);\n\n\t// game channel\n\trpc ListChannelGameEntry(ListChannelGameEntryRequest) returns(ListChannelGameEntryResponse);\n\trpc InsertGameProject(InsertGameProjectRequest) returns (InsertGameProjectResponse);\n}\n",
							SubOriginString: []string{"	option (base.soptions) = {ds_rpc: true, lua_export: true};\n\trpc FullSetGameHotfixData(FullSetGameHotfixDataRequest) returns (FullSetGameHotfixDataResponse);\n\trpc GetGameHotfixData(GetGameHotfixDataRequest) returns (GetGameHotfixDataResponse) {\n\t\toption (base.moptions) = {ds_rpc: true, lua_export: true};\n\t};\n\trpc GetGameHotfixDataView(GetGameHotfixDataViewRequest) returns (GetGameHotfixDataViewResponse);\n\n\t// game channel\n\trpc ListChannelGameEntry(ListChannelGameEntryRequest) returns(ListChannelGameEntryResponse);\n\trpc InsertGameProject(InsertGameProjectRequest) returns (InsertGameProjectResponse);\n"},
							SubList: [][]*Block{
								{
									{
										Key:             "(base.soptions)",
										Type:            ProtoBlockTypeOption,
										OriginString:    "	option (base.soptions) = {ds_rpc: true, lua_export: true};\n",
										SubOriginString: []string{"ds_rpc: true, lua_export: true"},
										SubList: [][]*Block{
											{
												{
													Key:             "ds_rpc",
													Type:            ProtoBlockTypeOptionItem,
													OriginString:    "ds_rpc: true",
													SubOriginString: []string{},
													SubList:         [][]*Block{},
												},
												{
													Key:             "lua_export",
													Type:            ProtoBlockTypeOptionItem,
													OriginString:    " lua_export: true",
													SubOriginString: []string{},
													SubList:         [][]*Block{},
												},
											},
										},
									},
									{
										Key:             "FullSetGameHotfixData",
										Type:            ProtoBlockTypeRPC,
										OriginString:    "	rpc FullSetGameHotfixData(FullSetGameHotfixDataRequest) returns (FullSetGameHotfixDataResponse);\n",
										SubOriginString: []string{""},
										SubList:         [][]*Block{nil},
									},
									{
										Key:             "GetGameHotfixData",
										Type:            ProtoBlockTypeRPC,
										OriginString:    "	rpc GetGameHotfixData(GetGameHotfixDataRequest) returns (GetGameHotfixDataResponse) {\n\t\toption (base.moptions) = {ds_rpc: true, lua_export: true};\n\t};\n",
										SubOriginString: []string{"		option (base.moptions) = {ds_rpc: true, lua_export: true};\n"},
										SubList: [][]*Block{
											{
												{
													Key:             "(base.moptions)",
													Type:            ProtoBlockTypeOption,
													OriginString:    "		option (base.moptions) = {ds_rpc: true, lua_export: true};\n",
													SubOriginString: []string{"ds_rpc: true, lua_export: true"},
													SubList: [][]*Block{
														{
															{
																Key:             "ds_rpc",
																Type:            ProtoBlockTypeOptionItem,
																OriginString:    "ds_rpc: true",
																SubOriginString: []string{},
																SubList:         [][]*Block{},
															},
															{
																Key:             "lua_export",
																Type:            ProtoBlockTypeOptionItem,
																OriginString:    " lua_export: true",
																SubOriginString: []string{},
																SubList:         [][]*Block{},
															},
														},
													},
												},
											},
										},
									},
									{
										Key:             "GetGameHotfixDataView",
										Type:            ProtoBlockTypeRPC,
										OriginString:    "	rpc GetGameHotfixDataView(GetGameHotfixDataViewRequest) returns (GetGameHotfixDataViewResponse);\n",
										SubOriginString: []string{""},
										SubList:         [][]*Block{nil},
									},
									{
										Key:             "ListChannelGameEntry",
										Type:            ProtoBlockTypeRPC,
										OriginString:    "	rpc ListChannelGameEntry(ListChannelGameEntryRequest) returns(ListChannelGameEntryResponse);\n",
										SubOriginString: []string{""},
										SubList:         [][]*Block{nil},
									},
									{
										Key:             "InsertGameProject",
										Type:            ProtoBlockTypeRPC,
										OriginString:    "	rpc InsertGameProject(InsertGameProjectRequest) returns (InsertGameProjectResponse);\n",
										SubOriginString: []string{""},
										SubList:         [][]*Block{nil},
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

func fillWantParent(p *Block, want *Block) {
	want.Parent = p
	// fill want parent
	for _, v := range want.SubList {
		for _, v2 := range v {
			if v2 != nil {
				fillWantParent(want, v2)
			}
		}
	}
}

func TestProtoBlock_Merge(t *testing.T) {
	c := NewProtoBufCodeBlockParser()
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
			b:    c.Parse(protoContent1),
			args: args{
				income: c.Parse(`
message UpdateGameDetailByGameIDRequest {
	string game_id = 1 [ds_rpc: true, lua_export: true];
}
message UpdateGameDetailByGameIDResponse {
	repeated AssetTokenConfig asset_token_config = 1 [ds_rpc: true, lua_export: true];
}
message BatchGetAssetTokenConfigResponse {
	repeated AssetTokenConfig asset_token_config = 1;
	repeated AssetTokenConfig info_list = 1 [ds_rpc: true, lua_export: true];
}
`),
			},
			want: c.Parse(`

message UpdateGameDetailByGameIDRequest {
	string game_id = 1[ds_rpc: true, lua_export: true];
	GameDetailUpdater updater = 2;
	optional string op_user_id = 3;
}
message UpdateGameDetailByGameIDResponse {
	repeated AssetTokenConfig asset_token_config = 1 [ds_rpc: true, lua_export: true];
}

message BatchGetAssetTokenConfigResponse {
	repeated AssetTokenConfig info_list = 1 [json_name = "nodes", ds_rpc: true, lua_export: true];
	repeated AssetTokenConfig asset_token_config = 1;
}


service GamedevService {
	option (base.soptions) = {ds_rpc: true, lua_export: true};
	rpc FullSetGameHotfixData(FullSetGameHotfixDataRequest) returns (FullSetGameHotfixDataResponse);
	rpc GetGameHotfixData(GetGameHotfixDataRequest) returns (GetGameHotfixDataResponse) {
		option (base.moptions) = {ds_rpc: true, lua_export: true};
	};
	rpc GetGameHotfixDataView(GetGameHotfixDataViewRequest) returns (GetGameHotfixDataViewResponse);

	// game channel
	rpc ListChannelGameEntry(ListChannelGameEntryRequest) returns(ListChannelGameEntryResponse);
	rpc InsertGameProject(InsertGameProjectRequest) returns (InsertGameProjectResponse);
}
`),
		},
		{
			name: "merge2",
			b: c.Parse(`
message UpdateGameDetailByGameIDResponse {
}
`),
			args: args{
				income: c.Parse(`
message UpdateGameDetailByGameIDResponse {
	repeated AssetTokenConfig asset_token_config = 1 [ds_rpc: true, lua_export: true];
}
`),
			},
			want: c.Parse(`
message UpdateGameDetailByGameIDResponse {
	repeated AssetTokenConfig asset_token_config = 1 [ds_rpc: true, lua_export: true];
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
				fmt.Println("getMergeTestJsonString:\n```" + getMergeTestJsonString(oldB, tt.args.income, got, tt.want) + "``` gotOrigin: ```" + got.OriginString + "```")
			}
		})
	}
}
