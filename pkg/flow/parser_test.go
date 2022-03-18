package flow

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseData(t *testing.T) {
	type args struct {
		content []byte
	}
	tests := []struct {
		name    string
		args    args
		want    NodesData
		wantErr bool
	}{
		{
			name: "main test",
			args: args{
				content: []byte(
					`
					{
						"3": {
							"class": "",
							"data": {
								"auth": "jira1",
								"method": "POST",
								"request": "http://localhost:8282/rest/api/2/issue"
							},
							"html": "\n  <div>\n    <div class=\"title-box\">Request</div>\n    <div class=\"box\">\n      <p>Enter request url</p>\n      <input type=\"url\" placeholder=\"https://createmyissue.com\" name=\"url\" data-action=\"focus\" df-request>\n      <p>Enter method</p>\n      <input type=\"text\" placeholder=\"POST\" name=\"method\" data-action=\"focus\" df-method>\n      <p>Enter additional headers</p>\n      <textarea data-action=\"focus\" df-headers placeholder=\"json/yaml key:value\"></textarea>\n      <p>Enter auth</p>\n      <input type=\"text\" placeholder=\"myauth\" name=\"auth\" data-action=\"focus\" df-auth>\n    </div\n  </div>\n  ",
							"id": 3,
							"inputs": {
								"input_1": {
									"connections": [
										{
											"input": "output_1",
											"node": "2"
										}
									]
								}
							},
							"name": "request",
							"outputs": {
								"output_1": {
									"connections": [
										{
											"node": "6",
											"output": "input_1"
										}
									]
								}
							},
							"pos_x": 646,
							"pos_y": -144.5,
							"typenode": false
						}
					}`,
				),
			},
			want: map[string]NodeData{
				"3": {
					Name: "request",
					Data: map[string]interface{}{
						"auth":    "jira1",
						"method":  "POST",
						"request": "http://localhost:8282/rest/api/2/issue",
					},
					Inputs: NodeConnection{
						"input_1": Connections{
							Connections: []Connection{
								{
									Node: "2",
								},
							},
						},
					},
					Outputs: NodeConnection{
						"output_1": Connections{
							Connections: []Connection{
								{
									Node:   "6",
									Output: "input_1",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseData(tt.args.content)
			fmt.Printf("%+v\n", got)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseData() = %v, want %v", got, tt.want)
			}
		})
	}
}
