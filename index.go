package main

import "net/http"

const indexHtmlContent = `<script src="https://code.jquery.com/jquery-3.3.1.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js"></script>
<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/diff2html/2.11.3/diff2html.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/mustache.js/3.1.0/mustache.min.js"></script>
<script src="https://unpkg.com/accessible-nprogress/dist/accessible-nprogress.min.js"></script>

<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css">
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/diff2html/2.11.3/diff2html.min.css">
<link rel='stylesheet' href='https://unpkg.com/accessible-nprogress/dist/accessible-nprogress.min.css'/>

<style>
    @media (min-width: 768px) {
        .modal-xxl {
            width: 100%;
            max-width: 1200px;
        }
    }

    a {
        cursor: pointer;
    }
</style>

<body style="margin: 20px 20px 20px 20px">

<script id="error-template" type="text/template">
    {{#error}}
    <div class="alert alert-danger alert-dismissible fade show" role="alert">
        {{{error}}}
        <button type="button" class="close" data-dismiss="alert" aria-label="Close">
            <span aria-hidden="true">&times;</span>
        </button>
    </div>
    {{/error}}
</script>
<script id="info-template" type="text/template">
    {{#info}}
    <div class="alert alert-success alert-dismissible fade show" role="alert">
        {{{info}}}
        <button type="button" class="close" data-dismiss="alert" aria-label="Close">
            <span aria-hidden="true">&times;</span>
        </button>
    </div>
    {{/info}}
</script>

<script id="settings-template" type="text/template">
    <nav class="navbar navbar-light bg-light">
        <span class="navbar-brand mb-0 h1">Commits for Project at :
            <span style="font-family: monospace; font-weight: bold">{{settings.path}}</span>
        </span>
    </nav>
</script>

<script id="action-buttons-template" type="text/template">
    <div style="margin: 20px">
        <button id="btn-update" class="btn btn-warning" type="button" onclick="update()">
            {{#updating}}
            <span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
            {{/updating}}
            Sync from Repository
        </button>
    </div>
</script>

<script id="dropdown-template" type="text/template">
    <div class="dropdown">
        <button class="btn btn-secondary dropdown-toggle" type="button" id="dropdownMenuButton" data-toggle="dropdown"
                aria-haspopup="true" aria-expanded="false" style="margin: 20px">
            Team list with number of commits
        </button>
        <div class="dropdown-menu" aria-labelledby="dropdownMenuButton">
            {{#authors}}
            <a class="dropdown-item" onclick="showCommits('{{index}}')">
                {{author.name}}&lt;{{author.email}}&gt;
                <span class="badge badge-primary badge-pill">{{commits.length}}</span>
            </a>
            {{/authors}}
        </div>
    </div>
</script>

<script id="list-group-template" type="text/template">
    {{#author}}
    <div class="alert alert-info" role="alert">
        User: {{author.name}}&lt;{{author.email}}&gt; has {{commits.length}} commit:
    </div>
    {{/author}}

    <div class="list-group">
        {{#commits}}
        <a class="list-group-item list-group-item-action flex-column align-items-start"
           onclick="showDiff('{{title}}','{{hash}}')">
            <div class="d-flex w-100 justify-content-between">
                <h6 class="mb-1">{{{title}}}</h6>
                <small>{{since}}</small>
            </div>
            <p class="mb-1">{{{body}}}</p>
        </a>
        {{/commits}}
    </div>
</script>

<script id="diff-dialog-template" type="text/template">
    <div class="modal fade" id="diffModalLong" tabindex="-1" role="dialog" aria-labelledby="diffModalLongTitle"
         aria-hidden="true">
        <div class="modal-dialog modal-xxl" role="document">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title" id="diffModalLongTitle">{{{title}}}</h5>
                    {{#bbUrl}}
                    <a href="{{bbUrl}}" target="_blank">
                        <img style="padding: 5px 20px 0 20px;"
                             src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABgAAAAYCAYAAADgdz34AAAACXBIWXMAAAsTAAALEwEAmpwYAAABYUlEQVRIidWUP0oDURCHv3nZNFaKIIgsqVwEbeztrIRACg/hAexcC8HFQ+QgXkDxBLpKbEQUERW0szA7FnFl9/2JkKyg072Zeb9v5r1h4L+bACwf6o4IWw0rHw9S6UcAInSBbqP6igH65utw2aQ4QMFI0wAUhrxpAFIBiDYPKDUNQPvjFwDtUQdSOpJMb4HYk/sM3Ad0OsCsq87dIJUYICp9Crl4ACqcXafS86knmZ4AG06gMjTmGyr+ZxJlPVA9wKrPqRUtU/GH/iFeOdJ525lkugTMeW8Ung7MmEkqhm4XRlgL5UvL08H7mElS4wIKwgCDB3BzIK/Ag7ci3z9oEPB0tScv5SGygjmw6ACgl2R6brk7gWJqL1EHCDnKplssMwQmxiXU91p1itAGVkZhaRgrPjVAGNNBy3AxNWBYL1LshCTTR2BhQv23wb7UdpM9RYiwq7CNuvCfTOF0wsL+sH0C8bZgkDkx4mYAAAAASUVORK5CYII="/>
                    </a>
                    {{/bbUrl}}
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                        <span aria-hidden="true">&times;</span>
                    </button>
                </div>
                <div class="modal-body">
                    {{{body}}}
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
                </div>
            </div>
        </div>
    </div>
</script>

<div class="error"></div>
<div class="info"></div>
<div class="settings"></div>

<div>
    <div class="row">
        <div class="col-6">
            <div class="dropdown"></div>
        </div>
        <div class="col-6">
            <div class="action-buttons" style="float: right"></div>
        </div>
    </div>
</div>

<div class="list-group"></div>
<div class="diff-dialog"></div>

</body>

<script>
    function bind(elementId, json) {
        var template = $('#' + elementId + "-template").html();
        Mustache.parse(template);
        $("." + elementId).html(Mustache.render(template, json));
    }
</script>
<script>

    var authors;
    var repoUrl;

    $(document).ready(function () {

        $.ajaxSetup({
            beforeSend: function () {
                NProgress.start();
            },
            complete: function () {
                NProgress.done();
            },
            success: function () {
            }
        });

        $.ajax({
            url: "/api/settings"
        }).done(function (data) {
            var settings = JSON.parse(data);
            repoUrl = settings.path;
            bind("settings", {
                "settings": settings
            })
        });

        getAuthors();
    });

    function getAuthors() {
        $.ajax({
            url: "/api"
        }).done(function (data) {
            authors = JSON.parse(data);
            authors.forEach(function (author, index) {
                author.index = index
            });
            bind("dropdown", {
                "authors": authors
            })
        });
    }

    // ---------------------------------------

    bind("action-buttons", {});

    function update() {
        toggleButton("btn-update", "updating", false);
        $.ajax({url: "/api/update"}).done(function () {

            getAuthors();

            toggleButton("btn-update", "updating", true);
            bind("info", {"info": "Success"});

        }).fail(function (resp) {
            if (resp.status === 401) {
                bind("error", {
                    "error": "Cannot sync from repository, authentication required. For now, try to use <b>git pull</b>."
                })
            }
            toggleButton("btn-update", "updating", true);
        });
    }

    function toggleButton(btnId, modelVar, enable) {
        var obj = {};
        obj[modelVar] = !enable;
        bind("action-buttons", obj);
        $("#" + btnId).prop('disabled', !enable);
    }

    // ---------------------------------------

    function showCommits(index) {
        var author = authors[index].author;
        var commits = authors[index].commits.map(function (commit) {
            var messageArr = commit.message.split("\n");
            var title = messageArr[0];
            var body = messageArr.slice(1, messageArr.length).join('<br/>');
            return {
                "title": escapeQuotes(title),
                "body": escapeQuotes(body),
                "hash": commit.hash,
                "since": 'Since ' + timeSince(new Date(commit.when)) + ' ago'
            }
        });
        bind("list-group", {
            "commits": commits,
            "author": author
        })
    }

    // ---------------------------------------

    function showDiff(title, hash) {

        function isBitbucketUrl() {
            return repoUrl.indexOf('/scm/') > 0;
        }

        if (isBitbucketUrl()) {
            var parts = repoUrl.split('/');
            var bbUrl = parts.slice(0, 3).join('/') + '/projects/' + parts[4] + '/repos/' +
                parts[5].split('.')[0] + '/commits/' + hash;
        }

        $.ajax({
            url: "/api/diff/" + hash
        }).done(function (data) {
            var diffHtml = Diff2Html.getPrettyHtml(
                data, {inputFormat: 'diff', showFiles: true, matching: 'lines', outputFormat: 'line-by-line'}
            );
            bind("diff-dialog", {
                "title": title,
                "body": diffHtml,
                bbUrl: bbUrl
            });
            $('#diffModalLong').modal('show');
        });
    }

</script>

<script>

    function escapeQuotes(str) {
        return str.replace(/'/g, '&apos;').replace(/"/g, '&quot;');
    }

    //https://stackoverflow.com/a/3177838/171950
    function timeSince(date) {
        var seconds = Math.floor((new Date() - date) / 1000);
        var interval = Math.floor(seconds / 31536000);

        if (interval > 1) {
            return interval + " years";
        }
        interval = Math.floor(seconds / 2592000);
        if (interval > 1) {
            return interval + " months";
        }
        interval = Math.floor(seconds / 86400);
        if (interval > 1) {
            return interval + " days";
        }
        interval = Math.floor(seconds / 3600);
        if (interval > 1) {
            return interval + " hours";
        }
        interval = Math.floor(seconds / 60);
        if (interval > 1) {
            return interval + " minutes";
        }
        return Math.floor(seconds) + " seconds";
    }

</script>
`

func Index(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/html")
	writer.Write([]byte(indexHtmlContent))
}
