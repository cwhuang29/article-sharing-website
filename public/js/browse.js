const landingPageEndpoint = "/articles/weekly-update";
const checkPermissionEndpoint = "/admin/check-permisssion";
const modifyEndpoint = "/admin/update/article";
const deleteEndpoint = "/admin/delete/article";
const bookmarkEndpoint = "/articles/bookmark";
let modalMode = "";
let actionURL = "";

const openModalBody = (mode, title) => {
  modalMode = mode;
  confirmModalTitle.innerText = title;
  confirmModalBody.classList.add("is-active");
};

const closeModalBody = () => {
  modalMode = "";
  confirmModalTitle.innerText = "";
  confirmModalBody.classList.remove("is-active");
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
  const headers = { "X-CSRF-TOKEN": csrfToken };

  return fetchData(url, {
    method: "DELETE",
    headers: headers,
    cache: "no-cache",
  })
    .then(checkStatus)
    .then(fetchSucceed)
    .catch(fetchFailed);
};

const modifyOrDelete = async () => {
  const endpoint = modalMode == "update" ? modifyEndpoint : modalMode == "delete" ? deleteEndpoint : "";
  if (!endpoint) {
    return;
  }

  const articleId = new URLSearchParams(window.location.search).get("articleId");
  const baseURL = new URL(window.location.href);
  const url = new URL(endpoint, baseURL);
  const checkUrl = new URL(checkPermissionEndpoint, baseURL);

  url.searchParams.set("articleId", articleId); // e.g. /admin/update/article?articleId=8
  checkUrl.searchParams.set("articleId", articleId);

  let res = await fetchData(checkUrl, {
    cache: "no-cache",
  })
    .then(checkStatus)
    .then(fetchSucceed)
    .catch(fetchFailed);

  if (!res) {
    closeModalBody();
    return;
  }
  if (modalMode == "update") {
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
  mode = "update";
  title = "Are you sure you want to modify this article?";
  openModalBody(mode, title);
};

const deleteArticle = () => {
  mode = "delete";
  title = "Are you sure you want to delete this article?";
  openModalBody(mode, title);
};

const switchBookmarkIcon = (isBookmarked) => {
  if (isBookmarked) {
    bookmarkIconNo.style.display = "none";
    bookmarkIconYes.style.display = "block";
  } else {
    bookmarkIconNo.style.display = "block";
    bookmarkIconYes.style.display = "none";
  }
};

const getBookmarkURL = () => {
  const isBookmarked = parseInt(bookmarkParent.dataset.bookmarked) === 0 ? 1 : 0;
  const articleId = new URLSearchParams(window.location.search).get("articleId");
  const baseURL = new URL(window.location.href);
  const url = new URL(bookmarkEndpoint + `/${articleId}`, baseURL);

  url.searchParams.set("bookmarked", isBookmarked);
  return url;
};

const bookmarkArticle = async () => {
  url = getBookmarkURL();
  await fetchData(url, { method: "PUT" })
    .then(checkStatus)
    .then((resp) => {
      resp.json().then(function (data) {
        switchBookmarkIcon(data.isBookmarked);
        bookmarkParent.dataset.bookmarked = data.isBookmarked;
      });
    })
    .catch(fetchFailed);
};

const initialBookmark = async () => {
  // Don't show error in this feature
  url = getBookmarkURL();
  await fetchData(url)
    .then(checkStatus)
    .then((resp) => {
      resp.json().then(function (data) {
        switchBookmarkIcon(data.isBookmarked);
        bookmarkParent.dataset.bookmarked = data.isBookmarked;
      });
    });
  bookmarkParent.style.display = "block";
};

const browseHandler = () => {
  if (getCookie("login_email")) {
    initialBookmark();
    bookmarkParent = document.getElementById("bookmarkParent");
    bookmarkIconNo = document.getElementById("bookmarkIconNo");
    bookmarkIconYes = document.getElementById("bookmarkIconYes");
    bookmarkParent.addEventListener("click", bookmarkArticle);
  }

  if (getCookie("is_admin")) {
    adminSection = document.getElementById("adminSection");
    modifyBtn = document.getElementById("modifyBtn");
    deleteBtn = document.getElementById("deleteBtn");
    confirmModalBody = document.getElementsByClassName("modal")[0];
    confirmModalClose = document.getElementsByClassName("modal-close")[0];
    confirmModalTitle = document.getElementById("confirm-modal-title");
    yesBtn = document.getElementById("yesBtn");
    noBtn = document.getElementById("noBtn");

    adminSection.style.display = "block";
    confirmModalClose.addEventListener("click", closeModalBody);
    noBtn.addEventListener("click", closeModalBody);
    yesBtn.addEventListener("click", modifyOrDelete);
    modifyBtn.addEventListener("click", modifyArticle);
    deleteBtn.addEventListener("click", deleteArticle);
  }
};

onDOMContentLoaded = (function () {
  browseHandler();
})();
