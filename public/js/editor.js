const TAGS_LIMIT = 5;
const TAGS_CHAR_LIMIT = 20;
const FILES_UPLOAD_LIMIT = 10;
const errInputMsg = {
    empty: "This field can't be empty.",
    long: "This field can have no more than 255 characters.",
    dateFormat: "The format of date should be yyyy-mm-dd.", // For safari which don't support input type date
    dateIllegal: "The date is illegal.",
    dateTooOld: "The date chosen should be greater than 1960-01-01.",
    dateFuture: "The date chosen should be smaller than the current year.",
    tagsTooMany: `You can target up to ${TAGS_LIMIT} tags at a time.`,
    tagsTooLong: `Each tag can contain at most ${TAGS_CHAR_LIMIT} characters.`,
};
const editorPlaceholder = `Tip:\nUpload images in advance so that you can get the URL to embed images.\n\nShortcuts:\nCtrl-B   -   toggleBold\nCtrl-I    -   toggleItalic\nCtrl-K   -   drawLink\nCtrl-H   -   toggleHeadingSmaller\nShift-Ctrl-H  -  toggleHeadingBigger`;
const createImagesEndpoint = "/admin/create/images";
const createArticleEndpoint = "/admin/create/article";
const updateArticleEndpoint = "/admin/update/article";

const loadMarkdownEditor = () => {
    // https://github.com/Ionaru/easy-markdown-editor
    return (easyMDE = new EasyMDE({
        element: document.getElementById("content-text-area"),
        previewRender: function (plainText) {
            c(marked(plainText));
            return marked(plainText);
        },
        autoDownloadFontAwesome: true,
        showIcons: ["italic", "|", "bold", "strikethrough", "code", "redo", "|", "undo"],
        // showIcons: ['strikethrough', 'code', 'table', 'redo', 'heading', 'undo', 'heading-bigger', 'heading-smaller', 'heading-1', 'heading-2', 'heading-3', 'clean-block', 'horizontal-rule'],
        lineNumbers: false,
        initialValue: "",
        minHeight: "250px",
        maxHeight: "400px",
        placeholder: editorPlaceholder,
        imageAccept: "image/png, image/jpeg", // Should check again at server side
        spellChecker: false,
        tabSize: 4,
        toolbarTips: true,
        imageMaxSize: 1024 * 1024 * 4, // 4 Mb
    }));
};

const tagsConstructor = (e) => {
    if ((e.key === "Enter" || e.keyCode === 13) && e.shiftKey) {
        if (tagsCount >= 5) {
            document.getElementById("err_msg_tags").innerText = errInputMsg.tagsTooMany;
            return;
        }
        var val = tagsInputBox.value.trim();
        if (val == "") {
            return;
        } else if (val.length > TAGS_CHAR_LIMIT) {
            document.getElementById("err_msg_tags").innerText = errInputMsg.tagsTooLong;
            return;
        }
        document.getElementById("err_msg_tags").innerText = "";
        tagsCount += 1;
        var newTag = `<span class="tag is-warning is-medium" name="tags" style="margin-right: 9px;">${val}<button class="delete is-small"></button></span>`;
        tagsList.innerHTML += newTag;
        tagsInputBox.value = "";
    }
};
const tagsDeconstructor = (e) => {
    if (e.target.tagName == "BUTTON") {
        // c(e.currentTarget); The #tag-list element
        tagsCount -= 1;
        e.target.parentNode.remove();
    }
};

const getInputValue = (easyMDE) => {
    var title = document.getElementsByName("title")[0].value.trim();
    var subtitle = document.getElementsByName("subtitle")[0].value.trim();
    var date = document.getElementsByName("date")[0].value;
    var authors = [...document.getElementsByName("authors")]
        .filter((author) => {
            return author.checked;
        })
        .map((author) => {
            return author.value;
        });
    var category = document.getElementsByName("category")[0].value;
    var tags = [...document.getElementsByName("tags")]
        .filter((tag) => {
            return tag.tagName == "SPAN";
        })
        .map((tag) => {
            return tag.textContent.trim();
        });
    // Insert newlines into head and tail of images. Otherwise the <img> tag will not embedded in <figure> tag after transformation
    var content = easyMDE.value().replace(/(!\[.*\]\(.*\))/g, "\n\n$1\n\n");

    return [title, subtitle, date, authors, category, tags, content];
};

const validateInput = (title, subtitle, date, authors, category, tags, content) => {
    var canSubmit = true;

    if (title.length == 0) {
        canSubmit = false;
        document.getElementById("err_msg_title").innerText = errInputMsg.empty;
    } else if (title.length > 255) {
        canSubmit = false;
        document.getElementById("err_msg_title").innerText = errInputMsg.long;
    } else {
        document.getElementById("err_msg_title").innerText = "";
    }

    if (subtitle.length > 255) {
        // Subtitle can be empty
        canSubmit = false;
        document.getElementById("err_msg_subtitle").innerText = errInputMsg.long;
    } else {
        document.getElementById("err_msg_subtitle").innerText = "";
    }

    // var currentTime = new Date();
    if (!/^\d\d\d\d-\d\d-\d\d$/.test(date)) {
        canSubmit = false;
        document.getElementById("err_msg_date").innerText = errInputMsg.dateFormat;
    } else if (!/^(19[6-9][0-9]|2\d\d\d)-(3[01]|[12][0-9]|0[1-9])-(0[1-9]|[12][0-9]|3[01])$/.test(date)) {
        // The regex can't detect 02/30, 04/31 ... . Nevertheless, Beckend will fix this error
        canSubmit = false;
        document.getElementById("err_msg_date").innerText = errInputMsg.dateIllegal;
    } else if (Number(date.split("-")[0]) < 1960) {
        // Can't detect if the input date is in the current year but in a future month and/or day. Nevertheless, Beckend will fix this error
        canSubmit = false;
        document.getElementById("err_msg_date").innerText = errInputMsg.dateTooOld;
        // } else if (Number(date.split('-')[0]) > currentTime.getFullYear()) {
        //     canSubmit = false;
        //     document.getElementById('err_msg_date').innerText = errInputMsg.dateFuture;
    } else {
        document.getElementById("err_msg_date").innerText = "";
    }

    if (authors.length == 0) {
        canSubmit = false;
        document.getElementById("err_msg_authors").innerText = errInputMsg.empty;
    } else {
        document.getElementById("err_msg_authors").innerText = "";
    }

    if (tags.length > TAGS_LIMIT) {
        canSubmit = false;
        document.getElementById("err_msg_tags").innerText = errInputMsg.tagsTooMany;
    } else {
        document.getElementById("err_msg_tags").innerText = "";
    }
    for (let t of tags) {
        if (t.length > TAGS_CHAR_LIMIT) {
            canSubmit = false;
            document.getElementById("err_msg_tags").innerText = errInputMsg.tagsTooLong;
            break;
        }
    }

    if (content.length == 0) {
        canSubmit = false;
        document.getElementById("err_msg_content").innerText = errInputMsg.empty;
    } else {
        document.getElementById("err_msg_content").innerText = "";
    }

    return canSubmit;
};

const dec2hex = (dec) => {
    return dec.toString(16).padStart(2, "0");
};

const generateId = (len) => {
    var arr = new Uint8Array((len || 40) / 2);
    window.crypto.getRandomValues(arr);
    return Array.from(arr, dec2hex).join("");
};

const checkStatus = async (resp) => {
    if (resp.status >= 400) {
        return Promise.reject(resp);
    } else {
        return Promise.resolve(resp);
    }
};

const creationSucceed = async (resp) => {
    // [...resp.headers.entries()].forEach(header => console.log(header[0], header[1]));
    /*
     * console.log(resp.redirected, resp.url);
     * Respond status code 201: false "http://127.0.0.1:8080/admin/create/article" (same URL). Location header can be retrieved
     * Respond status code 302: true "http://127.0.0.1:8080/articles/browse?articleId=66". Location can NOT be retrieved
     *                          The requirement of resp.redirected == true is that server sends respond with status code 3XX and Location header
     *                          As setting redirect to "follow", browser will send a request via the Location header (but won't render website)
     *                          Thus the following manual redirect (window.location.header = resp.url) will send the same request to server again (waste bandwidth)
     */
    window.location.href = resp.headers.get("Location");
    return Promise.resolve();
};

const creationFailed = async (resp) => {
    if (resp.status == 500) {
        resp.json().then(function (data) {
            showErrMsg(`<div><p><strong>Error !</strong></p><p>${data.err}</p></div>`);
        });
    } else if (resp.status == 400) {
        resp.json().then(function (data) {
            if (data.bindingError) {
                showErrMsg("<div><p><strong>Some severe errors occurred !</strong></p><p>Please reload the page and try again.</p></div>");
            } else {
                for (var key in data.err) {
                    document.getElementById(`err_msg_${key}`).innerText = data.err[key];
                }
            }
        });
    } else {
        showErrMsg("<div><p><strong>Some severe errors occurred !</strong></p><p>Please reload the page and try again.</p></div>");
    }
};

const postData = async (url = "", method = "POST", data = {}) => {
    return (response = await fetch(url, {
        method: method,
        mode: "cors", // no-cors, *cors, same-origin
        cache: "no-cache", // *default, no-cache, reload, force-cache, only-if-cached
        credentials: "same-origin", // include, *same-origin, omit
        headers: {
            "Content-Type": "application/json", // 'application/x-www-form-urlencoded',
        },
        redirect: "follow", // manual, *follow, error
        referrerPolicy: "no-referrer", // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
        body: JSON.stringify(data),
    }));
};

const uploadMultipleImages = async (url, data) => {
    return (response = await fetch(url, {
        // Do NOT set Content-Type header (browser will set this header automatically)
        method: "POST",
        body: data,
    }));
};

const submitArticle = async (method, url, formData, formNameIDMapping, title, subtitle, date, authors, category, tags, content) => {
    if (formData != null) {
        var success = await uploadMultipleImages(createImagesEndpoint, formData)
            .then((response) => response.json())
            .then((data) => {
                for (const fileName in data) {
                    content = content.replaceAll(formNameIDMapping[fileName], data[fileName]);
                }
                return Promise.resolve(1);
            })
            .catch((err) => {
                showErrMsg(`<div><p><strong>Image upload failed !</strong></p><p>Please try again later.</p></div>`);
                c("Images upload error:", err);
                return 0;
            });
        if (!success) {
            return;
        }
    }

    await postData(url, method, {
        title: title,
        subtitle: subtitle,
        date: date,
        authors: authors,
        category: category,
        tags: tags,
        content: content,
    })
        .then(checkStatus)
        .then(creationSucceed)
        .catch(creationFailed);
};

const getImagesData = () => {
    let hasImage = false;
    let formNameIDMapping = {};
    const imagesData = new FormData();
    const fileField = document.querySelectorAll('input[type="file"]');
    for (const f of fileField) {
        if (f.files[0] === undefined) {
            continue;
        }
        hasImage = true;
        formNameIDMapping[f.files[0].name] = f.nextElementSibling.nextElementSibling.nextElementSibling.innerText; // Uploaded file name - fake ID of the file
        imagesData.append("uploadImage", f.files[0]);
        // File {name: "test_image.png", lastModified: 1567836656000, lastModifiedDate: Sat Sep 07 2019 14:10:56 GMT+0800 (Taipei Standard Time), webkitRelativePath: "", size: 939134, …}
    }
    if (!hasImage) {
        return { imagesData: null, formNameIDMapping: null };
    }
    return { imagesData: imagesData, formNameIDMapping: formNameIDMapping };
};

const submitPost = async () => {
    submitBtn.classList.add("is-loading");

    let [title, subtitle, date, authors, category, tags, content] = getInputValue(easyMDE);
    let res = validateInput(title, subtitle, date, authors, category, tags, content);

    if (!res) {
        submitBtn.classList.remove("is-loading");
    } else {
        const { imagesData, formNameIDMapping } = getImagesData();
        await submitArticle("POST", createArticleEndpoint, imagesData, formNameIDMapping, title, subtitle, date, authors, category, tags, content);
        submitBtn.classList.remove("is-loading");
    }
};

const savePost = async () => {
    saveBtn.classList.add("is-loading");

    let [title, subtitle, date, authors, category, tags, content] = getInputValue(easyMDE);
    let res = validateInput(title, subtitle, date, authors, category, tags, content);

    if (!res) {
        savetBtn.classList.remove("is-loading");
    } else {
        const { imagesData, formNameIDMapping } = getImagesData();
        let articleId = new URLSearchParams(window.location.search).get("articleId");
        let para = "?" + new URLSearchParams({ articleId: articleId });
        await submitArticle("PUT", updateArticleEndpoint + para, imagesData, formNameIDMapping, title, subtitle, date, authors, category, tags, content);
        saveBtn.classList.remove("is-loading");
    }
};

onDOMContentLoaded = (function () {
    let easyMDE = loadMarkdownEditor();

    let loc = window.location;
    let baseURL = loc.protocol + "//" + loc.host + "/upload/images/";

    const FILE_ID_LENGTH = 8;
    let filesCount = 1;
    let noUploadDefaultMsg = "No image uploaded";
    let fileGroups = document.getElementById("fileInputGroups");

    const createFileIDField = () => {
        var id = generateId(FILE_ID_LENGTH);
        imgURL = baseURL + id;

        var d = document.createElement("span");
        d.classList.add("file-name");
        d.classList.add("fake-id");
        d.style.paddingRight = "265px";
        d.style.cursor = "default";
        d.style.userSelect = "all";
        d.style.WebkitTransition.userSelect = "all"; // Chrome 49+
        d.textContent = imgURL;
        d.addEventListener("click", function (e) {
            e.preventDefault();
        });
        return d;
    };
    const createFileUploadTag = () => {
        if (filesCount >= FILES_UPLOAD_LIMIT) {
            return null;
        }
        filesCount += 1;
        let fileInputTemplate = `<label class='file-label'>
                <input class='file-input' type='file' name='resume'>
                <span class='file-cta'>
                  <span class='file-icon'> 📂 </span>
                  <span class='file-label'>Upload images</span>
                </span>
                <span class='file-name'>${noUploadDefaultMsg}</span>
              </label>`;
        var d = document.createElement("div");
        d.classList.add("file");
        d.classList.add("has-name");
        d.classList.add("is-warning");
        d.classList.add("is-small");
        d.classList.add("pb-1");
        d.innerHTML = fileInputTemplate;
        return d;
    };
    const removeFileUploadTag = (target) => {
        target.remove();
        filesCount -= 1;
        if (filesCount == 0 || filesCount == FILES_UPLOAD_LIMIT - 1) {
            fileGroups.appendChild(createFileUploadTag());
        }
    };
    fileGroups.addEventListener("click", (e) => {
        if (e.target.tagName == "BUTTON") {
            removeFileUploadTag(e.target.parentNode.parentNode.parentNode);
            return;
        }
        fileInput = e.target.closest("input[type=file]");
        if (!fileInput) {
            return;
        }
        fileInput.onchange = () => {
            if (fileInput.files.length > 0) {
                var fileDeleteBtn = "<button class='delete is-small mr-2'></button>";
                var originalHTML = fileInput.nextElementSibling.textContent;
                fileInput.nextElementSibling.innerHTML = fileDeleteBtn + originalHTML;

                var val = fileInput.nextElementSibling.nextElementSibling.textContent;
                fileInput.nextElementSibling.nextElementSibling.textContent = fileInput.files[0].name;
                fileInput.parentNode.appendChild(createFileIDField());
                if (val == noUploadDefaultMsg) {
                    var newNode = createFileUploadTag();
                    if (newNode) {
                        fileGroups.appendChild(newNode);
                    }
                }
            }
        };
    });

    tagsCount = 0;
    tagsInputBox = document.querySelector("input[name='tags']");
    tagsList = document.querySelector("#tags-list");
    tagsInputBox.addEventListener("keyup", tagsConstructor);
    tagsList.addEventListener("click", tagsDeconstructor);

    submitBtn = document.querySelector("#submit_button");
    if (submitBtn) {
        submitBtn.addEventListener("click", submitPost);
    }

    cancelBtn = document.getElementById("cancel_button");
    if (cancelBtn) {
        cancelBtn.addEventListener("click", () => {
            let articleId = new URLSearchParams(window.location.search).get("articleId");
            let para = "?" + new URLSearchParams({ articleId: articleId });
            window.location.href = "/articles/browse" + para;
        });
    }
    saveBtn = document.getElementById("save_button");
    if (saveBtn) {
        saveBtn.addEventListener("click", savePost);
    }
})();
/*
    const uploadSingleImages = () => {
        const formData = new FormData();
        const fileField = document.querySelector('input[type="file"]');

        formData.append('uploadImage', fileField.files[0]);
        c(fileField.files[0]);
        // File {name: "test_image.png", lastModified: 1567836656000, lastModifiedDate: Sat Sep 07 2019 14:10:56 GMT+0800 (Taipei Standard Time), webkitRelativePath: "", size: 939134, …}
        fetch('/admin/create/images', { // Do NOT set Content-Type header (browser will set Content-Type automatically)
            method: 'POST',
            body: formData
        });
    };
    const createArticles = (title, subtitle, date, authors, category, tags, content) => {
        var request = new XMLHttpRequest();
        request.open('GET', `/admin/create/article`, true);
        request.onload = function() {
          if (request.status >= 200 && request.status < 400) {
            console.log(request.responseText);
          }
        };
        request.send();
    }
*/
