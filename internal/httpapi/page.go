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
    :root {
      --bg: #0e1116;
      --panel: #161b22;
      --line: #2b3340;
      --text: #e7ecf3;
      --muted: #93a0b4;
      --accent: #7dd3c7;
      --get: #3dd68c;
      --post: #6ea8fe;
      --put: #f5c542;
      --del: #ff7b72;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      font: 15px/1.5 ui-sans-serif, system-ui, Segoe UI, sans-serif;
      background: radial-gradient(1200px 500px at 10% -10%, #1a2a28 0%, transparent 50%), var(--bg);
      color: var(--text);
    }
    header {
      display: flex;
      align-items: center;
      gap: 1rem;
      padding: 0.9rem 1.4rem;
      border-bottom: 1px solid var(--line);
      background: rgba(14,17,22,.85);
      backdrop-filter: blur(8px);
      position: sticky; top: 0; z-index: 2;
    }
    header a { color: var(--accent); text-decoration: none; font-weight: 650; }
    header span { color: var(--muted); font-size: 13px; }
    .layout { display: grid; grid-template-columns: 260px 1fr; min-height: calc(100vh - 54px); }
    nav {
      border-right: 1px solid var(--line);
      padding: 1rem 0.8rem 2rem;
      background: var(--panel);
    }
    nav h2 { font-size: 11px; letter-spacing: .12em; text-transform: uppercase; color: var(--muted); margin: 0 0 .6rem .4rem; }
    nav a {
      display: block; color: var(--text); text-decoration: none;
      padding: .4rem .55rem; border-radius: 8px; margin-bottom: 2px;
    }
    nav a:hover, nav a.active { background: #212833; }
    nav .src { display:block; color: var(--muted); font-size: 12px; }
    main { padding: 1.5rem 2rem 4rem; max-width: 920px; }
    .hero h1 { margin: 0 0 .4rem; font-size: 1.8rem; }
    .hero p { color: var(--muted); margin: 0 0 1.2rem; }
    .cards { display: grid; gap: 1rem; }
    .card {
      background: var(--panel); border: 1px solid var(--line); border-radius: 14px;
      padding: 1rem 1.1rem; text-decoration: none; color: inherit; display: block;
    }
    .card:hover { border-color: var(--accent); }
    .badge {
      display: inline-block; min-width: 3.4rem; text-align: center;
      font: 700 11px/1.2 ui-monospace, SFMono-Regular, Consolas, monospace;
      padding: .22rem .4rem; border-radius: 6px; margin-right: .55rem;
    }
    .GET { background: #12351f; color: var(--get); }
    .POST { background: #152a4a; color: var(--post); }
    .PUT, .PATCH { background: #3a2e0c; color: var(--put); }
    .DELETE { background: #3a1518; color: var(--del); }
    .op {
      border: 1px solid var(--line); border-radius: 14px; background: var(--panel);
      margin: 1.1rem 0; overflow: hidden;
    }
    .op-h { padding: .85rem 1rem; border-bottom: 1px solid var(--line); }
    .path { font-family: ui-monospace, SFMono-Regular, Consolas, monospace; }
    .op-b { padding: 1rem; }
    .desc { color: var(--muted); }
    table { width: 100%; border-collapse: collapse; font-size: 13px; }
    th, td { text-align: left; padding: .35rem .4rem; border-bottom: 1px solid var(--line); }
    pre {
      background: #0b0e13; border: 1px solid var(--line); border-radius: 10px;
      padding: .8rem; overflow: auto; font-size: 12.5px;
    }
    textarea {
      width: 100%; min-height: 120px; background: #0b0e13; color: var(--text);
      border: 1px solid var(--line); border-radius: 10px; padding: .7rem; font: 12.5px ui-monospace, Consolas, monospace;
    }
    button {
      background: var(--accent); color: #06221e; border: 0; border-radius: 8px;
      font-weight: 700; padding: .5rem .9rem; cursor: pointer;
    }
    .try { margin-top: .8rem; }
    .note { font-size: 13px; color: var(--muted); }
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
          <a class="card" href="http://localhost:3011" target="_blank" rel="noreferrer">
            <strong>P11 CI dashboard</strong>
            <div class="note">GitHub Actions runs (local <code>pf-developer-ci-dash</code>)</div>
          </a>
          <a class="card" href="http://localhost:3013" target="_blank" rel="noreferrer">
            <strong>P11 PR review BFF</strong>
            <div class="note">GitHub PR API proxy (<code>pf-developer-review</code>)</div>
          </a>
          <a class="card" href="http://localhost:3010" target="_blank" rel="noreferrer">
            <strong>P11 repo scanner</strong>
            <div class="note">Static checks (<code>pf-developer-scanner</code>)</div>
          </a>
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
                <input data-try-path value="{{.Path}}" style="width:100%;margin:.3rem 0 .6rem;padding:.45rem .6rem;border-radius:8px;border:1px solid var(--line);background:#0b0e13;color:var(--text);font-family:ui-monospace,Consolas,monospace">
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
