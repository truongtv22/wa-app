package main

import (
	"html/template"
	"net/http"
	"strings"
)

const dashboardLoginMaxFormBytes = 4096

type dashboardLoginPageData struct {
	Next  string
	Error string
	Copy  dashboardLoginCopy
}

type dashboardLoginCopy struct {
	Lang          string
	Title         string
	Heading       string
	Subtitle      string
	PasswordLabel string
	SubmitLabel   string
	Hint          string
	InvalidForm   string
	WrongPassword string
}

const dashboardLocaleCookieName = "wa_lang"

var dashboardLoginCopies = map[string]dashboardLoginCopy{
	"zh-CN": {
		Lang:          "zh-CN",
		Title:         "WA 登录",
		Heading:       "登录 WA 管理",
		Subtitle:      "输入访问密码，登录状态保留 7 天",
		PasswordLabel: "访问密码",
		SubmitLabel:   "登录",
		Hint:          "关闭浏览器后仍会保持登录；如需退出可访问 /logout。",
		InvalidForm:   "登录请求无效",
		WrongPassword: "密码不正确",
	},
	"en": {
		Lang:          "en",
		Title:         "WA Sign In",
		Heading:       "Sign in to WA admin",
		Subtitle:      "Enter the access password. The session lasts 7 days.",
		PasswordLabel: "Access password",
		SubmitLabel:   "Sign in",
		Hint:          "The session stays active after closing the browser. Visit /logout to sign out.",
		InvalidForm:   "Invalid sign-in request",
		WrongPassword: "Incorrect password",
	},
	"vi": {
		Lang:          "vi",
		Title:         "Đăng nhập WA",
		Heading:       "Đăng nhập quản trị WA",
		Subtitle:      "Nhập mật khẩu truy cập, phiên đăng nhập giữ trong 7 ngày.",
		PasswordLabel: "Mật khẩu truy cập",
		SubmitLabel:   "Đăng nhập",
		Hint:          "Đóng trình duyệt vẫn giữ đăng nhập; truy cập /logout để đăng xuất.",
		InvalidForm:   "Yêu cầu đăng nhập không hợp lệ",
		WrongPassword: "Mật khẩu không đúng",
	},
}

var dashboardLoginTemplate = template.Must(template.New("dashboard-login").Parse(dashboardLoginHTML))

func handleDashboardLogin(w http.ResponseWriter, r *http.Request, auth dashboardAuthConfig) {
	next := dashboardSafeRedirectTarget(r.URL.Query().Get("next"))
	locale := dashboardLocaleFromRequest(r)
	setDashboardLocaleCookie(w, locale)
	switch r.Method {
	case http.MethodGet:
		if authenticatedDashboardRequest(r, auth) {
			http.Redirect(w, r, next, http.StatusSeeOther)
			return
		}
		renderDashboardLogin(w, http.StatusOK, next, "", locale)
	case http.MethodPost:
		handleDashboardLoginPost(w, r, auth, next, locale)
	default:
		w.Header().Set("Allow", http.MethodGet+", "+http.MethodPost)
		writeDashboardError(w, r, http.StatusMethodNotAllowed, "common.method_not_allowed", "方法不允许")
	}
}

func handleDashboardLoginPost(w http.ResponseWriter, r *http.Request, auth dashboardAuthConfig, next string, fallbackLocale string) {
	r.Body = http.MaxBytesReader(w, r.Body, dashboardLoginMaxFormBytes)
	if err := r.ParseForm(); err != nil {
		renderDashboardLogin(w, http.StatusBadRequest, next, dashboardLoginCopyFor(fallbackLocale).InvalidForm, fallbackLocale)
		return
	}
	locale := normalizeDashboardLocale(r.FormValue("lang"))
	if locale == "" {
		locale = fallbackLocale
	}
	setDashboardLocaleCookie(w, locale)
	next = dashboardSafeRedirectTarget(r.FormValue("next"))
	if !validDashboardPassword(r.FormValue("password"), auth) {
		renderDashboardLogin(w, http.StatusUnauthorized, next, dashboardLoginCopyFor(locale).WrongPassword, locale)
		return
	}
	cookie, err := newDashboardSessionCookie(r, auth)
	if err != nil {
		writeDashboardError(w, r, http.StatusInternalServerError, "login.create_session_failed", "创建 dashboard 会话失败")
		return
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, next, http.StatusSeeOther)
}

func handleDashboardLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodGet+", "+http.MethodPost)
		writeDashboardError(w, r, http.StatusMethodNotAllowed, "common.method_not_allowed", "方法不允许")
		return
	}
	http.SetCookie(w, expiredDashboardSessionCookie(r))
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func renderDashboardLogin(w http.ResponseWriter, status int, next string, message string, locale string) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_ = dashboardLoginTemplate.Execute(w, dashboardLoginPageData{Next: next, Error: message, Copy: dashboardLoginCopyFor(locale)})
}

func dashboardLocaleFromRequest(r *http.Request) string {
	if locale := normalizeDashboardLocale(r.URL.Query().Get("lang")); locale != "" {
		return locale
	}
	if cookie, err := r.Cookie(dashboardLocaleCookieName); err == nil {
		if locale := normalizeDashboardLocale(cookie.Value); locale != "" {
			return locale
		}
	}
	if locale := normalizeDashboardLocale(r.Header.Get("Accept-Language")); locale != "" {
		return locale
	}
	return "zh-CN"
}

func normalizeDashboardLocale(raw string) string {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "" {
		return ""
	}
	switch {
	case strings.HasPrefix(value, "vi"):
		return "vi"
	case strings.HasPrefix(value, "en"):
		return "en"
	case strings.HasPrefix(value, "zh"):
		return "zh-CN"
	default:
		return ""
	}
}

func setDashboardLocaleCookie(w http.ResponseWriter, locale string) {
	http.SetCookie(w, &http.Cookie{
		Name:     dashboardLocaleCookieName,
		Value:    locale,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 365,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	})
}

func dashboardLoginCopyFor(locale string) dashboardLoginCopy {
	if copy, ok := dashboardLoginCopies[locale]; ok {
		return copy
	}
	return dashboardLoginCopies["zh-CN"]
}

const dashboardLoginHTML = `<!doctype html>
<html lang="{{.Copy.Lang}}">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>{{.Copy.Title}}</title>
<style>
:root{color-scheme:light dark;font-family:Inter,ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;background:#07130d;color:#e9fff0}body{min-height:100dvh;margin:0;display:grid;place-items:center;background:radial-gradient(circle at top,#14532d 0,#07130d 48%,#020617 100%)}main{width:min(92vw,380px);border:1px solid rgba(255,255,255,.12);border-radius:24px;background:rgba(2,6,23,.76);box-shadow:0 24px 70px rgba(0,0,0,.36);padding:32px;backdrop-filter:blur(14px)}.brand{display:flex;align-items:center;gap:12px;margin-bottom:24px}.mark{display:grid;place-items:center;width:44px;height:44px;border-radius:14px;background:#22c55e;color:#052e16;font-weight:800}.title{margin:0;font-size:22px}.sub{margin:4px 0 0;color:#a7f3d0;font-size:13px}label{display:grid;gap:8px;margin-top:16px;color:#bbf7d0;font-size:13px}input{box-sizing:border-box;width:100%;border:1px solid rgba(187,247,208,.22);border-radius:12px;background:rgba(15,23,42,.92);color:#f8fafc;padding:12px 14px;font:inherit;outline:none}input:focus{border-color:#22c55e;box-shadow:0 0 0 3px rgba(34,197,94,.18)}button{width:100%;margin-top:22px;border:0;border-radius:12px;background:#22c55e;color:#052e16;padding:12px 14px;font:inherit;font-weight:700;cursor:pointer}button:hover{background:#4ade80}.hint{margin:14px 0 0;color:#86efac;font-size:12px;line-height:1.6}.error{margin:0 0 14px;border:1px solid rgba(248,113,113,.28);border-radius:12px;background:rgba(127,29,29,.32);color:#fecaca;padding:10px 12px;font-size:13px}
.langs{display:flex;justify-content:flex-end;gap:10px;margin-bottom:18px;font-size:12px}.langs a{color:#86efac;text-decoration:none}.langs a:hover{text-decoration:underline}
</style>
</head>
<body>
<main>
  <div class="langs"><a href="/login?lang=zh-CN&amp;next={{.Next}}">中文</a><a href="/login?lang=en&amp;next={{.Next}}">EN</a><a href="/login?lang=vi&amp;next={{.Next}}">VI</a></div>
  <div class="brand"><div class="mark">WA</div><div><h1 class="title">{{.Copy.Heading}}</h1><p class="sub">{{.Copy.Subtitle}}</p></div></div>
  {{if .Error}}<p class="error">{{.Error}}</p>{{end}}
  <form method="post" action="/login" autocomplete="on">
    <input type="hidden" name="next" value="{{.Next}}">
    <input type="hidden" name="lang" value="{{.Copy.Lang}}">
    <label>{{.Copy.PasswordLabel}}<input name="password" type="password" autocomplete="current-password" required autofocus></label>
    <button type="submit">{{.Copy.SubmitLabel}}</button>
    <p class="hint">{{.Copy.Hint}}</p>
  </form>
</main>
</body>
</html>`
