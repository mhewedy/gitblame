package main

import "net/http"

const indexHtmlContent = `<script src="https://code.jquery.com/jquery-3.3.1.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js"></script>
<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/diff2html/2.11.3/diff2html.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/mustache.js/3.1.0/mustache.min.js"></script>

<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css">
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/diff2html/2.11.3/diff2html.min.css">
<style>
    @media (min-width: 768px) {
        .modal-xxl {
            width: 100%;
            max-width: 1200px;
        }
    }
</style>

<body style="margin: 20px 20px 20px 20px">
<!--Templates -->
<script id="main" type="text/template">
    <div class="dropdown">
        <button class="btn btn-secondary dropdown-toggle" type="button" id="dropdownMenuButton" data-toggle="dropdown"
                aria-haspopup="true" aria-expanded="false" style="margin-bottom: 20px">
            Team List with number of commits
        </button>
        <div class="dropdown-menu" aria-labelledby="dropdownMenuButton">
            {{#data.authors}}
            <a class="dropdown-item" href="#dd" onclick="showCommits('{{index}}')">
                {{author.name}}&lt;{{author.email}}&gt;
                <span class="badge badge-primary badge-pill">{{commits.length}}</span>
            </a>
            {{/data.authors}}
        </div>
    </div>

    {{#data.author}}
    <div class="alert alert-info" role="alert">
        User: {{data.author.name}}&lt;{{data.author.email}}&gt; has {{data.commits.length}} commit:
    </div>
    {{/data.author}}

    <div class="list-group">
        {{#data.commits}}
        <a href="#lg" class="list-group-item list-group-item-action flex-column align-items-start"
           onclick="showDiff('{{title}}','{{hash}}')" data-toggle="modal" data-target="#exampleModalLong">
            <div class="d-flex w-100 justify-content-between">
                <h6 class="mb-1">{{{title}}}</h6>
                <small>{{since}}</small>
            </div>
            <p class="mb-1">{{{body}}}</p>
        </a>
        {{/data.commits}}
    </div>

    <!-- Modal -->
    <div class="modal fade" id="exampleModalLong" tabindex="-1" role="dialog" aria-labelledby="exampleModalLongTitle"
         aria-hidden="true">
        <div class="modal-dialog modal-xxl" role="document">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title" id="exampleModalLongTitle">Modal title</h5>
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                        <span aria-hidden="true">&times;</span>
                    </button>
                </div>
                <div class="modal-body">
                    ...
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
                    <button type="button" class="btn btn-primary">Save changes</button>
                </div>
            </div>
        </div>
    </div>

</script>
<!-- //Templates -->
<div class="target-output"></div>
</body>

<script>
    function bind(elementId, json) {
        var template = $('#' + elementId).html();
        Mustache.parse(template);
        $(".target-output").html(Mustache.render(template, json));
    }
</script>
<script>
    var authors;
    var commits;
    var author;

    $(document).ready(function () {
        $.ajax({
            url: "/api"
        }).done(function (data) {
            authors = JSON.parse(data);
            authors.forEach(function (author, index) {
                author.index = index
            });
            bind("main", {
                "data": {"authors": authors}
            })
        });
    });

    function showCommits(index) {

        author = authors[index].author;
        commits = authors[index].commits.map(function (commit) {
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
        bind("main", {
            "data": {
                "authors": authors,
                "commits": commits,
                "author": author
            }
        })
    }

    function showDiff(title, hash) {
        $.ajax({
            url: "/api/diff/" + hash
        }).done(function (data) {
            var diffHtml = Diff2Html.getPrettyHtml(
                data, {inputFormat: 'diff', showFiles: true, matching: 'lines', outputFormat: 'line-by-line'}
            );

            bind("main", {
                "data": {
                    "authors": authors,
                    "commits": commits,
                    "author": author,
                    "diff": {
                        "title": title,
                        "body": diffHtml
                    }
                }
            })
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
