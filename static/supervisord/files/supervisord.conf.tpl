[program:mock-{{ .Domain }}-{{ .Port }}]
environment = ROOT="{{ .Root }}",PORT="{{ .Port }}",NAME="{{ .Domain }}"
autorestart = true
redirect_stderr = true
stdout_logfile = /var/log/supervisord/mock-{{ .Domain }}.log
command = java -cp "/var/wiremock/lib/*:/var/wiremock/extensions/*" wiremock.Run --port {{ .Port }} --root-dir {{ .Root }} --global-response-templating --record-mappings --verbose
