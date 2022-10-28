package writer

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
	type args struct {
		p       *ProtoBlock
		content string
	}
	tests := []struct {
		name string
		args args
		want *ProtoBlock
	}{
		{
			name: "test1",
			args: args{
				content: protoContent1,
			},
			want: &ProtoBlock{
				OriginString:    protoContent1,
				SubOriginString: protoContent1,
				SubList: []*ProtoBlock{
					{
						Key:             "UpdateGameDetailByGameIDRequest",
						Type:            "message",
						OriginString:    "message UpdateGameDetailByGameIDRequest {\n\tstring game_id = 1;\n\tGameDetailUpdater updater = 2;\n\toptional string op_user_id = 3;\n  }\n",
						SubOriginString: "\n\tstring game_id = 1;\n\tGameDetailUpdater updater = 2;\n\toptional string op_user_id = 3;\n  ",
						SubList: []*ProtoBlock{
							{
								Key:          "game_id",
								Type:         "message field",
								OriginString: "\tstring game_id = 1;\n",
								SubList:      nil,
							},
							{
								Key:          "updater",
								Type:         "message field",
								OriginString: "\tGameDetailUpdater updater = 2;\n",
								SubList:      nil,
							},
							{
								Key:          "op_user_id",
								Type:         "message field",
								OriginString: "\toptional string op_user_id = 3;\n",
								SubList:      nil,
							},
						},
					},
					{
						Key:             "UpdateGameDetailByGameIDResponse",
						Type:            "message",
						OriginString:    "  message UpdateGameDetailByGameIDResponse {\n  }\n",
						SubOriginString: "\n  ",
						SubList:         nil,
					},
					{
						Key:             "BatchGetAssetTokenConfigResponse",
						Type:            "message",
						OriginString:    "  message BatchGetAssetTokenConfigResponse {\n\trepeated AssetTokenConfig info_list = 1 [json_name = \"nodes\"];\n  }\n",
						SubOriginString: "\n\trepeated AssetTokenConfig info_list = 1 [json_name = \"nodes\"];\n  ",
						SubList: []*ProtoBlock{
							{
								Key:             "info_list",
								Type:            "message field",
								OriginString:    "\trepeated AssetTokenConfig info_list = 1 [json_name = \"nodes\"];\n",
								SubOriginString: "json_name = \"nodes\"",
								SubList: []*ProtoBlock{
									{
										Key:          "json_name",
										Type:         "option item",
										OriginString: "json_name = \"nodes\"",
										SubList:      nil,
									},
								},
							},
						},
					},
					{
						Key:             "GamedevService",
						Type:            "service",
						OriginString:    "service GamedevService {\n\toption (base.soptions) = {ds_rpc: true, lua_export: true};\n\trpc FullSetGameHotfixData(FullSetGameHotfixDataRequest) returns (FullSetGameHotfixDataResponse);\n\trpc GetGameHotfixData(GetGameHotfixDataRequest) returns (GetGameHotfixDataResponse) {\n\t  option (base.moptions) = {ds_rpc: true, lua_export: true};\n\t};\n\trpc GetGameHotfixDataView(GetGameHotfixDataViewRequest) returns (GetGameHotfixDataViewResponse);\n  \n\t// game channel\n\trpc ListChannelGameEntry(ListChannelGameEntryRequest) returns(ListChannelGameEntryResponse);\n\trpc InsertGameProject(InsertGameProjectRequest) returns (InsertGameProjectResponse);\n  }\n",
						SubOriginString: "\n\toption (base.soptions) = {ds_rpc: true, lua_export: true};\n\trpc FullSetGameHotfixData(FullSetGameHotfixDataRequest) returns (FullSetGameHotfixDataResponse);\n\trpc GetGameHotfixData(GetGameHotfixDataRequest) returns (GetGameHotfixDataResponse) {\n\t  option (base.moptions) = {ds_rpc: true, lua_export: true};\n\t};\n\trpc GetGameHotfixDataView(GetGameHotfixDataViewRequest) returns (GetGameHotfixDataViewResponse);\n  \n\t// game channel\n\trpc ListChannelGameEntry(ListChannelGameEntryRequest) returns(ListChannelGameEntryResponse);\n\trpc InsertGameProject(InsertGameProjectRequest) returns (InsertGameProjectResponse);\n  ",
						SubList: []*ProtoBlock{
							{
								Key:             "(base.soptions)",
								Type:            "options",
								OriginString:    "\toption (base.soptions) = {ds_rpc: true, lua_export: true};\n",
								SubOriginString: "ds_rpc: true, lua_export: true",
								SubList: []*ProtoBlock{
									{
										Key:          "ds_rpc",
										Type:         "option item",
										OriginString: "ds_rpc: true, ",
										SubList:      nil,
									},
									{
										Key:          "lua_export",
										Type:         "option item",
										OriginString: "lua_export: true",
										SubList:      nil,
									},
								},
							},
							{
								Key:          "FullSetGameHotfixData",
								Type:         "rpc",
								OriginString: "\trpc FullSetGameHotfixData(FullSetGameHotfixDataRequest) returns (FullSetGameHotfixDataResponse);\n",
								SubList:      nil,
							},
							{
								Key:             "GetGameHotfixData",
								Type:            "rpc",
								OriginString:    "\trpc GetGameHotfixData(GetGameHotfixDataRequest) returns (GetGameHotfixDataResponse) {\n\t  option (base.moptions) = {ds_rpc: true, lua_export: true};\n\t};\n",
								SubOriginString: "\n\t  option (base.moptions) = {ds_rpc: true, lua_export: true};\n\t",
								SubList: []*ProtoBlock{
									{
										Key:             "(base.moptions)",
										Type:            "options",
										OriginString:    "\t  option (base.moptions) = {ds_rpc: true, lua_export: true};\n",
										SubOriginString: "ds_rpc: true, lua_export: true",
										SubList: []*ProtoBlock{
											{
												Key:          "ds_rpc",
												Type:         "option item",
												OriginString: "ds_rpc: true, ",
												SubList:      nil,
											},
											{
												Key:          "lua_export",
												Type:         "option item",
												OriginString: "lua_export: true",
												SubList:      nil,
											},
										},
									},
								},
							},
							{
								Key:          "GetGameHotfixDataView",
								Type:         "rpc",
								OriginString: "\trpc GetGameHotfixDataView(GetGameHotfixDataViewRequest) returns (GetGameHotfixDataViewResponse);\n",
								SubList:      nil,
							},
							{
								Key:          "ListChannelGameEntry",
								Type:         "rpc",
								OriginString: "\trpc ListChannelGameEntry(ListChannelGameEntryRequest) returns(ListChannelGameEntryResponse);\n",
								SubList:      nil,
							},
							{
								Key:          "InsertGameProject",
								Type:         "rpc",
								OriginString: "\trpc InsertGameProject(InsertGameProjectRequest) returns (InsertGameProjectResponse);\n",
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
			got := ProtoBlockFromString(tt.args.content)
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

func fillWantParent(p *ProtoBlock, want *ProtoBlock) {
	want.Parent = p
	// fill want parent
	for _, v := range want.SubList {
		if v != nil {
			fillWantParent(want, v)
		}
	}
}

func TestProtoBlock_Merge(t *testing.T) {
	type args struct {
		income *ProtoBlock
	}
	incomeStr := `
				message BatchGetAssetTokenConfigResponse {
					repeated AssetTokenConfig asset_token_config = 1;
					repeated AssetTokenConfig info_list = 1 [ds_rpc: true, lua_export: true];
				}
				`
	tests := []struct {
		name string
		b    *ProtoBlock
		args args
		want *ProtoBlock
	}{
		{
			name: "merge",
			b:    ProtoBlockFromString(protoContent1),
			args: args{
				income: ProtoBlockFromString(incomeStr),
			},
			want: ProtoBlockFromString(`

message UpdateGameDetailByGameIDRequest {
	string game_id = 1;
	GameDetailUpdater updater = 2;
	optional string op_user_id = 3;
  }
  message UpdateGameDetailByGameIDResponse {
  }
  
  message BatchGetAssetTokenConfigResponse {
	repeated AssetTokenConfig info_list = 1 [json_name = "nodes"ds_rpc: true, lua_export: true];
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
			fmt.Println("Old:\n```" + string(jsOld) + "```")
			got := tt.b.Merge(tt.args.income)

			if !assert.Equal(t, tt.want, got) {
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
