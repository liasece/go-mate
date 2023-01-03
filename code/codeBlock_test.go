package code

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodeBlock_addSub(t *testing.T) {
	graphqlC := NewGraphqlCodeBlockParser()
	type args struct {
		income *Block
	}
	tests := []struct {
		name string
		args args
		b    *Block
		want *Block
	}{
		{
			name: "test1",
			args: args{
				income: graphqlC.Parse(`
extend type Query {
  gameEntries(
    filter: GameEntryFilter!
    sorts: [GameEntrySorter!]
    skip: Int!
    limit: Int!
    test: Int!
  ): GameEntryConnection!
    @HasPermission(auth: { prefixAny: [GAME, PLAYER] })
}
`).SubList[0].SubList[0].SubList[4],
			},
			b: graphqlC.Parse(`
extend type Query {
  gameEntries(filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!): GameEntryConnection!
    @HasPermission(auth: { prefixAny: [GAME, PLAYER] })
}
`).SubList[0].SubList[0],
			want: graphqlC.Parse(`
extend type Query {
  gameEntries(filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!, test: Int!): GameEntryConnection!
    @HasPermission(auth: { prefixAny: [GAME, PLAYER] })
}
`).SubList[0].SubList[0],
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
			tt.b.addSub(tt.args.income)

			if !assert.Equal(t, tt.want, tt.b) {
				fmt.Println("Old:\n```" + string(jsOld) + "```")
				fmt.Println("Income:\n```" + string(incomeJs) + "```")
				js, err := json.MarshalIndent(tt.b, "", "\t")
				if err != nil {
					t.Errorf("json.MarshalIndent error: %v", err)
				}
				fmt.Println("Got:\n```" + string(js) + "```")
				fmt.Println("Got origin:\n```" + tt.b.OriginString + "```")
			}
		})
	}
}

func TestCodeBlock_getSubJoinString(t *testing.T) {
	graphqlC := NewGraphqlCodeBlockParser()
	type args struct {
		income *Block
	}
	tests := []struct {
		name string
		args args
		b    *Block
		want string
	}{
		{
			name: "test1",
			args: args{
				income: graphqlC.Parse(`
extend type Query {
  gameEntries(
    filter: GameEntryFilter!
    sorts: [GameEntrySorter!]
    skip: Int!
    limit: Int!
    test: Int!
  ): GameEntryConnection!
    @HasPermission(auth: { prefixAny: [GAME, PLAYER] })
}
`).SubList[0].SubList[0].SubList[4],
			},
			b: graphqlC.Parse(`
extend type Query {
  gameEntries(filter: GameEntryFilter!, sorts: [GameEntrySorter!], offset: Int!, limit: Int!): GameEntryConnection!
    @HasPermission(auth: { prefixAny: [GAME, PLAYER] })
}
`).SubList[0].SubList[0],
			want: ", ",
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
			fmt.Println("Old:\n```" + string(jsOld) + "```")
			fmt.Println("Income:\n```" + string(incomeJs) + "```")
			got := tt.b.getSubJoinString()
			assert.Equal(t, tt.want, got)
		})
	}
}
