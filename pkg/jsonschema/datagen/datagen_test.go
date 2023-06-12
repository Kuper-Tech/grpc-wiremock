package datagen

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_datagen_Generate(t *testing.T) {
	tests := []struct {
		name   string
		schema string
	}{
		// {"seconds":88051}
		{schema: `{"properties":{"seconds":{"type":"integer"}},"required":["seconds"],"type":"object"}`},

		// {"accessToken":"70sACB-","expiresAt":"2021-06-11T09:13:05+03:00","refreshToken":"-A65q_*I-4-z**"}
		{schema: `{"properties":{"accessToken":{"type":"string"},"expiresAt":{"format":"date-time","type":"string"},"refreshToken":{"type":"string"}},"required":["accessToken","expiresAt","refreshToken"],"type":"object"}`},

		// {"accessToken":"lToa--*Ml","expiresAt":"2018-06-18T23:41:01+03:00"}
		{schema: `{"properties":{"accessToken":{"type":"string"},"expiresAt":{"format":"date-time","type":"string"}},"required":["accessToken","expiresAt"],"type":"object"}`},

		// {"code":"-_*-H__-B*-8-_","fingerprint":"-XfT60a--_g-_S*","phone":"**-M___-6"}
		{schema: `{"properties":{"code":{"description":"OTP (One Time Password)","type":"string"},"fingerprint":{"type":"string"},"phone":{"description":"телефонный номер","type":"string"}},"required":["code","fingerprint","phone"],"type":"object"}`},

		// {"invalid-params":[{"name":"7**-w_0t","reason":"__5_--3a--_"},{"name":"6_0__Vi43-*5_f--","reason":"o-Q_zI*_*mb6_"},{"name":"A1*R*U","reason":"4nrD-t@_"},{"name":"_kPkg--*-_*","reason":"*7-I*-_W"},{"name":"-__-_4-nVf_6*_F--63","reason":"*4__*_O--p_B_1-"},{"name":"9---ll_*t-kWdi-","reason":"w_62L-l"},{"name":"_9LH_*8131_*6*4_l-_","reason":"q**F-*X_"},{"name":"**_HR1e--p__HS_-Y","reason":"H9B_*3C-_7**N3___-"},{"name":"-@_u_M0--_","reason":"U_*H_-*-s_u-_-R-"}],"title":"-**2_t"}
		{schema: `{"properties":{"invalid-params": {"items": {"properties": {"name": {"type": "string"},"reason": {"type": "string"}},"type":"object"},"type": "array"}, "title": {"type": "string"}},"type":"object"}`},

		// {"invalid-params":[{"name":"Fg**_r","reason":"w*-xiOt4P2-"},{"name":"1*Q-_-*","reason":"*ag-_y2*k1"},{"name":"3_q--2p*AQ_9N_*-","reason":"Bu1_-5--R6obQc"}],"status":8635,"title":"c-n___rgG-__*__","type":"@__ynA-5-*-P-"}
		{schema: `{"properties":{"invalid-params":{"description":"Cписок","items":{"properties":{"name":{"type":"string"},"reason":{"type":"string"}},"type":"object"},"type":"array"},"status":{"description":"HTTP-код","type":"integer"},"title":{"description":"Человекочитаемое","type":"string"},"type":{"description":"Уникальный","type":"string"}},"type": "object","required":["type","status","title"]}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := New().Generate([]byte(tt.schema))
			require.NoError(t, err)
			require.True(t, isValidJson(data))

			//fmt.Println(string(data))
		})
	}
}

func isValidJson(content []byte) bool {
	m := map[string]interface{}{}
	if err := json.Unmarshal(content, &m); err == nil {
		return true
	}
	return false
}
