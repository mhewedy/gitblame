package main

import "net/http"

const indexHtmlContent = `
<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css"
      crossorigin="anonymous">
<link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/diff2html/2.11.3/diff2html.min.css">

<script src="https://code.jquery.com/jquery-3.3.1.min.js" crossorigin="anonymous"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js"
        crossorigin="anonymous"></script>
<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js" crossorigin="anonymous"></script>

<script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/diff2html/2.11.3/diff2html.min.js"></script>
<style>
    @media (min-width: 768px) {
        .modal-xxl {
            width: 100%;
            max-width: 1200px;
        }
    }
	hr {
	  margin-top: 3px;
	  margin-bottom: 3px;
	  border: 0;
	  border-top: 1px solid rgba(0, 0, 0, 0.1);
	}
</style>

<script>
    $(document).ready(function () {
        $.ajax({
            url: "/api"
        }).done(function (data) {
            displayAuthors(JSON.parse(data));
        });
    });

    function displayAuthors(authorsWithCommits) {
        let i = 0;
        for (let authorWithCommits of authorsWithCommits) {
            $("#container").append(authorsCollapseHtml(i, authorWithCommits.author, authorWithCommits.commits));
            i++;
        }
    }

    function authorsCollapseHtml(index, author, commits) {
        let cellColor = (index % 2) === 0 ? '#f2f2f2' : '#ffffff';
        return '<div style="color: #3572b0; background: ' + cellColor + ';padding: 10px 0 10px 10px;" ' +
            ' >' +
            '<a aria-expanded="false" style="color: #3572b0; cursor: pointer;" data-toggle="collapse" data-target="#collapse' + index + '" >' +
            author.name + ' &lt;' + author.email + '&gt; ' + '(' + commits.length + ')' +
            '</a>' +
            '</div>' +
            '<div class="collapse" id="collapse' + index + '">' +
            '    <div class="card card-body">' +
            commitsHtml(commits) +
            '    </div>' +
            '</div>';
    }

    function commitsHtml(commits) {
        let response = '';
        let index = 0;
        for (commit of commits) {
            response += '<div style="margin-bottom: 5px;" >' +
                '<a data-toggle="modal" data-target="#diffModal" onclick="hashClicked(\'' + commit.hash + '\')" ' + 
				'style="color: #3572b0; text-align: left; cursor: pointer;">' +
                commit.message.split('\n')[0] +
                '</a>' +
				'<span style="margin-left: 20px; color: gray; font-size: small; float: right">(' + commit.when + ')</span>' +
                '</div>' +
				'<hr />';
            index++;
        }
        return response;
    }

    function hashClicked(hash) {
        $(".modal-body").empty();
        $.ajax({
            url: "/api/diff/" + hash
        }).done(function (data) {
            openDialog(data);
        });
    }

    function openDialog(diff) {
        var diffHtml = Diff2Html.getPrettyHtml(
            diff,
            {inputFormat: 'diff', showFiles: true, matching: 'lines', outputFormat: 'line-by-line'}
        );
        $(".modal-body").append(diffHtml);
    }

</script>

<div id="container">
</div>

<div class="modal" id="diffModal" tabindex="-1" role="dialog" aria-labelledby="diffModalTitle"
     aria-hidden="true">
    <div class="modal-dialog modal-xxl" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                    <span aria-hidden="true">&times;</span>
                </button>
            </div>
            <div class="modal-body" style="overflow: scroll;">
            </div>
        </div>
    </div>
</div>`

func Index(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/html")
	writer.Write([]byte(indexHtmlContent))
}
