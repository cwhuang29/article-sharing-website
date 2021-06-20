const landingPageEndpoint = '/articles/weekly-update';
const checkPermissionEndpoint = '/admin/check-permisssion';
const modifyEndpoint = '/admin/update/article';
const deleteEndpoint = '/admin/delete/article';
const bookmarkEndpoint = '/articles/bookmark';
const likeEndpoint = '/articles/like';
let modalMode = '';
let actionURL = '';

const openModalBody = (mode, title) => {
  modalMode = mode;
  confirmModalTitle.innerText = title;
  confirmModalBody.classList.add('is-active');
};

const closeModalBody = () => {
  modalMode = '';
  confirmModalTitle.innerText = '';
  confirmModalBody.classList.remove('is-active');
};

const fetchSucceed = async (resp) => {
  return Promise.resolve(1);
};

const fetchFailed = async (resp) => {
  resp.json().then(function (data) {
    showErrMsg(data.errHead, data.errBody);
  });
  return Promise.resolve(0);
};

const checkStatus = async (resp) => {
  if (resp.status >= 400) {
    return Promise.reject(resp);
  }
  return Promise.resolve(resp);
};

const fetchDeleteReq = async (url) => {
  const csrfToken = document.querySelector('meta[name="csrf-token"]').content;
  const headers = { 'X-CSRF-TOKEN': csrfToken };

  return fetchData(url, {
    method: 'DELETE',
    headers: headers,
    cache: 'no-cache',
  })
    .then(checkStatus)
    .then(fetchSucceed)
    .catch(fetchFailed);
};

const modifyOrDelete = async () => {
  const endpoint = modalMode == 'update' ? modifyEndpoint : modalMode == 'delete' ? deleteEndpoint : '';
  if (!endpoint) {
    return;
  }

  const articleId = new URLSearchParams(window.location.search).get('articleId');
  const baseURL = new URL(window.location.href);
  const url = new URL(endpoint, baseURL);
  const checkUrl = new URL(checkPermissionEndpoint, baseURL);

  url.searchParams.set('articleId', articleId);
  checkUrl.searchParams.set('articleId', articleId);

  let res = await fetchData(checkUrl, {
    cache: 'no-cache',
  })
    .then(checkStatus)
    .then(fetchSucceed)
    .catch(fetchFailed);

  if (!res) {
    closeModalBody();
    return;
  }
  if (modalMode == 'update') {
    window.location = url;
  } else {
    res = await fetchDeleteReq(url);
    if (res == 1) {
      window.location = landingPageEndpoint;
    }
    closeModalBody();
  }
};

const modifyArticle = () => {
  mode = 'update';
  title = 'Are you sure you want to modify this article?';
  openModalBody(mode, title);
};

const deleteArticle = () => {
  mode = 'delete';
  title = 'Are you sure you want to delete this article?';
  openModalBody(mode, title);
};

const switchIconStatus = (yes, yesIcon, noIcon) => {
  if (yes) {
    yesIcon.style.display = 'block';
    noIcon.style.display = 'none';
  } else {
    noIcon.style.display = 'block';
    yesIcon.style.display = 'none';
  }
};

const generateURL = (endpoint, paraKey, paraValue) => {
  const articleId = new URLSearchParams(window.location.search).get('articleId');
  const baseURL = new URL(window.location.href);
  const url = new URL(endpoint + `/${articleId}`, baseURL);

  url.searchParams.set(paraKey, paraValue);
  return url;
};

const bookmarkSuccess = async (resp) => {
  await resp.json().then(function (data) {
    switchIconStatus(data.isBookmarked, bookmarkIconYes, bookmarkIconNo);
    bookmarkParent.dataset.bookmarked = data.isBookmarked;
  });
};

const likeSuccess = async (resp) => {
  await resp.json().then(function (data) {
    switchIconStatus(data.isLiked, likeIconYes, likeIconNo);
    likeParent.dataset.liked = data.isLiked;
  });
};

const updateArticleStatus = async (url, method, updateSuccessFunc) => {
  await fetchData(url, { method: method }).then(checkStatus).then(updateSuccessFunc).catch(fetchFailed);
};

const updateBookmarkStatus = () => {
  const paraKey = 'bookmarked';
  const paraValue = parseInt(bookmarkParent.dataset.bookmarked) === 0 ? 1 : 0;
  const url = generateURL(bookmarkEndpoint, paraKey, paraValue);
  updateArticleStatus(url, 'PUT', bookmarkSuccess);
};

const updateLikeStatus = () => {
  const paraKey = 'liked';
  const paraValue = parseInt(likeParent.dataset.liked) === 0 ? 1 : 0;
  const url = generateURL(likeEndpoint, paraKey, paraValue);
  updateArticleStatus(url, 'PUT', likeSuccess);
};

const initialBookmark = async () => {
  const paraKey = 'bookmarked';
  const paraValue = parseInt(bookmarkParent.dataset.bookmarked) === 0 ? 0 : 1;
  const url = generateURL(bookmarkEndpoint, paraKey, paraValue);
  await updateArticleStatus(url, 'GET', bookmarkSuccess);
  bookmarkParent.style.display = 'block';
};

const initialLike = async () => {
  const paraKey = 'liked';
  const paraValue = parseInt(likeParent.dataset.liked) === 0 ? 0 : 1;
  const url = generateURL(likeEndpoint, paraKey, paraValue);
  await updateArticleStatus(url, 'GET', likeSuccess);
  likeParent.style.display = 'block';
};

const browseHandler = () => {
  if (getCookie('login_email')) {
    bookmarkIconNo = document.getElementById('bookmarkIconNo');
    bookmarkIconYes = document.getElementById('bookmarkIconYes');
    bookmarkParent = document.getElementById('bookmarkParent');
    bookmarkParent.addEventListener('click', updateBookmarkStatus);

    likeIconNo = document.getElementById('likeIconNo');
    likeIconYes = document.getElementById('likeIconYes');
    likeParent = document.getElementById('likeParent');
    likeParent.addEventListener('click', updateLikeStatus);

    initialBookmark();
    initialLike();
  }

  if (getCookie('is_admin')) {
    adminSection = document.getElementById('adminSection');
    modifyBtn = document.getElementById('modifyBtn');
    deleteBtn = document.getElementById('deleteBtn');
    confirmModalBody = document.getElementsByClassName('modal')[0];
    confirmModalClose = document.getElementsByClassName('modal-close')[0];
    confirmModalTitle = document.getElementById('confirm-modal-title');
    yesBtn = document.getElementById('yesBtn');
    noBtn = document.getElementById('noBtn');

    adminSection.style.display = 'block';
    confirmModalClose.addEventListener('click', closeModalBody);
    noBtn.addEventListener('click', closeModalBody);
    yesBtn.addEventListener('click', modifyOrDelete);
    modifyBtn.addEventListener('click', modifyArticle);
    deleteBtn.addEventListener('click', deleteArticle);
  }
};

onDOMContentLoaded = (function () {
  browseHandler();
})();
