var ractive = new Ractive({
    target: '#target',
    template: '#template',
    data: {greeting: 'Hello', name: 'world'}
});

var authors;
var repoUrl;
var stats;
var currentAuthorIndex;

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

    getSettings();
    getAuthors();
    getStats();
});

function getSettings() {
    $.ajax({
        url: "/api/settings"
    }).done(function (data) {
        var settings = JSON.parse(data);
        repoUrl = settings.path;
        ractive.set("settings", settings);
    });
}

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

function getStats() {
    $.ajax({
        url: "/api/stats"
    }).done(function (data) {
        stats = JSON.parse(data);
        currentAuthorIndex && showCommits(currentAuthorIndex);
    });
}

// ---------------------------------------

function update() {
    toggleButton("btn-update", "updating", false);

    $.ajax({url: "/api/update"}).done(function () {
        getAuthors();

        success("Success");
        toggleButton("btn-update", "updating", true);
    }).fail(function (resp) {
        if (resp.status === 401) {
            error("Cannot sync from repository, authentication required. For now, try to use <b>git pull</b>.")
        } else if (resp.status === 404) {
            success("Already up to date");
        } else {
            error("Internal server error, check logs.")
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
    currentAuthorIndex = index;
    var author = authors[index].author;
    var commits = authors[index].commits.map(function (commit) {
        var messageArr = escapeHtml(commit.message).split("\n");
        var title = messageArr[0];
        var body = messageArr.slice(1, messageArr.length).join('<br/>');
        return {
            "title": title,
            "body": body,
            "hash": commit.hash,
            "since": 'Since ' + timeSince(new Date(commit.when)) + ' ago',
            "stat": stats && stats[commit.hash]
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

// ---------------

function success(msg) {
    ractive.set("info", msg);
    ractive.set("type", "success");
    showAlert();
}

function error(msg) {
    ractive.set("info", msg);
    ractive.set("type", "danger");
    showAlert();
}

function hideAlert() {
    $(".alert").css('display', 'none');
}

function showAlert() {
    $(".alert").css('display', 'block');
}
