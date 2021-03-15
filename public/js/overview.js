const fetchNewContentAnchor = 0.8;
const limit = 10;
let offset = 0,
  prevOffset = -1;

const formatArticle = (id, title, subtitle, tags, category, content) => {
  let isMobile = /iPhone|iPad|iPod|Android/i.test(navigator.userAgent);

  if (isMobile && title.length > 55) {
    title = title.substr(0, 55) + " ...";
  }
  title = title.replaceAll(" ", "&nbsp;");

  if (isMobile) {
    subtitle = subtitle.replaceAll(" ", "&nbsp;");
  }

  let tagHTML = "";
  tags.forEach((t) => {
    tagHTML += `<a href="/articles/tags?query=${t}"><span class="tag is-warning">${t}</span></a>&nbsp;`;
  });

  category = `<a href="/articles/${category}"><span class="tag is-primary">${category}</span></a>`;

  let img = /<img.*>/.exec(content); // Note that there is a <p></p> tag surrounded
  let img_add_class = '<img class="article-list-img-h" ';
  if (img != null) {
    truncEnd = img[0].length - 4;
    if (img[0][img[0].length - 5] == "<") {
      // If the image is embedded in a list, then it will be surrounded by <li></li> instead of <p></p>
      truncEnd -= 1;
    }
    let img_tag = `<div class="column is-5" style="text-align: center;">` + img_add_class + img[0].substring(4, truncEnd) + "</div>"; // substring(): to remove "</p>"
    content =
      '<div class="column is-7">' +
      content
        .replaceAll(/<img.*>/g, "")
        .replaceAll(/<pre>/g, "")
        .replaceAll(/<\/pre>/g, "") +
      "</div>" +
      img_tag;
  } else {
    content =
      '<div class="column is-12">' +
      content
        .replaceAll(/<img.*>/g, "")
        .replaceAll(/<pre>/g, "")
        .replaceAll(/<\/pre>/g, "") +
      "</div>";
  }

  return `<div class="tile is-ancestor">
                <div class="tile is-parent">
                    <div class="tile is-child box article-list-container">
                        <div data-articleid=${id}></div>
                        <div class="article-list-tag">
                            ${tagHTML}
                            ${category}
                        </div>
                        <p class="title">${title}</p>
                        <p class="subtitle">${subtitle}</p>
                        <div class="columns overview-content">${content}</div>
                    </div>
                </div>
            </div>`;
};

const appendNewContent = (content) => {
  ele = document.createElement("div");
  ele.classList.add("articles-parent");
  ele.innerHTML = content;

  lastArticlesParents.insertAfter = newContent;
  lastArticlesParents.parentNode.insertBefore(ele, lastArticlesParents.nextSibling);
  lastArticlesParents = ele;
};

const checkStatus = async (resp) => {
  const contentType = resp.headers.get("content-type");

  if (contentType && contentType.indexOf("application/json") !== -1 && resp.status < 400) {
    return Promise.resolve(resp);
  }
  return Promise.reject(resp);
};

const fetchSucceed = async (resp) => {
  await resp.json().then((data) => {
    prevOffset = offset;
    offset += data.size;

    data.articleList = data.articleList || []; // If there is no data, the empty array returned by backend becomes null
    if (data.articleList.length == 0) {
      return;
    }

    newContent = "";
    data.articleList.forEach((a) => {
      newContent += formatArticle(a.ID, a.Title, a.Subtitle, a.Tags, a.Category, a.Content);
    });
    appendNewContent(newContent);
  });
  return Promise.resolve(true);
};

const fetchFailed = async (resp) => {
  await resp.json().then((data) => {
    c("Error: ", data.errHead, data.errBody);
  });
  return Promise.resolve(false);
};

const fetchOlderContent = async (count) => {
  if (offset == prevOffset) {
    return;
  }

  let urlPath = location.pathname.split("/");
  let path = urlPath[urlPath.length - 1];

  if (path == "weekly-update") {
    return;
  } else if (path == "tags") {
    type = "tag";
    query = new URLSearchParams(window.location.search).get("query");
  } else {
    type = "category";
    query = path; // Either "pharma" or "medication"
  }

  let para = "?" + new URLSearchParams({ type: type, query: query, offset: offset, limit: limit });
  let url = "fetch" + para;
  let res = await fetchData(url).then(checkStatus).then(fetchSucceed).catch(fetchFailed);

  if (!res) {
    if (count < 3) {
      fetchOlderContent(++count);
    } else {
      showErrMsg("Failed to Fetch Content", "Please reload the page and try again.");
    }
  }
};

const fetchInitialContent = async () => {
  if (offset == 0) {
    await fetchOlderContent(0);
    if (offset == 0) {
      showNoticeMsg("Oops ... ", "There is no articles.");
    }
  }
};

onDOMContentLoaded = (function () {
  lastArticlesParents = document.getElementsByClassName("articles-parent")[0];
  offset = Number(document.getElementById("articles-count").innerText) || 0; // offset == # means we'll skip # articles in next fetch

  fetchInitialContent();

  document.getElementById("articles-container").addEventListener("click", (e) => {
    window.location.href = "/articles/browse?articleId=" + e.target.closest("div.tile.is-child").children[0].dataset.articleid;
  });

  let lastFetch = 0,
    delay = 800;
  window.onscroll = () => {
    if (lastFetch < Date.now() - delay && window.innerHeight + window.pageYOffset >= document.body.offsetHeight * fetchNewContentAnchor) {
      lastFetch = Date.now();
      fetchOlderContent(0);
    }
  };
})();
