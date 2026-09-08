package decoy

import (
"html/template"
"net/http"
)

var decoyTemplate = template.Must(template.New("decoy").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>Maintenance</title>
<style>
body { font-family: sans-serif; background-color: #f4f4f4; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; }
.card { background: white; padding: 40px; border-radius: 10px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center; }
h1 { color: #333; }
</style>
</head>
<body>
<div class="card">
<h1>System Under Maintenance</h1>
<p>Please check back later.</p>
</div>
</body>
</html>`))

func ServeDecoy(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "text/html")
decoyTemplate.Execute(w, nil)
}
