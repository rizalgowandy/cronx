package page

import (
	"html/template"
	"sync"
)

const jobsTemplate = `
<!--
	This page is only being used for development to restructure the code,
	the real html page is on jobs.go.
-->
<!DOCTYPE html>
<html lang="en">
<head>
	<!-- Standard Meta -->
	<meta charset="UTF-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
	<meta http-equiv="refresh" content="30"/>
	<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0">
	<!-- Site Properties -->
	<title>Cronx</title>
	<link
	   rel="stylesheet"
	   type="text/css"
	   href="https://cdn.jsdelivr.net/npm/semantic-ui@2.4.2/dist/semantic.min.css">
	<script
	   src="https://code.jquery.com/jquery-3.1.1.min.js"
	   integrity="sha256-hVVnYaiADRTO2PzUGmuLJr8BLUSjGIZsDYGmIJLv2b8="
	   crossorigin="anonymous"></script>
	<script
	   src="https://cdn.jsdelivr.net/npm/semantic-ui@2.4.2/dist/semantic.min.js"
	   crossorigin="anonymous"></script>
	<script
	   src="https://cdnjs.cloudflare.com/ajax/libs/html2canvas/0.5.0-beta4/html2canvas.min.js"
	   integrity="sha512-OqcrADJLG261FZjar4Z6c4CfLqd861A3yPNMb+vRQ2JwzFT49WT4lozrh3bcKxHxtDTgNiqgYbEUStzvZQRfgQ=="
	   crossorigin="anonymous"></script>
	<script src="https://cdn.jsdelivr.net/npm/canvas2image@1.0.5/canvas2image.min.js"></script>
	<script type='text/javascript'>
		function screenshot() {
			html2canvas(document.querySelector('#data_table')).then(function(canvas) {
				Canvas2Image.saveAsPNG(canvas, canvas.width, canvas.height);
			});
		}
	</script>
	<style>
        body > .ui.container {
            margin-top: 3em;
            padding-bottom: 3em;
        }
	</style>
	<title>Cronx</title>
</head>
<body>
<div class="ui container">
	<div class="ui left fixed vertical stackable inverted main menu">
		<div class="header item">
			<i class="stopwatch icon"></i>
			Cronx
		</div>
		<a class="item active" href="javascript:window.location.reload()">
			<i class="tasks icon"></i>
			Jobs
		</a>
		<a class="item" href="/histories">
			<i class="history icon"></i>
			Histories
		</a>
		<div class="item" onclick="screenshot()">
			<button class="fluid ui labeled inverted green icon button">
				<i class="camera icon"></i>
				<div class="left aligned">Screenshot</div>
			</button>
		</div>
	</div>
	<div class="ui five steps">
		<div class="step">
			<i class="arrow down icon"></i>
			<div class="content">
				<div class="title">Down</div>
				<div class="description">Job fails to be registered</div>
			</div>
		</div>
		<div class="step">
			<i class="arrow up icon"></i>
			<div class="content">
				<div class="title">Up</div>
				<div class="description">Job has just been created</div>
			</div>
		</div>
		<div class="step">
			<i class="sync icon"></i>
			<div class="content">
				<div class="title">Running</div>
				<div class="description">Job is currently running</div>
			</div>
		</div>
		<div class="step">
			<i class="hourglass end icon"></i>
			<div class="content">
				<div class="title">Idle</div>
				<div class="description">Job is waiting for next execution time</div>
			</div>
		</div>
		<div class="step">
			<i class="attention icon"></i>
			<div class="content">
				<div class="title">Error</div>
				<div class="description">Job fails on the prev run</div>
			</div>
		</div>
	</div>
	<div id="data_table">
		<table class="ui sortable selectable center aligned celled table">
			<thead>
			<tr>
				<th id="id"
                        {{if eq (index .Sort.Columns "id") "ASC"}} class="sorted ascending"
                        {{else if eq (index .Sort.Columns "id") "DESC"}} class="sorted descending"
                        {{end}}
				>ID
				</th>
				<th id="name"
                        {{if eq (index .Sort.Columns "name") "ASC"}} class="sorted ascending"
                        {{else if eq (index .Sort.Columns "name") "DESC"}} class="sorted descending"
                        {{end}}
				>Name
				</th>
				<th id="status"
                        {{if eq (index .Sort.Columns "status") "ASC"}} class="sorted ascending"
                        {{else if eq (index .Sort.Columns "status") "DESC"}} class="sorted descending"
                        {{end}}
				>Status
				</th>
				<th id="prev_run"
                        {{if eq (index .Sort.Columns "prev_run") "ASC"}} class="sorted ascending"
                        {{else if eq (index .Sort.Columns "prev_run") "DESC"}} class="sorted descending"
                        {{end}}
				>Prev run
				</th>
				<th id="next_run"
                        {{if eq (index .Sort.Columns "next_run") "ASC"}} class="sorted ascending"
                        {{else if eq (index .Sort.Columns "next_run") "DESC"}} class="sorted descending"
                        {{end}}
				>Next run
				</th>
				<th id="latency"
                        {{if eq (index .Sort.Columns "latency") "ASC"}} class="sorted ascending"
                        {{else if eq (index .Sort.Columns "latency") "DESC"}} class="sorted descending"
                        {{end}}
				>Latency
				</th>
			</tr>
			</thead>
			<tbody>
            {{range .Data}}
				<tr
                        {{if eq .Job.Status "RUNNING"}} class="warning"
                        {{else if eq .Job.Status "IDLE"}} class="positive"
                        {{else if eq .Job.Status "DOWN"}} class="error"
                        {{else if eq .Job.Status "ERROR"}} class="error"
                        {{end}}
				>
					<td>{{.ID}}</td>
					<td class="left aligned">
                        {{if gt .Job.TotalWave 1 }}
                            {{.Job.Name}} ({{.Job.Wave}}/{{.Job.TotalWave}})
                        {{else}}
                            {{.Job.Name}}
                        {{end}}
					</td>
					<td>
                        {{if eq .Job.Status "RUNNING"}}
							<div class="ui yellow label">
								<i class="sync icon"></i>
                                {{.Job.Status}}
							</div>
                        {{else if eq .Job.Status "IDLE"}}
							<div class="ui green label">
								<i class="hourglass end icon"></i>
                                {{.Job.Status}}
							</div>
                        {{else if eq .Job.Status "DOWN"}}
							<div class="ui red label">
								<i class="arrow down icon"></i>
                                {{.Job.Status}}
							</div>
                        {{else if eq .Job.Status "ERROR"}}
							<div class="ui red label">
								<i class="attention icon"></i>
                                {{.Job.Status}}
							</div>
                        {{else}}
							<div class="ui label">
								<i class="arrow up icon"></i>
                                {{.Job.Status}}
							</div>
                        {{end}}
					</td>
					<td>
                        {{if eq .Job.Status "ERROR"}}
                            {{if not .Prev.IsZero}}
                                {{.Prev.Format "2006-01-02 15:04:05"}}
                            {{end}}
							<br/>
                            {{.Job.Error}}
                        {{else}}
                            {{if not .Prev.IsZero}}
                                {{.Prev.Format "2006-01-02 15:04:05"}}
                            {{end}}
                        {{end}}
					</td>
					<td>
                        {{if eq .Job.Status "DOWN"}}
                            {{.Job.Error}}
                        {{else}}
                            {{if not .Next.IsZero}}
                                {{.Next.Format "2006-01-02 15:04:05"}}
                            {{end}}
                        {{end}}
					</td>
					<td>{{.Job.Latency}}</td>
				</tr>
            {{end}}
			</tbody>
		</table>
	</div>
</div>
</body>
</html>
`

var (
	jobsPageOnce  sync.Once
	jobsPage      *template.Template
	jobsPageError error
)

func GetJobsPageTemplate() (*template.Template, error) {
	jobsPageOnce.Do(func() {
		t := template.New(jobsTemplateName)
		jobsPage, jobsPageError = t.Parse(jobsTemplate)
	})

	return jobsPage, jobsPageError
}
