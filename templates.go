package vcd

import (
	"crypto/md5"
	"encoding/hex"
	"html/template"
	"strings"
	"time"
)

var (
	listTmpl, newTmpl, showTmpl *template.Template
)

type newTmplArgs struct {
	Error string
}

func init() {
	baseTmpl := template.Must(template.New("base").Funcs(template.FuncMap{
		"eightiesTime": eightiesTime,
		"hashVessel":   hashVessel,
		"now":          time.Now,
		"removeVesselBlankLines": removeVesselBlankLines,
	}).Parse(`
{{define "layout" -}}
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>{{block "title" .}}{{end}}[VCD ARCHIVE]</title>
	<link rel="icon" type="image/vnd.microsoft.icon" href="http://vcd.phuce.com/favicon.ico"/>
	<link rel="alternate" type="application/rss+xml" href="http://vcd.phuce.com/rss.xml" title="VCD ARCHIVE LOG"/>
	<link rel="stylesheet" type="text/css" href="/styles.css"/>
</head>
<body>
	<div class="content">
		<div class="tty_h1 header invert">
			<div class="left">[<a href="/">VESSEL&nbsp;CONFIGURATION&nbsp;DATA&nbsp;ARCHIVE</a></div>
			<div class="right">{{now | eightiesTime}}]</div>
			<div style="clear: both;"></div>
		</div>
		{{block "content" .}}{{end}}
		<div class="tty_h1 footer invert">
			<div class="left">[</div>
			<div class="right">]</div>
			<div class="center"><a href="http://www.captainforever.com/">WWW.CAPTAINFOREVER.COM</a></div>
			<div style="clear: both;"></div>
		</div>
		<div>&nbsp;</div>
	</div>
	<script type="text/javascript">
		var gaJsHost = (("https:" == document.location.protocol) ? "https://ssl." : "http://www.");
		document.write(unescape("%3Cscript src='" + gaJsHost + "google-analytics.com/ga.js' type='text/javascript'%3E%3C/script%3E"));
	</script>
	<script type="text/javascript">
		try {
			var pageTracker = _gat._getTracker("UA-9140800-2");
			pageTracker._trackPageview();
		} catch(err) {}
	</script>
</body>
</html>
{{end}}

{{define "vessel"}}
<div class="vessel">
    <div class="tty_h2"><span>ARCHIVE DATA FILE <a href="/vessels/{{.ID}}">{{.ID}}.VCD</a></span></div>
    <div class="stats">
      <div class="logged">logged by pilot {{. | hashVessel}} on {{.CreatedAt | eightiesTime}}</div>
      <div class="piloted">{{.PilotClicks}} known pilot(s)</div>
      <div style="clear:both;"></div>
   </div>
   <br/>
   <div class="data">{{.Data | removeVesselBlankLines}}</div>
   <div class="link">[<a href="/vessels/{{.ID}}/pilot">PILOT THIS VESSEL</a>]</div>
</div>
{{end}}
	`))

	listTmpl = template.Must(template.Must(baseTmpl.Clone()).Parse(`
{{define "content"}}
<div class="tty_h1">
    <div class="left">[<span class="greenblock">DOWNLOADED {{. | len}} VESSEL{{if ne 1 (. | len)}}S{{end}}</span></div>
    <div class="right">[<a href="/new">UPLOAD NEW VESSEL</a>]</div>
    <div style="clear: both;"></div>
</div>
<br/>
{{range .}}
    {{template "vessel" .}}
    <br/>
{{end}}
{{end}}
	`))

	newTmpl = template.Must(template.Must(baseTmpl.Clone()).Parse(`
{{define "content"}}
<div class="tty_h1">
    <div class="left">[<span class="greenblock">ALLOCATING NEW ARCHIVE ENTRY...</span></div>
    <div class="right">[<a href="/">INDEX</a>][<a href="/new">UPLOAD NEW VESSEL</a>]</div>
    <div style="clear: both;"></div>
</div>
<br/>
<div class="tty_dialog" style="width: 812px; height: 480px;">
    <div class="borders">
        <div class="edge t"></div>
        <div class="edge r"></div>
        <div class="edge b"></div>
        <div class="edge l"></div>
        <div class="corner tl"></div><div class="corner tr"></div>
        <div class="corner bl"></div><div class="corner br"></div>
    </div>
    <div class="text" style="top: 28px;">
        INPUT VESSEL CONFIGURATION<br/>
        Enter HTML export of vehicle data for storage:<br/>
        {{if .Error}}<span class="redblock">{{.Error}}</span>{{end}}<br/>
        <form method="post">
            <textarea id="vessel_data" name="vessel_data"></textarea><br/>
            <br/>
            <input class="button" type="submit" value="TRANSMIT">
        </form>
    </div>
</div>
{{end}}
	`))

	showTmpl = template.Must(template.Must(baseTmpl.Clone()).Parse(`
{{define "content"}}
<div class="tty_h1">
    <div class="left">[<span class="greenblock">DOWNLOADED VESSEL</span></div>
    <div class="right">[<a href="/">INDEX</a>][<a href="/new">UPLOAD NEW VESSEL</a>]</div>
    <div style="clear: both;"></div>
</div>
<br/>
{{template "vessel" .}}
{{end}}
	`))
}

func eightiesTime(t time.Time) template.HTML {
	if t.IsZero() {
		t = time.Now()
	}
	t = t.AddDate(-20, 0, 0)
	return template.HTML(t.Format("2006-01-02&nbsp;15:04:05"))
}

func hashVessel(v *Vessel) (string, error) {
	h := md5.New()
	_, err := h.Write(v.IP)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func now() time.Time {
	return time.Now()
}

func removeVesselBlankLines(data string) template.HTML {
	var lines []string
	for _, line := range strings.Split(data, "<br/>") {
		if len(strings.Replace(line, "&nbsp;", "", -1)) == 0 {
			continue
		}
		lines = append(lines, line)
	}
	return template.HTML(strings.Join(lines, "<br/>"))
}
