package g79client

import "testing"

// TestSelectAndroidPatchVersion 覆盖 patch 版本选择与降级逻辑。
//
// 背景：网易会先更新 packlist 的 netease.version（例如 3.10.0.300878），
// 而 patchlist 与 patch 资源仍停留在旧版本（例如 3.9.41.300665）。
// 此时不能直接报错，否则整个 g79 客户端初始化失败、登录不可用。
func TestSelectAndroidPatchVersion(t *testing.T) {
	original := EngineVersion
	defer func() { EngineVersion = original }()

	versions := []string{
		"2.12.49.254480",
		"3.7.46.293055",
		"3.8.49.297835",
		"3.9.20.297977",
		"3.9.41.300665",
	}

	tests := []struct {
		name    string
		engine  string
		want    string
		wantErr bool
	}{
		{name: "前缀一致时取最新的同前缀版本", engine: "3.9.20.222222", want: "3.9.41.300665"},
		{name: "引擎版本为 3.10 时降级到 3.x 最新", engine: "3.10.0.300878", want: "3.9.41.300665"},
		{name: "完全无匹配时回退到最新 patch", engine: "9.9.9.999999", want: "3.9.41.300665"},
		{name: "非法引擎版本报错", engine: "3", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			EngineVersion = tt.engine
			got, err := selectAndroidPatchVersion(versions)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("期望报错，实际返回 %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("未期望的报错: %v", err)
			}
			if got != tt.want {
				t.Fatalf("选择结果错误: got=%q want=%q", got, tt.want)
			}
		})
	}

	if _, err := selectAndroidPatchVersion(nil); err == nil {
		t.Fatal("空版本列表应当报错")
	}
}
