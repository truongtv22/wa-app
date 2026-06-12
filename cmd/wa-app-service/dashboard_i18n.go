package main

import "net/http"

var dashboardTranslations = map[string]map[string]string{
	"en": {
		"common.method_not_allowed":                   "method not allowed",
		"common.invalid_request_body":                 "invalid request body",
		"common.request_body_must_be_json":            "request body must be json",
		"common.service_not_configured":               "wa-app service is not configured",
		"common.authentication_required":              "authentication required",
		"account.wa_account_id_required":              "wa_account_id is required",
		"contact.id_required":                         "contact id is required",
		"login.create_session_failed":                 "create dashboard session failed",
		"health.workflow_register_native":             "WA native registration flow",
		"workflow.phone_required":                     "phone is required",
		"workflow.country_calling_code_required":      "country_calling_code is required",
		"workflow.phone_parse_failed":                 "phone parse failed",
		"workflow.phone_country_code_mismatch":        "phone country calling code does not match country_calling_code",
		"workflow.phone_not_possible":                 "the phone number length does not match the country rule. Check the country calling code and phone number.",
		"op.load_long_connection_status_failed":       "load long connection status failed",
		"op.load_wa_accounts_failed":                  "load WA accounts failed",
		"op.load_wa_account_failed":                   "load WA account failed",
		"op.delete_wa_account_failed":                 "delete WA account failed",
		"op.create_wa_account_failed":                 "create WA account failed",
		"op.load_wa_client_profiles_failed":           "load WA client profiles failed",
		"op.load_wa_otp_history_failed":               "load WA OTP history failed",
		"op.load_wa_messages_failed":                  "load WA messages failed",
		"op.mark_wa_messages_read_failed":             "mark WA messages read failed",
		"op.delete_wa_messages_failed":                "delete WA messages failed",
		"op.send_wa_text_message_failed":              "send WA text message failed",
		"op.load_wa_contacts_failed":                  "load WA contacts failed",
		"op.delete_wa_contact_failed":                 "delete WA contact failed",
		"op.resolve_wa_contacts_failed":               "resolve WA contacts failed",
		"op.probe_wa_phone_failed":                    "probe WA phone failed",
		"op.build_login_state_check_request_failed":   "build login-state check request failed",
		"op.get_wa_2fa_status_failed":                 "get WA 2FA status failed",
		"op.set_wa_2fa_settings_failed":               "set WA 2FA settings failed",
		"op.set_wa_account_email_failed":              "set WA account email failed",
		"op.request_wa_account_email_otp_failed":      "request WA account email OTP failed",
		"op.verify_wa_account_email_otp_failed":       "verify WA account email OTP failed",
		"op.set_wa_account_profile_name_failed":       "set WA account profile name failed",
		"op.set_wa_account_profile_picture_failed":    "set WA account profile picture failed",
		"op.remove_wa_account_profile_picture_failed": "remove WA account profile picture failed",
		"proto.set_profile_name_request_json":         "request body must be a SetAccountProfileNameRequest JSON object",
		"proto.set_profile_picture_request_json":      "request body must be a SetAccountProfilePictureRequest JSON object",
		"proto.remove_profile_picture_request_json":   "request body must be a RemoveAccountProfilePictureRequest JSON object",
		"proto.resolve_contacts_request_json":         "request body must be a ResolveWAContactsRequest JSON object",
	},
	"vi": {
		"common.method_not_allowed":                   "phương thức không được hỗ trợ",
		"common.invalid_request_body":                 "request body không hợp lệ",
		"common.request_body_must_be_json":            "request body phải là json",
		"common.service_not_configured":               "dịch vụ wa-app chưa được cấu hình",
		"common.authentication_required":              "cần xác thực",
		"account.wa_account_id_required":              "thiếu wa_account_id",
		"contact.id_required":                         "thiếu contact id",
		"login.create_session_failed":                 "tạo session dashboard thất bại",
		"health.workflow_register_native":             "Luồng đăng ký native WA",
		"workflow.phone_required":                     "thiếu số điện thoại",
		"workflow.country_calling_code_required":      "thiếu country_calling_code",
		"workflow.phone_parse_failed":                 "phân tích số điện thoại thất bại",
		"workflow.phone_country_code_mismatch":        "mã gọi quốc gia của số điện thoại không khớp country_calling_code",
		"workflow.phone_not_possible":                 "độ dài số điện thoại không khớp quy tắc quốc gia. Hãy kiểm tra mã gọi quốc gia và số điện thoại.",
		"op.load_long_connection_status_failed":       "tải trạng thái kết nối dài thất bại",
		"op.load_wa_accounts_failed":                  "tải danh sách tài khoản WA thất bại",
		"op.load_wa_account_failed":                   "tải tài khoản WA thất bại",
		"op.delete_wa_account_failed":                 "xóa tài khoản WA thất bại",
		"op.create_wa_account_failed":                 "tạo tài khoản WA thất bại",
		"op.load_wa_client_profiles_failed":           "tải client profile WA thất bại",
		"op.load_wa_otp_history_failed":               "tải lịch sử OTP WA thất bại",
		"op.load_wa_messages_failed":                  "tải tin nhắn WA thất bại",
		"op.mark_wa_messages_read_failed":             "đánh dấu đã đọc tin nhắn WA thất bại",
		"op.delete_wa_messages_failed":                "xóa tin nhắn WA thất bại",
		"op.send_wa_text_message_failed":              "gửi tin nhắn văn bản WA thất bại",
		"op.load_wa_contacts_failed":                  "tải liên hệ WA thất bại",
		"op.delete_wa_contact_failed":                 "xóa liên hệ WA thất bại",
		"op.resolve_wa_contacts_failed":               "resolve liên hệ WA thất bại",
		"op.probe_wa_phone_failed":                    "dò số WA thất bại",
		"op.build_login_state_check_request_failed":   "dựng request kiểm tra trạng thái đăng nhập thất bại",
		"op.get_wa_2fa_status_failed":                 "lấy trạng thái 2FA WA thất bại",
		"op.set_wa_2fa_settings_failed":               "đặt cấu hình 2FA WA thất bại",
		"op.set_wa_account_email_failed":              "đặt email tài khoản WA thất bại",
		"op.request_wa_account_email_otp_failed":      "yêu cầu OTP email tài khoản WA thất bại",
		"op.verify_wa_account_email_otp_failed":       "xác minh OTP email tài khoản WA thất bại",
		"op.set_wa_account_profile_name_failed":       "đặt tên hồ sơ tài khoản WA thất bại",
		"op.set_wa_account_profile_picture_failed":    "đặt ảnh hồ sơ tài khoản WA thất bại",
		"op.remove_wa_account_profile_picture_failed": "xóa ảnh hồ sơ tài khoản WA thất bại",
		"proto.set_profile_name_request_json":         "request body phải là JSON của SetAccountProfileNameRequest",
		"proto.set_profile_picture_request_json":      "request body phải là JSON của SetAccountProfilePictureRequest",
		"proto.remove_profile_picture_request_json":   "request body phải là JSON của RemoveAccountProfilePictureRequest",
		"proto.resolve_contacts_request_json":         "request body phải là JSON của ResolveWAContactsRequest",
	},
}

func dashboardT(r *http.Request, key string, fallback string) string {
	return dashboardTLocale(dashboardLocaleFromRequest(r), key, fallback)
}

func dashboardTLocale(locale string, key string, fallback string) string {
	if locale == "zh-CN" {
		return fallback
	}
	if translated := dashboardTranslations[locale][key]; translated != "" {
		return translated
	}
	return fallback
}

func writeDashboardError(w http.ResponseWriter, r *http.Request, status int, key string, fallback string) {
	writeJSON(w, status, map[string]string{"error": dashboardT(r, key, fallback)})
}
