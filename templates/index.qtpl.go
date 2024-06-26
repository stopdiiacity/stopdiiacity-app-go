// Code generated by qtc from "index.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

package templates

import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

func StreamGenerate(qw422016 *qt422016.Writer, prefixGroups []PrefixGroup) {
	qw422016.N().S(`
<!DOCTYPE html>
<html lang="en">
<head>
    <title>Stop Diia City</title>
    <meta charset="UTF-8">
    <meta name="description" content="Stop Diia City">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png">
    <link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
    <link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">
    <link rel="manifest" href="/site.webmanifest">
    <link type="text/css" rel="stylesheet" href="styles.css"/>
    <meta name="google-site-verification" content="jHJ25SxN1Yu6iMLtF2Psv_-QzNmp0mS4DRnmM5z6Glw"/>

    <!-- Latest compiled and minified CSS -->
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/3.4.1/css/bootstrap.min.css" integrity="sha384-HSMxcRTRxnN+Bdg0JdbxYKrThecOKuH5zCYotlSAcp1+c8xmyTe9GYg1l9a69psu" crossorigin="anonymous">

    <!-- Optional theme -->
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/3.4.1/css/bootstrap-theme.min.css" integrity="sha384-6pzBo3FDv/PJ8r2KRkGHifhEocL+1X2rVCTTkUfGk7/0pbek5mMa1upzvWbrUbOZ" crossorigin="anonymous">

    <!-- Latest compiled and minified JavaScript -->
    <script src="https://stackpath.bootstrapcdn.com/bootstrap/3.4.1/js/bootstrap.min.js" integrity="sha384-aJ21OjlMXNL5UyIl/XNwTMqvzeRMZH2w8c5cRVpzpU8Y5bApTppSuUkhZXN0VxHd" crossorigin="anonymous"></script>
</head>
<body>
<div class="container">
    <div class="page-header">
        <h1>Stop Diia City <iframe src="https://ghbtns.com/github-btn.html?user=stopdiiacity&repo=stopdiiacity-app-go&type=star&count=true&size=large" frameborder="0" scrolling="0" width="170" height="30" title="GitHub"></iframe></h1>
    </div>
    <div class="row">
        <div class="col-md-4">
            <h2>
                <a href="https://dou.ua/forums/topic/33673/" target="_blank">Обговорення на DOU</a>
            </h2>
            <ul class="list-group">
                <li class="list-group-item"><a href="https://addons.mozilla.org/en-GB/firefox/addon/stopdiiacity/" target="_blank">Завантажити розширення StopDiiaCity для Firefox з офіційного сайту</a></li>
                <li class="list-group-item"><a href="https://chrome.google.com/webstore/detail/stopdiiacity/omhhpgmnkkhepkpifbgbbjenlibgmolo" target="_blank">Завантажити розширення StopDiiaCity для Chrome з офіційного сайту</a></li>
                <li class="list-group-item"><a href="/swagger/index.html" target="_blank">Swagger documentation</a></li>
            </ul>
            <form id="js-form" class="form-group" action="javascript:void(0);">
                <fieldset>
                    <legend>Пошук за лінком</legend>
                    <label>
                        <input class="form-control" name="url" type="url" placeholder="company's url"/>
                    </label>
                    <input class="btn default" name="reset" type="reset" value="Reset"/>
                    <input class="btn btn-primary" name="submit" type="submit" value="Verify"/>
                </fieldset>
            </form>
        </div>
        <div class="col-md-8">
            <h2>Companies</h2>
            `)
	qw422016.N().S(`<ul class="list-group">`)
	for _, prefixGroup := range prefixGroups {
		for _, url := range prefixGroup.Prefixes {
			qw422016.N().S(`<li class="list-group-item"><a href="`)
			qw422016.E().S(url)
			qw422016.N().S(`">`)
			qw422016.E().S(url)
			qw422016.N().S(`</a></li>`)
		}
	}
	qw422016.N().S(`</ul>`)
	qw422016.N().S(`
        </div>
    </div>
    <div class="row">
        <div class="col-md-12">
            <a href="https://www.vultr.com/?ref=8741375"><img src="/vultr_banner_728x90.png" width="728" height="90"></a>
        </div>
    </div>
</div>

<script type="application/javascript">
    {
        const $form = document.getElementById("js-form");

        const $url = $form.elements["url"];

        $form.onsubmit = function () {
            const url = $url.value;

            fetch("/verify.json", {
                method: "POST",
                body: JSON.stringify({
                    "urls": [url],
                }),
            }).then(function (response) {
                return response.json();
            }).then(function (json) {
                if (json.exists) {
                    alert(`)
	qw422016.N().S("`")
	qw422016.N().S(`URL ${url} is StopDiiaCity!`)
	qw422016.N().S("`")
	qw422016.N().S(`);
                } else {
                    alert(`)
	qw422016.N().S("`")
	qw422016.N().S(`URL ${url} is safe!`)
	qw422016.N().S("`")
	qw422016.N().S(`);
                }
            }).catch(console.error)
        }
    }
</script>
</body>
</html>
`)
}

func WriteGenerate(qq422016 qtio422016.Writer, prefixGroups []PrefixGroup) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	StreamGenerate(qw422016, prefixGroups)
	qt422016.ReleaseWriter(qw422016)
}

func Generate(prefixGroups []PrefixGroup) string {
	qb422016 := qt422016.AcquireByteBuffer()
	WriteGenerate(qb422016, prefixGroups)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}
