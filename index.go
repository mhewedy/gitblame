package main

import "net/http"

const indexHtmlContent = `<script src="https://code.jquery.com/jquery-3.3.1.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js"></script>
<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/diff2html/2.11.3/diff2html.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/ractive/1.3.8/ractive.min.js"></script>
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

<script id="template" type="text/ractive">

    {{#error}}
    <div class="alert alert-danger alert-dismissible fade show" role="alert">
        {{{error}}}
        <button type="button" class="close" data-dismiss="alert" aria-label="Close">
            <span aria-hidden="true">&times;</span>
        </button>
    </div>
    {{/error}}

    {{#info}}
    <div class="alert alert-success alert-dismissible fade show" role="alert">
        {{{info}}}
        <button type="button" class="close" data-dismiss="alert" aria-label="Close">
            <span aria-hidden="true">&times;</span>
        </button>
    </div>
    {{/info}}

    <nav class="navbar navbar-light bg-light">
        <span class="navbar-brand mb-0 h1">Commits for Project at :
            <span style="font-family: monospace; font-weight: bold">{{settings.path}}</span>
        </span>
    </nav>

    <div>
        <div class="row">
            <div class="col-6">
                <div class="dropdown">
                    <button class="btn btn-secondary dropdown-toggle" type="button" id="dropdownMenuButton"
                            data-toggle="dropdown"
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
            </div>
            <div class="col-6">
                <div style="margin: 20px; float: right">
                    <button id="btn-update" class="btn btn-warning" type="button" onclick="update()">
                        {{#updating}}
                        <span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
                        {{/updating}}
                        Sync from Repository
                    </button>
                </div>
            </div>
        </div>
    </div>


    {{#author}}
    <h6>
        <div class="alert alert-secondary" role="alert">
            User: {{author.name}}&lt;{{author.email}}&gt; has {{commits.length}} commit:
        </div>
    </h6>
    {{/author}}

    <div class="list-group-flush">
        {{#commits}}
        <a class="list-group-item list-group-item-action flex-column align-items-start"
           onclick="showDiff('{{title}}','{{hash}}')" id="commit-{{hash}}">
            <div class="d-flex w-100 justify-content-between">
                <h6 class="mb-1">{{{title}}}</h6>
                <small>{{since}}</small>
            </div>
            <p class="mb-1">{{{body}}}</p>
        </a>
        {{/commits}}
    </div>

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

<div id="target"></div>

</body>

<script>

    var ractive = new Ractive({
        target: '#target',
        template: '#template',
        data: {greeting: 'Hello', name: 'world'}
    });

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
            ractive.set("settings", settings);
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
            ractive.set("authors", authors);
            ractive.set("commits", []);
            ractive.set("author", {});
        });
    }

    // ---------------------------------------

    function update() {
        toggleButton("btn-update", "updating", false);
        $.ajax({url: "/api/update"}).done(function () {

            getAuthors();

            toggleButton("btn-update", "updating", true);
            ractive.set("info", "Success");

        }).fail(function (resp) {
            if (resp.status === 401) {
                ractive.set("error", "Cannot sync from repository, authentication required. For now, try to use <b>git pull</b>.");
            }
            if (resp.status === 404) {
                ractive.set("info", "Already up to date");
            }

            toggleButton("btn-update", "updating", true);
        });
    }

    function toggleButton(btnId, modelVar, enable) {
        ractive.set(modelVar, !enable);
        $("#" + btnId).prop('disabled', !enable);
    }

    // ---------------------------------------

    function showCommits(index) {
        var author = authors[index].author;
        var commits = authors[index].commits.map(function (commit) {
            var messageArr = escapeHtml(commit.message).split("\n");
            var title = messageArr[0];
            var body = messageArr.slice(1, messageArr.length).join('<br/>');
            return {
                "title": title,
                "body": body,
                "hash": commit.hash,
                "since": 'Since ' + timeSince(new Date(commit.when)) + ' ago'
            }
        });

        ractive.set("commits", commits);
        ractive.set("author", author);
    }

    // ---------------------------------------

    function showDiff(title, hash) {

        function isBitbucketUrl() {
            return repoUrl.indexOf('/scm/') > 0;
        }

        function applyActiveClass() {
            $("#commit-" + hash).addClass('active').css("color", "#fff");
        }

        if (isBitbucketUrl()) {
            var parts = repoUrl.split('/');
            var bbUrl = parts.slice(0, 3).join('/') + '/projects/' + parts[4] + '/repos/' +
                parts[5].split('.')[0] + '/commits/' + hash;
        }

        applyActiveClass();

        $.ajax({
            url: "/api/diff/" + hash
        }).done(function (data) {
            var diffHtml = Diff2Html.getPrettyHtml(
                data, {inputFormat: 'diff', showFiles: true, matching: 'lines', outputFormat: 'line-by-line'}
            );

            ractive.set("title", title);
            ractive.set("body", diffHtml);
            ractive.set("bbUrl", bbUrl);

            $('#diffModalLong').modal('show');
            $('#diffModalLong').on('hidden.bs.modal', function (e) {
                $("#commit-" + hash).removeClass('active').css("color", "#000");
            })
        });
    }

</script>

<script>

    //https://github.com/janl/mustache.js/blob/master/mustache.js
    var entityMap = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#39;',
        '/': '&#x2F;',
        '=': '&#x3D;'
    };

    function escapeHtml(string) {
        return String(string).replace(/[&<>"'=\/]/g, function fromEntityMap(s) {
            return entityMap[s];
        });
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
