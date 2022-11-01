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
		p       *CodeBlock
		content string
	}
	tests := []struct {
		name string
		args args
		want *CodeBlock
	}{
		{
			name: "test1",
			args: args{
				content: protoContent1,
			},
			want: &CodeBlock{
				OriginString:    protoContent1,
				SubOriginString: protoContent1,
				Type:            CodeBlockType{MergeAble: true},
				SubList: []*CodeBlock{
					{
						Key:             "UpdateGameDetailByGameIDRequest",
						Type:            ProtoBlockTypeMessage,
						OriginString:    "message UpdateGameDetailByGameIDRequest {\n\tstring game_id = 1;\n\tGameDetailUpdater updater = 2;\n\toptional string op_user_id = 3;\n}",
						SubOriginString: "string game_id = 1;\n\tGameDetailUpdater updater = 2;\n\toptional string op_user_id = 3;",
						SubList: []*CodeBlock{
							{
								Key:          "game_id",
								Type:         ProtoBlockTypeMessageField,
								OriginString: "string game_id = 1;",
								SubList:      nil,
							},
							{
								Key:          "updater",
								Type:         ProtoBlockTypeMessageField,
								OriginString: "GameDetailUpdater updater = 2;",
								SubList:      nil,
							},
							{
								Key:          "op_user_id",
								Type:         ProtoBlockTypeMessageField,
								OriginString: "optional string op_user_id = 3;",
								SubList:      nil,
							},
						},
					},
					{
						Key:             "UpdateGameDetailByGameIDResponse",
						Type:            ProtoBlockTypeMessage,
						OriginString:    "message UpdateGameDetailByGameIDResponse {\n}",
						SubOriginString: "",
						SubList:         nil,
					},
					{
						Key:             "BatchGetAssetTokenConfigResponse",
						Type:            ProtoBlockTypeMessage,
						OriginString:    "message BatchGetAssetTokenConfigResponse {\n\trepeated AssetTokenConfig info_list = 1 [json_name = \"nodes\"];\n}",
						SubOriginString: "repeated AssetTokenConfig info_list = 1 [json_name = \"nodes\"];",
						SubList: []*CodeBlock{
							{
								Key:             "info_list",
								Type:            ProtoBlockTypeMessageField,
								OriginString:    "repeated AssetTokenConfig info_list = 1 [json_name = \"nodes\"];",
								SubOriginString: "json_name = \"nodes\"",
								SubList: []*CodeBlock{
									{
										Key:          "json_name",
										Type:         ProtoBlockTypeOptionItem,
										OriginString: "json_name = \"nodes\"",
										SubList:      nil,
									},
								},
							},
						},
					},
					{
						Key:             "GamedevService",
						Type:            ProtoBlockTypeService,
						OriginString:    "service GamedevService {\n\toption (base.soptions) = {ds_rpc: true, lua_export: true};\n\trpc FullSetGameHotfixData(FullSetGameHotfixDataRequest) returns (FullSetGameHotfixDataResponse);\n\trpc GetGameHotfixData(GetGameHotfixDataRequest) returns (GetGameHotfixDataResponse) {\n\t\toption (base.moptions) = {ds_rpc: true, lua_export: true};\n\t};\n\trpc GetGameHotfixDataView(GetGameHotfixDataViewRequest) returns (GetGameHotfixDataViewResponse);\n\n\t// game channel\n\trpc ListChannelGameEntry(ListChannelGameEntryRequest) returns(ListChannelGameEntryResponse);\n\trpc InsertGameProject(InsertGameProjectRequest) returns (InsertGameProjectResponse);\n}",
						SubOriginString: "option (base.soptions) = {ds_rpc: true, lua_export: true};\n\trpc FullSetGameHotfixData(FullSetGameHotfixDataRequest) returns (FullSetGameHotfixDataResponse);\n\trpc GetGameHotfixData(GetGameHotfixDataRequest) returns (GetGameHotfixDataResponse) {\n\t\toption (base.moptions) = {ds_rpc: true, lua_export: true};\n\t};\n\trpc GetGameHotfixDataView(GetGameHotfixDataViewRequest) returns (GetGameHotfixDataViewResponse);\n\n\t// game channel\n\trpc ListChannelGameEntry(ListChannelGameEntryRequest) returns(ListChannelGameEntryResponse);\n\trpc InsertGameProject(InsertGameProjectRequest) returns (InsertGameProjectResponse);",
						SubList: []*CodeBlock{
							{
								Key:             "(base.soptions)",
								Type:            ProtoBlockTypeOption,
								OriginString:    "option (base.soptions) = {ds_rpc: true, lua_export: true};",
								SubOriginString: "ds_rpc: true, lua_export: true",
								SubList: []*CodeBlock{
									{
										Key:          "ds_rpc",
										Type:         ProtoBlockTypeOptionItem,
										OriginString: "ds_rpc: true",
										SubList:      nil,
									},
									{
										Key:          "lua_export",
										Type:         ProtoBlockTypeOptionItem,
										OriginString: "lua_export: true",
										SubList:      nil,
									},
								},
							},
							{
								Key:          "FullSetGameHotfixData",
								Type:         ProtoBlockTypeRPC,
								OriginString: "rpc FullSetGameHotfixData(FullSetGameHotfixDataRequest) returns (FullSetGameHotfixDataResponse);",
								SubList:      nil,
							},
							{
								Key:             "GetGameHotfixData",
								Type:            ProtoBlockTypeRPC,
								OriginString:    "rpc GetGameHotfixData(GetGameHotfixDataRequest) returns (GetGameHotfixDataResponse) {\n\t\toption (base.moptions) = {ds_rpc: true, lua_export: true};\n\t};",
								SubOriginString: "option (base.moptions) = {ds_rpc: true, lua_export: true};",
								SubList: []*CodeBlock{
									{
										Key:             "(base.moptions)",
										Type:            ProtoBlockTypeOption,
										OriginString:    "option (base.moptions) = {ds_rpc: true, lua_export: true};",
										SubOriginString: "ds_rpc: true, lua_export: true",
										SubList: []*CodeBlock{
											{
												Key:          "ds_rpc",
												Type:         ProtoBlockTypeOptionItem,
												OriginString: "ds_rpc: true",
												SubList:      nil,
											},
											{
												Key:          "lua_export",
												Type:         ProtoBlockTypeOptionItem,
												OriginString: "lua_export: true",
												SubList:      nil,
											},
										},
									},
								},
							},
							{
								Key:          "GetGameHotfixDataView",
								Type:         ProtoBlockTypeRPC,
								OriginString: "rpc GetGameHotfixDataView(GetGameHotfixDataViewRequest) returns (GetGameHotfixDataViewResponse);",
								SubList:      nil,
							},
							{
								Key:          "ListChannelGameEntry",
								Type:         ProtoBlockTypeRPC,
								OriginString: "rpc ListChannelGameEntry(ListChannelGameEntryRequest) returns(ListChannelGameEntryResponse);",
								SubList:      nil,
							},
							{
								Key:          "InsertGameProject",
								Type:         ProtoBlockTypeRPC,
								OriginString: "rpc InsertGameProject(InsertGameProjectRequest) returns (InsertGameProjectResponse);",
								SubList:      nil,
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

func fillWantParent(p *CodeBlock, want *CodeBlock) {
	want.Parent = p
	// fill want parent
	for _, v := range want.SubList {
		if v != nil {
			fillWantParent(want, v)
		}
	}
}

func TestProtoBlock_Merge(t *testing.T) {
	c := NewProtoBufCodeBlockParser()
	type args struct {
		income *CodeBlock
	}
	incomeStr := `

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
				`
	tests := []struct {
		name string
		b    *CodeBlock
		args args
		want *CodeBlock
	}{
		{
			name: "merge",
			b:    c.Parse(protoContent1),
			args: args{
				income: c.Parse(incomeStr),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsOld, err := json.MarshalIndent(tt.b, "", "\t")
			if err != nil {
				t.Errorf("json.MarshalIndent error: %v", err)
			}
			got := tt.b.Merge(tt.args.income)

			if !assert.Equal(t, tt.want, got) {
				fmt.Println("Old:\n```" + string(jsOld) + "```")
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
