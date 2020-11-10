package main

const tmpl = `
<html>

<head>
    <meta name="viewport" content="width=device-width,initial-scale=1"/>
    <title>LanBucket</title>
</head>

<body style="background-color: #3E606F; text-align: center;">
    <h1 style="color: #FCFFF5;">LanBucket</h1>
    <div style="border-radius: 10px; padding: 10px; margin: 0 auto; max-width: 900px; min-height: 450px; background: #FCFFF5;">
        <table border="0" width="100%">
            <tr>
                <td width="80%">文件</td>
                <td width="20%">大小</td>
            </tr>
            {{ range $k, $v := .files }}
            <tr>
                <td><a href="/file?name={{ $v.Name }}">{{ $v.Name }}</a></td>
                <td>{{ $v.Size }}B</td>
            </tr>
            {{ end }}
        </table>
    </div>
    <p style="color: #FCFFF5;">Copyright © 2020 <strong>ZX</strong></p>
</body>

</html>
`
