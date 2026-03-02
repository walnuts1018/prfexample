package handler

import "github.com/labstack/echo/v5"

func (h Handler) AssertLinks(c *echo.Context) error {
	return c.JSON(200, []map[string]any{
		{
			"relation": []string{
				"delegate_permission/common.handle_all_urls",
				"delegate_permission/common.get_login_creds",
			},
			"target": map[string]any{
				"namespace":    "android_app",
				"package_name": "dev.walnuts.test.prf_example",
				"sha256_cert_fingerprints": []string{
					"95:01:CA:38:FD:57:9B:D7:DD:21:2F:D1:F1:16:E0:76:4B:48:C5:36:AE:71:67:6A:2C:56:92:7B:B4:09:4D:63", // TODO: ~/.android/debug.keystore
				},
			},
		},
	})
}
