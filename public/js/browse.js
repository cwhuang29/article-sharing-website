const checkPermissionEndpoint = "/admin/check-permisssion";
const modifyEndpoint = "/admin/update/article";
const deleteEndpoint = "/admin/delete/article";
const landingPage = "/articles/weekly-update";
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

const checkStatus = async (resp) => {
    if (resp.status >= 400) {
        return Promise.reject(resp);
    }
    return Promise.resolve(resp);
};

const creationFailed = async (resp) => {
    if (resp.status < 500) {
        resp.json().then(function (data) {
            showErrMsg(data.err);
        });
    } else {
        showErrMsg("<div><p><strong>Some severe errors occurred !</strong></p><p>Please reload the page and try again.</p></div>");
    }
    return Promise.resolve(0);
};

const creationSucceed = async (resp) => {
    return Promise.resolve(1);
};

const fetchDeleteReq = async (url) => {
    return fetch(url, {
        method: "DELETE",
        mode: "cors",
        cache: "no-cache",
        credentials: "same-origin",
        redirect: "follow",
        referrerPolicy: "no-referrer",
    })
        .then(checkStatus)
        .then(creationSucceed)
        .catch(creationFailed);
};

const modifyOrDelete = async () => {
    baseURL = modalMode == "update" ? modifyEndpoint : modalMode == "delete" ? deleteEndpoint : "";
    if (!baseURL) {
        return;
    }

    let articleId = new URLSearchParams(window.location.search).get("articleId");
    let para = "?" + new URLSearchParams({ articleId: articleId });
    actionURL = baseURL + para; // e.g. /admin/update/article?articleId=8

    let res = await fetch(checkPermissionEndpoint + para, {
        method: "GET",
        mode: "cors",
        cache: "no-cache",
        credentials: "same-origin",
        redirect: "follow",
        referrerPolicy: "no-referrer",
    })
        .then(checkStatus)
        .then(creationSucceed)
        .catch(creationFailed);

    if (!res) {
        return;
    }

    if (modalMode == "update") {
        window.location = actionURL;
    } else {
        res = await fetchDeleteReq(actionURL);
        if (res == 1) {
            window.location = landingPage;
        }
    }
    closeModalBody();
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

onDOMContentLoaded = (function () {
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
})();