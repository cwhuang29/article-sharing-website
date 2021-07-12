/*
 * Note: about 95% of code is the same as overview.js, without the sessionStorage feature
 */

const bookmarkEndpoint = '/articles/bookmark';
const fetchNewContentAnchor = 0.8;
const limit = 10;
let offset = 0;
let stopFetching = false;

const appendNewContent = (content) => {
  ele = document.createElement('div');
  ele.innerHTML = content;

  let anchor = articlesContainer.lastElementChild;
  anchor.parentNode.insertBefore(ele, anchor.nextSibling);
};

const isMobile = () => {
  return /iPhone|iPad|iPod|Android/i.test(navigator.userAgent);
};

const toTitleCase = (s) => {
  if (typeof s == 'string' && s.length > 0) {
    return s[0].toUpperCase() + s.substr(1);
  }
  return '';
};

const encodeHTMLEntities = (val) => {
  /*
   * Input: <scrpit>console.log(1)</script>
   * Output: &lt;script&gt;console.log(1)&lt;/script&gt;
   */
  let e = document.createElement('textarea');
  e.innerHTML = val;
  return e.innerHTML;
};

const formatArticle = (article) => {
  let { id, adminOnly, title, subtitle, tags, category, outline, cover_photo } = article;

  const titleTag = `<p class="title">${title}</p>`;

  let subtitleTag = ''; // Don't show subtitle on mobile devices
  if (!isMobile()) {
    subtitleTag = `<p style="font-size: 110%; font-weight: 600; margin-bottom: 4.5px">${subtitle}</p>`;
  }

  let tagHTML = '';
  tags.forEach((t) => {
    tagHTML += `<a href="/articles/tags?q=${encodeURIComponent(t)}"><span class="tag is-warning">${encodeHTMLEntities(t)}</span></a>`;
  });

  categoryTag = `<a href="/articles/${category}"><span class="tag is-primary">${toTitleCase(category)}</span></a>`;

  adminTag = '';
  if (adminOnly) {
    adminTag = `<span class="tag is-danger">Admin Only</span>`;
  }

  let overviewContent;
  const outlineTag = `<p>${outline}</p>`;
  if (cover_photo) {
    let imgTag = `<div class="column is-4" style="text-align: right;"><img class="article-list-img-h" src="/${cover_photo}"></div>`;
    overviewContent = `<div class="column is-8">${subtitleTag}${outlineTag}</div>${imgTag}`;
  } else {
    overviewContent = `<div class="column is-12">${subtitleTag}${outlineTag}</div>`;
  }

  return `<div class="tile is-ancestor">
                <div class="tile is-parent">
                    <div class="tile is-child box article-list-container">
                        <div data-articleid=${id}></div>
                        <div class="article-list-tag">${adminTag}${categoryTag}${tagHTML}</div>
                        ${titleTag}
                        <div class="columns overview-content">${overviewContent}</div>
                    </div>
                </div>
            </div>`;
};

const checkStatus = async (resp) => {
  const contentType = resp.headers.get('content-type');

  if (contentType && contentType.indexOf('application/json') !== -1 && resp.status < 400) {
    return Promise.resolve(resp);
  }
  return Promise.reject(resp);
};

const fetchSucceed = async (resp) => {
  await resp.json().then((data) => {
    offset += data.size;
    if (data.size < limit) {
      stopFetching = true;
    }

    data.articleList = data.articleList || []; // If there is no data, the empty array returned by backend becomes null
    if (data.articleList.length == 0) {
      return;
    }

    let newContent = '';
    data.articleList.forEach((a) => {
      newContent += formatArticle(a);
    });
    appendNewContent(newContent);
  });
  return Promise.resolve(true);
};

const fetchFailed = async (resp) => {
  await resp.json().then((data) => {
    c('Error: ', data.errHead, data.errBody);
  });
  return Promise.resolve(false);
};

const fetchContent = async (count) => {
  if (stopFetching) {
    return;
  }

  const baseURL = new URL(window.location.href);
  const url = new URL(bookmarkEndpoint, baseURL);
  const paras = { offset: offset, limit: limit };

  for (const [key, val] of Object.entries(paras)) {
    url.searchParams.set(key, val);
  }

  const res = await fetchData(url).then(checkStatus).then(fetchSucceed).catch(fetchFailed);
  if (!res) {
    if (count < 3) {
      fetchContent(++count); // Try again
    } else {
      showErrMsg('Failed to Fetch Content', 'Please reload the page and try again.');
    }
  }
};

const initialFetch = async () => {
  await fetchContent(0);

  if (offset == 0) {
    showNoticeMsg("You haven't saved any articles");
  } else {
    /*
     * For an edge case: the initial articles' height is smaller than the window's height,
     * so user may think that there is no more content and stop scrolling down.
     * Since without scrolling, the fetch event will not be triggered, add one more fetch after the initial fetch.
     * This may happen at weekly-update page or user is visiting website with a vertical screen.
     */
    setTimeout(() => fetchContent(0), 200);
  }
};

const bookmarkHandler = () => {
  offset = 0; // offset == # means we'll skip # articles in next fetch

  initialFetch();

  articlesContainer = document.querySelector('#articles-container');
  articlesContainer.addEventListener('click', (e) => {
    window.location.href = '/articles/browse?articleId=' + e.target.closest('div.tile.is-child').children[0].dataset.articleid;
  });

  const delay = 300;
  let lastFetch = 0;
  window.onscroll = () => {
    if (lastFetch < Date.now() - delay && window.innerHeight + window.pageYOffset >= document.body.offsetHeight * fetchNewContentAnchor) {
      lastFetch = Date.now();
      fetchContent(0);
    }
  };
};

onDOMContentLoaded = (function () {
  bookmarkHandler();
})();
