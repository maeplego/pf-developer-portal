package httpapi

import "html/template"

func mustParse() *template.Template {
	return template.Must(template.New("page").Parse(pageHTML))
}

const pageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.Title}} · pf-developer-portal</title>
  <style>
    @import url("https://fonts.googleapis.com/css2?family=DM+Sans:ital,opsz,wght@0,9..40,400;0,9..40,500;0,9..40,600;0,9..40,700;1,9..40,400&display=swap");
    :root {
      color-scheme: light dark;
      --font: "DM Sans", system-ui, -apple-system, "Segoe UI", sans-serif;
      --bg: #f4f6fb;
      --bg-accent: radial-gradient(1200px 600px at 10% -10%, rgba(99, 102, 241, 0.12), transparent 55%),
        radial-gradient(900px 500px at 90% 0%, rgba(14, 165, 233, 0.1), transparent 50%), var(--bg);
      --surface: rgba(255, 255, 255, 0.82);
      --panel: #ffffff;
      --line: rgba(15, 23, 42, 0.08);
      --text: #0f172a;
      --muted: #64748b;
      --accent: #4f46e5;
      --accent-hover: #4338ca;
      --accent-soft: rgba(79, 70, 229, 0.12);
      --code-bg: #f8fafc;
      --get: #059669;
      --post: #4f46e5;
      --put: #d97706;
      --del: #dc2626;
      --radius: 14px;
      --radius-sm: 10px;
      --shadow-sm: 0 4px 16px rgba(15, 23, 42, 0.06);
      --header-h: 64px;
    }
    @media (prefers-color-scheme: dark) {
      :root {
        --bg: #0b1020;
        --bg-accent: radial-gradient(1200px 600px at 10% -10%, rgba(99, 102, 241, 0.18), transparent 55%),
          radial-gradient(900px 500px at 90% 0%, rgba(14, 165, 233, 0.12), transparent 50%), var(--bg);
        --surface: rgba(15, 23, 42, 0.72);
        --panel: #111827;
        --line: rgba(148, 163, 184, 0.14);
        --text: #e2e8f0;
        --muted: #94a3b8;
        --accent: #818cf8;
        --accent-hover: #a5b4fc;
        --accent-soft: rgba(129, 140, 248, 0.16);
        --code-bg: #0b0e13;
        --post: #818cf8;
        --shadow-sm: 0 6px 20px rgba(0, 0, 0, 0.25);
      }
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      font: 15px/1.55 var(--font);
      background: var(--bg-accent);
      color: var(--text);
      -webkit-font-smoothing: antialiased;
    }
    header {
      display: flex;
      align-items: center;
      gap: 1rem;
      min-height: var(--header-h);
      padding: 0 1.25rem;
      border-bottom: 1px solid var(--line);
      background: var(--surface);
      backdrop-filter: blur(14px);
      position: sticky; top: 0; z-index: 20;
    }
    header a { color: var(--accent); text-decoration: none; font-weight: 650; }
    header a:hover { color: var(--accent-hover); }
    header span { color: var(--muted); font-size: 13px; }
    .layout { display: grid; grid-template-columns: 260px 1fr; min-height: calc(100vh - var(--header-h)); }
    nav {
      border-right: 1px solid var(--line);
      padding: 1rem 0.8rem 2rem;
      background: var(--panel);
    }
    nav h2 { font-size: 11px; letter-spacing: .12em; text-transform: uppercase; color: var(--muted); margin: 0 0 .6rem .4rem; }
    nav a {
      display: block; color: var(--text); text-decoration: none;
      padding: .4rem .55rem; border-radius: 999px; margin-bottom: 2px;
    }
    nav a:hover, nav a.active { background: var(--accent-soft); color: var(--text); }
    nav .src { display:block; color: var(--muted); font-size: 12px; }
    main { padding: 1.5rem 2rem 4rem; max-width: 920px; }
    .hero {
      padding: 1.25rem 1.35rem;
      border-radius: calc(var(--radius) + 4px);
      background: linear-gradient(135deg, rgba(79, 70, 229, 0.12), rgba(14, 165, 233, 0.08));
      border: 1px solid var(--line);
      margin-bottom: 1.25rem;
    }
    .hero h1 { margin: 0 0 .4rem; font-size: clamp(1.5rem, 2.5vw, 1.8rem); letter-spacing: -0.03em; }
    .hero p { color: var(--muted); margin: 0 0 1rem; }
    .hero p:last-child { margin-bottom: 0; }
    .cards { display: grid; gap: 1rem; grid-template-columns: repeat(auto-fill, minmax(240px, 1fr)); }
    .card {
      background: var(--panel); border: 1px solid var(--line); border-radius: var(--radius);
      box-shadow: var(--shadow-sm); padding: 1rem 1.1rem; text-decoration: none; color: inherit; display: block;
      transition: border-color 0.15s ease, box-shadow 0.15s ease;
    }
    .card:hover { border-color: var(--accent); box-shadow: 0 8px 24px rgba(79, 70, 229, 0.12); }
    .badge {
      display: inline-block; min-width: 3.4rem; text-align: center;
      font: 700 11px/1.2 ui-monospace, SFMono-Regular, Consolas, monospace;
      padding: .22rem .45rem; border-radius: 999px; margin-right: .55rem;
    }
    .GET { background: rgba(5, 150, 105, 0.15); color: var(--get); }
    .POST { background: var(--accent-soft); color: var(--post); }
    .PUT, .PATCH { background: rgba(217, 119, 6, 0.15); color: var(--put); }
    .DELETE { background: rgba(220, 38, 38, 0.15); color: var(--del); }
    .op {
      border: 1px solid var(--line); border-radius: var(--radius); background: var(--panel);
      box-shadow: var(--shadow-sm); margin: 1.1rem 0; overflow: hidden;
    }
    .op-h { padding: .85rem 1rem; border-bottom: 1px solid var(--line); }
    .path { font-family: ui-monospace, SFMono-Regular, Consolas, monospace; }
    .op-b { padding: 1rem; }
    .desc { color: var(--muted); }
    table { width: 100%; border-collapse: collapse; font-size: 13px; }
    th, td { text-align: left; padding: .5rem .55rem; border-bottom: 1px solid var(--line); }
    th { color: var(--muted); font-size: 0.82rem; font-weight: 600; text-transform: uppercase; letter-spacing: 0.04em; }
    pre {
      background: var(--code-bg); border: 1px solid var(--line); border-radius: var(--radius-sm);
      padding: .8rem; overflow: auto; font-size: 12.5px;
    }
    textarea, input[data-try-path] {
      width: 100%; background: var(--code-bg); color: var(--text);
      border: 1px solid var(--line); border-radius: var(--radius-sm); padding: .7rem;
      font: 12.5px ui-monospace, Consolas, monospace;
    }
    textarea { min-height: 120px; }
    button {
      display: inline-flex; align-items: center; justify-content: center;
      background: linear-gradient(135deg, var(--accent), #6366f1); color: #fff; border: 0; border-radius: 999px;
      font: inherit; font-weight: 600; padding: .55rem 1rem; cursor: pointer;
      box-shadow: 0 8px 24px rgba(79, 70, 229, 0.28);
      transition: transform 0.15s ease, box-shadow 0.15s ease;
    }
    button:hover { transform: translateY(-1px); box-shadow: 0 12px 28px rgba(79, 70, 229, 0.34); }
    .try { margin-top: .8rem; }
    .note { font-size: 13px; color: var(--muted); }
    code { font-family: ui-monospace, SFMono-Regular, Consolas, monospace; font-size: 0.9em; padding: 0.12rem 0.35rem; border-radius: 6px; background: var(--accent-soft); color: var(--accent); }
    @media (max-width: 800px) { .layout { grid-template-columns: 1fr; } nav { display: none; } }
  </style>
</head>
<body>
  <header>
    <a href="/">Developer portal</a>
    <span>P11 · hand-placed OpenAPI · mock has no side effects</span>
  </header>
  <div class="layout">
    <nav>
      <h2>Catalog</h2>
      {{range .APIs}}
        <a href="/docs/{{.Slug}}" {{if and $.API (eq $.API.Slug .Slug)}}class="active"{{end}}>
          {{.Title}}
          <span class="src">{{.Slug}} · {{.Version}}</span>
        </a>
      {{end}}
    </nav>
    <main>
      {{if .Home}}
        <div class="hero">
          <h1>API catalog</h1>
          <p>Specs live in <code>specs/</code>. Reference on the left, mock under <code>/mock/{slug}</code>. Breaking-change CI, GitHub Actions dashboard, and PR review are sibling repos in this portfolio.</p>
        </div>
        <div class="cards">
          {{range .Tools}}
          <a class="card" href="{{.URL}}" target="_blank" rel="noreferrer">
            <strong>{{.Title}}</strong>
            <div class="note">{{.Note}}</div>
          </a>
          {{end}}
        </div>
        <div class="cards" style="margin-top:1rem">
          {{range .APIs}}
            <a class="card" href="/docs/{{.Slug}}">
              <strong>{{.Title}}</strong>
              <div class="note">{{.Summary}} · {{.Source}} · {{len .Paths}} operations</div>
            </a>
          {{end}}
        </div>
      {{else}}
        <div class="hero">
          <h1>{{.API.Title}}</h1>
          <p>{{.API.Summary}} · version {{.API.Version}} · mock base <code>{{.Mock}}</code></p>
          <p class="note">{{.API.FileName}} is the source of truth for this page.</p>
        </div>
        {{range .Ops}}
          <article class="op" id="{{.ID}}">
            <div class="op-h">
              <span class="badge {{.Method}}">{{.Method}}</span>
              <span class="path">{{.Path}}</span>
              {{if .Summary}} — {{.Summary}}{{end}}
            </div>
            <div class="op-b">
              {{if .Description}}<p class="desc">{{.Description}}</p>{{end}}
              {{if .Params}}
                <h3>Parameters</h3>
                <table>
                  <tr><th>Name</th><th>In</th><th>Type</th><th>Required</th></tr>
                  {{range .Params}}
                    <tr><td>{{.Name}}</td><td>{{.In}}</td><td>{{.Type}}</td><td>{{if .Required}}yes{{else}}no{{end}}</td></tr>
                  {{end}}
                </table>
              {{end}}
              {{if .BodyExample}}
                <h3>Request example</h3>
                <pre>{{.BodyExample}}</pre>
              {{end}}
              {{if .RespExample}}
                <h3>Response {{.RespStatus}} example</h3>
                <pre>{{.RespExample}}</pre>
              {{end}}
              <div class="try">
                <h3>Try the mock</h3>
                <p class="note">Calls this process only. Writes do not persist.</p>
                <label class="note">Path</label>
                <input data-try-path value="{{.Path}}" style="margin:.3rem 0 .6rem">
                {{if eq .Method "POST" "PUT" "PATCH"}}
                  <textarea data-try-body>{{.BodyExample}}</textarea>
                {{end}}
                <p><button type="button" data-try-method="{{.Method}}" data-try-mock="{{$.Mock}}">Send {{.Method}}</button></p>
                <pre data-try-out hidden></pre>
              </div>
            </div>
          </article>
        {{end}}
      {{end}}
    </main>
  </div>
  <script>
    document.querySelectorAll("[data-try-method]").forEach((btn) => {
      btn.addEventListener("click", async () => {
        const box = btn.closest(".try");
        const method = btn.dataset.tryMethod;
        const mock = btn.dataset.tryMock;
        let path = box.querySelector("[data-try-path]").value.trim();
        if (!path.startsWith("/")) path = "/" + path;
        const out = box.querySelector("[data-try-out]");
        const bodyEl = box.querySelector("[data-try-body]");
        const headers = { "Accept": "application/json" };
        const init = { method, headers };
        if (bodyEl) {
          headers["Content-Type"] = "application/json";
          headers["Idempotency-Key"] = "demo-try-" + Date.now();
          init.body = bodyEl.value;
        }
        out.hidden = false;
        try {
          const res = await fetch(mock + path, init);
          const text = await res.text();
          out.textContent = res.status + "\\n" + text;
        } catch (err) {
          out.textContent = String(err);
        }
      });
    });
  </script>
</body>
</html>
`
