package web

const (
	index = `<!DOCTYPE html>
	<html>
	
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width,initial-scale=1">
		<title>curseforge mod 下载</title>
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/xmdhs/searchqanda/style.css">
	</head>
	
	<body>
		<div class="container-lg px-3 my-5 markdown-body">
			<h1>curseforge mod 下载</h1>
			<form action="/curseforge/s" target="_blank"><input type="text" name="q"></form><br>
			<hr>
			<details>
    			<summary>说明</summary>
    			<p>输入 mod 的英文名，回车即可</p>
    			<p>本站并非镜像站，只是利用 curseforge 提供的 api 来实现搜索和显示下载链接的。对于显示的内容不保证最新，下载的文件不保证完整性。</p>
			</details>
		</div>
	</body>

	</html>`
	searchhtml = `<!DOCTYPE html>
	<html>
	
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width,initial-scale=1">
		<title>curseforge mod 下载 - {{.Name}}</title>
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/xmdhs/searchqanda/style.css">
	</head>
	
	<body>
		<div class="container-lg px-3 my-5 markdown-body">
			<h1>{{.Name}}</h1>{{range .List}}<h3><a href="{{ .Link}}" target="_blank">{{ .Title}}</a></h3>
			<blockquote>
				<p>{{ .Txt}}</p>
			</blockquote>{{end}} {{if .T}}<br>
			<hr>
			<p><a href="{{ .Link}}">more</a></p>{{else}} {{end}}
		</div>
		<script src="https://cdn.jsdelivr.net/npm/anchor-js/anchor.min.js"></script>
		<script>anchors.add();</script>
	</body>
	
	</html>`
)
