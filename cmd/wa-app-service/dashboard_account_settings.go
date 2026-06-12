package main

import (
	"net/http"

	waappv1 "github.com/byte-v-forge/wa-app/gen/go/byte/v/forge/waapp/v1"
)

func (s *dashboardHTTP) handleGetTwoFactorAuthStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, r, http.MethodGet)
		return
	}
	if s.service == nil {
		writeDashboardError(w, r, http.StatusServiceUnavailable, "common.service_not_configured", "wa-app 服务未配置")
		return
	}
	payload := queryPayload(r)
	resp, err := s.service.GetTwoFactorAuthStatus(r.Context(), &waappv1.GetTwoFactorAuthStatusRequest{
		Context:       &waappv1.RequestContext{RequestId: newRequestID("wa-account-2fa-status")},
		Selector:      accountSettingsSelector(payload),
		RemoteRefresh: boolField(payload, "remote_refresh"),
	})
	if err != nil {
		writeDashboardError(w, r, http.StatusInternalServerError, "op.get_wa_2fa_status_failed", "获取 WA 2FA 状态失败")
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (s *dashboardHTTP) handleSetTwoFactorAuthSettings(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccountSettingsPost(w, r) {
		return
	}
	payload, ok := readAccountSettingsPayload(w, r)
	if !ok {
		return
	}
	resp, err := s.service.SetTwoFactorAuthSettings(r.Context(), &waappv1.SetTwoFactorAuthSettingsRequest{
		Context:  accountSettingsRequestContext(payload, "wa-account-2fa"),
		Selector: accountSettingsSelector(payload),
		Pin:      &waappv1.SensitiveText{Value: textField(payload, "pin")},
	})
	if err != nil {
		writeDashboardError(w, r, http.StatusInternalServerError, "op.set_wa_2fa_settings_failed", "设置 WA 2FA 失败")
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (s *dashboardHTTP) handleSetAccountEmail(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccountSettingsPost(w, r) {
		return
	}
	payload, ok := readAccountSettingsPayload(w, r)
	if !ok {
		return
	}
	resp, err := s.service.SetAccountEmail(r.Context(), &waappv1.SetAccountEmailRequest{
		Context:       accountSettingsRequestContext(payload, "wa-account-email"),
		Selector:      accountSettingsSelector(payload),
		EmailAddress:  textField(payload, "email_address"),
		GoogleIdToken: &waappv1.SensitiveText{Value: textField(payload, "google_id_token")},
	})
	if err != nil {
		writeDashboardError(w, r, http.StatusInternalServerError, "op.set_wa_account_email_failed", "设置 WA 账号邮箱失败")
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (s *dashboardHTTP) handleRequestAccountEmailOtp(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccountSettingsPost(w, r) {
		return
	}
	payload, ok := readAccountSettingsPayload(w, r)
	if !ok {
		return
	}
	resp, err := s.service.RequestAccountEmailOtp(r.Context(), &waappv1.RequestAccountEmailOtpRequest{
		Context:        accountSettingsRequestContext(payload, "wa-account-email-otp"),
		Selector:       accountSettingsSelector(payload),
		LocaleLanguage: firstNonEmpty(textField(payload, "locale_language"), textField(payload, "language"), "en"),
		LocaleCountry:  firstNonEmpty(textField(payload, "locale_country"), textField(payload, "country"), "US"),
	})
	if err != nil {
		writeDashboardError(w, r, http.StatusInternalServerError, "op.request_wa_account_email_otp_failed", "请求 WA 账号邮箱 OTP 失败")
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (s *dashboardHTTP) handleVerifyAccountEmailOtp(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccountSettingsPost(w, r) {
		return
	}
	payload, ok := readAccountSettingsPayload(w, r)
	if !ok {
		return
	}
	resp, err := s.service.VerifyAccountEmailOtp(r.Context(), &waappv1.VerifyAccountEmailOtpRequest{
		Context:  accountSettingsRequestContext(payload, "wa-account-email-otp-verify"),
		Selector: accountSettingsSelector(payload),
		Code:     &waappv1.SensitiveText{Value: firstNonEmpty(textField(payload, "code"), textField(payload, "otp"))},
	})
	if err != nil {
		writeDashboardError(w, r, http.StatusInternalServerError, "op.verify_wa_account_email_otp_failed", "校验 WA 账号邮箱 OTP 失败")
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (s *dashboardHTTP) handleSetAccountProfileName(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccountSettingsPost(w, r) {
		return
	}
	req := &waappv1.SetAccountProfileNameRequest{}
	if !readProtoJSONPayload(w, r, 1<<20, req, "proto.set_profile_name_request_json", "请求体必须是 SetAccountProfileNameRequest JSON 对象") {
		return
	}
	ensureAccountSettingsContext(&req.Context, "wa-account-profile-name")
	resp, err := s.service.SetAccountProfileName(r.Context(), req)
	if err != nil {
		writeDashboardError(w, r, http.StatusInternalServerError, "op.set_wa_account_profile_name_failed", "设置 WA 账号资料名称失败")
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (s *dashboardHTTP) handleSetAccountProfilePicture(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccountSettingsPost(w, r) {
		return
	}
	req := &waappv1.SetAccountProfilePictureRequest{}
	if !readProtoJSONPayload(w, r, 4<<20, req, "proto.set_profile_picture_request_json", "请求体必须是 SetAccountProfilePictureRequest JSON 对象") {
		return
	}
	ensureAccountSettingsContext(&req.Context, "wa-account-profile-picture")
	resp, err := s.service.SetAccountProfilePicture(r.Context(), req)
	if err != nil {
		writeDashboardError(w, r, http.StatusInternalServerError, "op.set_wa_account_profile_picture_failed", "设置 WA 账号资料图片失败")
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (s *dashboardHTTP) handleRemoveAccountProfilePicture(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccountSettingsPost(w, r) {
		return
	}
	req := &waappv1.RemoveAccountProfilePictureRequest{}
	if !readProtoJSONPayload(w, r, 1<<20, req, "proto.remove_profile_picture_request_json", "请求体必须是 RemoveAccountProfilePictureRequest JSON 对象") {
		return
	}
	ensureAccountSettingsContext(&req.Context, "wa-account-profile-picture-remove")
	resp, err := s.service.RemoveAccountProfilePicture(r.Context(), req)
	if err != nil {
		writeDashboardError(w, r, http.StatusInternalServerError, "op.remove_wa_account_profile_picture_failed", "移除 WA 账号资料图片失败")
		return
	}
	writeProtoJSON(w, http.StatusOK, resp)
}

func (s *dashboardHTTP) requireAccountSettingsPost(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, r, http.MethodPost)
		return false
	}
	if s.service == nil {
		writeDashboardError(w, r, http.StatusServiceUnavailable, "common.service_not_configured", "wa-app 服务未配置")
		return false
	}
	return true
}

func readAccountSettingsPayload(w http.ResponseWriter, r *http.Request) (map[string]any, bool) {
	return readJSONPayload(w, r)
}

func queryPayload(r *http.Request) map[string]any {
	payload := map[string]any{}
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			payload[key] = values[len(values)-1]
		}
	}
	return payload
}

func accountSettingsRequestContext(payload map[string]any, prefix string) *waappv1.RequestContext {
	return &waappv1.RequestContext{
		RequestId: firstNonEmpty(textField(payload, "request_id"), newRequestID(prefix)),
	}
}

func ensureAccountSettingsContext(ctx **waappv1.RequestContext, prefix string) {
	if *ctx == nil {
		*ctx = &waappv1.RequestContext{}
	}
	if (*ctx).RequestId == "" {
		(*ctx).RequestId = newRequestID(prefix)
	}
}

func accountSettingsSelector(payload map[string]any) *waappv1.AccountLoginSelector {
	selector := objectField(payload, "selector")
	return &waappv1.AccountLoginSelector{
		LoginStateId:         firstNonEmpty(textField(payload, "login_state_id"), textField(selector, "login_state_id")),
		RegisteredIdentityId: firstNonEmpty(textField(payload, "registered_identity_id"), textField(selector, "registered_identity_id")),
		WaAccountId:          firstNonEmpty(textField(payload, "wa_account_id"), textField(selector, "wa_account_id")),
		ClientProfileId:      firstNonEmpty(textField(payload, "client_profile_id"), textField(selector, "client_profile_id")),
	}
}
