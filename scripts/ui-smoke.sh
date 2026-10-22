#!/usr/bin/env bash
# UI smoke: sticky checkboxes + Quill re-mount after SPA navigation.
# Usage: ./scripts/ui-smoke.sh [base_url]   default https://127.0.0.1:8443
set -uo pipefail

BASE="${1:-https://127.0.0.1:8443}"
PWCLI="/home/chubuk/.local/share/mimocode/builtin_skills/0.1.15/skills/playwright/scripts/playwright_cli.sh"
CFG="${HOME}/.playwright/cli.config.json"
PASS=0
FAIL=0

say() { printf '%s\n' "$*"; }
ok()  { PASS=$((PASS+1)); say "  OK  $*"; }
bad() { FAIL=$((FAIL+1)); say "  FAIL $*"; }

pw() { bash "$PWCLI" "$@"; }
pev() { pw eval --raw "$1" 2>/dev/null | tail -1 | tr -d '"' | tr -d '\r'; }

say "UI smoke against $BASE"

pw open "$BASE/login" --config "$CFG" >/dev/null
pw fill '#login' 'admin' >/dev/null
pw fill '#password' 'TestPass123!' >/dev/null
pw click 'button:has-text("Войти")' >/dev/null 2>&1 || pw click 'text=Войти' >/dev/null 2>&1 || true
sleep 3
path=$(pev 'location.pathname')
case "$path" in
  *dashboard*) ok "login (path=$path)" ;;
  *) bad "login (path=$path)"; pw close >/dev/null; say "smoke: abort"; exit 1 ;;
esac

# --- PDN sticky ---
pw goto "$BASE/admin/pdn" >/dev/null
sleep 2
if [ "$(pev 'document.getElementById("pdn-enabled") ? "1" : "0"')" != "1" ]; then
  bad "PDN page: checkboxes missing"
else
  pev 'document.getElementById("pdn-enabled").click(); "ok"' >/dev/null
  pev 'document.getElementById("pdn-log").click(); "ok"' >/dev/null
  toggled=$(pev 'JSON.stringify([document.getElementById("pdn-enabled").checked,document.getElementById("pdn-log").checked])')
  pev 'Admin.savePDN(); "ok"' >/dev/null
  sleep 3
  after=$(pev 'JSON.stringify([document.getElementById("pdn-enabled").checked,document.getElementById("pdn-log").checked])')
  if [ -z "$toggled" ] || [ "$toggled" = "[null,null]" ]; then
    bad "PDN: cannot read checkboxes ($toggled)"
  elif [ "$toggled" = "$after" ]; then
    ok "PDN checkboxes sticky ($after)"
  else
    bad "PDN bounced: $toggled -> $after"
  fi
fi

# --- server auto_renew sticky ---
pw goto "$BASE/admin/server" >/dev/null
sleep 2
if [ "$(pev 'document.getElementById("srv-auto-renew") ? "1" : "0"')" != "1" ]; then
  bad "server page: auto_renew missing"
else
  pev 'document.getElementById("srv-auto-renew").click(); "ok"' >/dev/null
  srv_t=$(pev 'String(document.getElementById("srv-auto-renew").checked)')
  pev 'Admin.saveServer(); "ok"' >/dev/null
  sleep 3
  srv_a=$(pev 'String(document.getElementById("srv-auto-renew").checked)')
  if [ "$srv_t" = "$srv_a" ] && [ -n "$srv_t" ] && [ "$srv_t" != "null" ]; then
    ok "auto_renew sticky ($srv_a)"
  else
    bad "auto_renew bounced: '$srv_t' -> '$srv_a'"
  fi
fi

# --- Quill visual after SPA re-mount ---
pw goto "$BASE/admin/email" >/dev/null
sleep 4
q1=$(pev '(EmailAdmin.quill && EmailAdmin.quill.container && EmailAdmin.quill.container.isConnected) ? String((EmailAdmin.quill.getText()||"").trim().length) : "-1"')
pev 'App.navigate("/dashboard"); "ok"' >/dev/null
sleep 1
pev 'App.navigate("/admin/email"); "ok"' >/dev/null
sleep 4
q2=$(pev '(EmailAdmin.quill && EmailAdmin.quill.container && EmailAdmin.quill.container.isConnected) ? String((EmailAdmin.quill.getText()||"").trim().length) : "-1"')
vis=$(pev 'String((document.getElementById("tpl-quill")||{innerText:""}).innerText.trim().length)')
if [ "${q1:-0}" -gt 10 ] 2>/dev/null && [ "${q2:-0}" -gt 10 ] 2>/dev/null && [ "${vis:-0}" -gt 10 ] 2>/dev/null; then
  ok "Quill content after SPA re-mount (len=$q2 visible=$vis)"
else
  bad "Quill empty after re-mount (first=$q1 second=$q2 visible=$vis)"
fi

pw close >/dev/null || true

say "----"
say "smoke: $PASS ok, $FAIL failed"
[ "$FAIL" -eq 0 ]
