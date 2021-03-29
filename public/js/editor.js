const TAGS_LIMIT = 5;
const TAGS_CHAR_LIMIT = 20;
const FILES_UPLOAD_LIMIT = 10;
const FILE_ID_LENGTH = 10;
const FILE_MAX_SIZE = 4 * 1000 * 1000; // 4MB
const ACCEPT_FILE_TYPE = {
  image: ["image/png", "image/jpeg", "image/gif", "image/webp", "image/apng"],
};
const errInputMsg = {
  empty: "This field can't be empty.",
  long: "This field can have no more than 255 characters.",
  dateFormat: "The format of date should be yyyy-mm-dd.", // For browsers don't support input type date
  dateIllegal: "The date is illegal.",
  dateTooOld: "The date chosen should be greater than 1960-01-01.",
  dateFuture: "The date chosen should be smaller than the current year.",
  tagsTooMany: `You can target up to ${TAGS_LIMIT} tags at a time.`,
  tagsTooLong: `Each tag can contain at most ${TAGS_CHAR_LIMIT} characters.`,
};
const editorPlaceholder = `Tip:\nUpload images in advance so that you can get the URL to embed them (4MB per image).\n\nShortcuts:\nCtrl-B   -   toggleBold\nCtrl-I    -   toggleItalic\nCtrl-K   -   drawLink\nCtrl-H   -   toggleHeadingSmaller\nShift-Ctrl-H  -  toggleHeadingBigger`;
const createArticleEndpoint = "/admin/create/article";
const updateArticleEndpoint = "/admin/update/article";
const baseURL = window.location.protocol + "//" + window.location.host + "/";
let easyMDE;
let INITIAL_INPUT_SIZE = 0; // Autosave only if user have typed something

const getLocalStorageKey = () => {
  let key = "article-create";
  if (window.location.pathname.indexOf("create") == -1) {
    key = "article-update-" + new URLSearchParams(window.location.search).get("articleId");
  }
  return key;
};

const clearInputLocalStorage = () => {
  window.localStorage.removeItem(getLocalStorageKey());
};

const calculateObjectValueSize = (obj) => {
  return Object.entries(obj).reduce((ttl, val) => (ttl += val[1].length || 0), 0); // In case of non-string type
};

const saveInputToLocalStorage = () => {
  const key = getLocalStorageKey();
  const values = getInputValue();
  const totalInputSize = calculateObjectValueSize(values);

  if (Math.abs(totalInputSize - INITIAL_INPUT_SIZE) > 0) {
    // If user is updating articles, the INITIAL_INPUT_SIZE may be super large
    window.localStorage.setItem(key, JSON.stringify(values));
  }
};

const setupLocalStorage = () => {
  const key = getLocalStorageKey();
  const values = window.localStorage.getItem(key);
  if (!values) {
    showNoticeMsg("Article will be automatically saved", "Enjoy your journey : )");
  } else {
    showNoticeMsg("You can now continue editing", "Enjoy your journey : )");
    writeLocalStorageValue(values);
  }
};

const encodeHTMLEntities = (val) => {
  /*
   * Input: <scrpit>console.log(1)</script>
   * Output: &lt;script&gt;console.log(1)&lt;/script&gt;
   */
  let e = document.createElement("textarea");
  e.innerHTML = val;
  return e.innerHTML;
};

const writeLocalStorageValue = (values) => {
  const { adminOnly, title, subtitle, date, authors, category, tags, content } = JSON.parse(values);

  document.querySelector("#adminOnly").checked = adminOnly;
  document.getElementsByName("title")[0].value = title;
  document.getElementsByName("subtitle")[0].value = subtitle;
  document.getElementsByName("date")[0].value = date;
  document.getElementsByName("category")[0].value = category;
  easyMDE.value(content); // Not working: document.getElementsByName("content")[0].value = content;

  [...document.getElementsByName("authors")].forEach((ele, idx) => {
    if (authors.includes(ele.value)) {
      document.getElementsByName("authors")[idx].checked = true;
    }
  });

  var tagsHTMLHead = '<span class="tag is-warning is-medium" name="tags" style="margin-right: 8px; margin-bottom: 5px">';
  var tagsHTMLTail = '<button class="delete is-small"></button></span>';
  var tagsBody = "";
  tags.forEach((ele) => {
    tagsBody += `${tagsHTMLHead}${encodeHTMLEntities(ele)}${tagsHTMLTail}`;
  });
  document.getElementById("tags-list").innerHTML = tagsBody;
};

const loadMarkdownEditor = () => {
  // https://github.com/Ionaru/easy-markdown-editor
  return new EasyMDE({
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
    imageAccept: "image/png, image/jpeg",
    spellChecker: false,
    tabSize: 4,
    toolbarTips: true,
    imageMaxSize: 1024 * 1024 * 4, // 4 Mb
  });
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

    var newTag = `<span class="tag is-warning is-medium" name="tags" style="margin-right: 8px; margin-bottom: 5px">${val}<button class="delete is-small"></button></span>`;
    tagsList.innerHTML += newTag;
    tagsInputBox.value = "";
  }
};

const tagsDeconstructor = (e) => {
  if (e.target.tagName.toLowerCase() == "button") {
    // c(e.currentTarget); The #tag-list element which registered this event listener's callback function
    tagsCount -= 1;
    e.target.parentNode.remove();
  }
};

const getInputValue = () => {
  var adminOnly = document.querySelector("#adminOnly").checked;
  var title = document.getElementsByName("title")[0].value.trim();
  var subtitle = document.getElementsByName("subtitle")[0].value.trim();
  var date = document.getElementsByName("date")[0].value;
  var authors = [...document.getElementsByName("authors")].filter((author) => author.checked).map((author) => author.value);
  // If toLowerCase() is omitted, sqlite can't find out records in function GetSameCategoryArticles() (but MySQL works fine)
  var category = document.getElementsByName("category")[0].value.toLowerCase(); // From "Medication" to "medication"
  var tags = [...document.getElementsByName("tags")].filter((tag) => tag.tagName.toLowerCase() == "span").map((tag) => tag.textContent.trim());
  var content = easyMDE.value();

  return { adminOnly: adminOnly, title: title, subtitle: subtitle, date: date, authors: authors, category: category, tags: tags, content: content };
};

const validateInput = (values) => {
  const { adminOnly, title, subtitle, date, authors, category, tags, content } = values;
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
  }
  return Promise.resolve(resp);
};

const fetchSucceed = async (resp) => {
  // [...resp.headers.entries()].forEach(header => console.log(header[0], header[1]));
  /*
   * console.log(resp.redirected, resp.url);
   * Respond status code 201: false "http://127.0.0.1:8080/admin/create/article" (same URL). Location header can be retrieved
   * Respond status code 302: true "http://127.0.0.1:8080/articles/browse?articleId=66". Location can NOT be retrieved
   *                          The requirement of resp.redirected == true is that server sends respond with status code 3XX and Location header
   *                          As setting redirect to "follow", browser will send a request via the Location header (but won't render website)
   *                          Thus the following manual redirect (window.location.header = resp.url) will send the same request to server again (waste bandwidth)
   */
  window.localStorage.removeItem(getLocalStorageKey());
  window.location.href = resp.headers.get("Location");
  return Promise.resolve();
};

const fetchFailed = async (resp) => {
  if (resp.status >= 400 && resp.status <= 500) {
    resp.json().then(function (data) {
      showErrMsg(data.errHead, data.errBody);
      for (var key in data.errTags) {
        document.getElementById(`err_msg_${key}`).innerText = data.err[key];
      }
    });
  } else {
    showErrMsg("An Error Occurred !", "Please reload the page and try again.");
  }
};

const submitArticle = async (method, url, formData) => {
  let csrfToken = document.querySelector('meta[name="csrf-token"]').content;
  const headers = new Headers({ "X-CSRF-TOKEN": csrfToken });

  await fetch(url, {
    method: method,
    headers: headers,
    body: formData,
  })
    .then(checkStatus)
    .then(fetchSucceed)
    .catch(fetchFailed);
};

const generateForm = (values) => {
  const { adminOnly, title, subtitle, date, authors, category, tags, content } = values;
  const formData = new FormData();

  formData.append("adminOnly", adminOnly);
  formData.append("title", title);
  formData.append("subtitle", subtitle);
  formData.append("date", date);
  formData.append("authors", authors);
  formData.append("category", category);
  formData.append("tags", tags);
  formData.append("content", content);

  const fileField = document.querySelectorAll('input[type="file"]');
  for (const f of fileField) {
    if (f.files[0] === undefined) {
      continue;
    } else if (f.files[0].size > FILE_MAX_SIZE) {
      document.getElementById("err_msg_content").innerText = `File size of ${f.files[0].name} is too large!`;
      return;
    } else if (ACCEPT_FILE_TYPE.image.indexOf(f.files[0].type) == -1) {
      document.getElementById("err_msg_content").innerText = `File type of ${f.files[0].name} is not permitted!`;
      return;
    }
    let fakeID = f.nextElementSibling.nextElementSibling.nextElementSibling.innerText.substr(-FILE_ID_LENGTH); // Fake image ID provided to user
    let newName = fakeID;
    formData.append("uploadImages", f.files[0], newName); // Note: f.files[0].name is readonly (can't change a name of a created file), so we add customized name as 3rd argument
  }
  return formData;
};

const submitHandler = async (method, endpoint, button) => {
  button.classList.add("is-loading");
  const values = getInputValue();
  let res = validateInput(values);

  if (!res) {
    button.classList.remove("is-loading");
  } else {
    const formData = generateForm(values);
    if (formData != null) {
      await submitArticle(method, endpoint, formData);
    }
    button.classList.remove("is-loading");
  }
};

const submitPost = () => {
  submitHandler("POST", createArticleEndpoint, submitBtn);
};

const savePost = () => {
  let articleId = new URLSearchParams(window.location.search).get("articleId");
  let para = "?" + new URLSearchParams({ articleId: articleId });

  submitHandler("PUT", updateArticleEndpoint + para, saveBtn);
};

onDOMContentLoaded = (function () {
  easyMDE = loadMarkdownEditor();

  setupLocalStorage();

  const values = getInputValue();
  INITIAL_INPUT_SIZE = calculateObjectValueSize(values);

  window.setInterval(() => saveInputToLocalStorage(), 5000);
  document.getElementById("clearAutosaveButton").addEventListener("click", () => clearInputLocalStorage());
  document.getElementById("saveNowButton").addEventListener("click", () => saveInputToLocalStorage());

  // Since the prompt message can't be customized, and the default message is "Changes you made may not be saved"
  // which is quite confusing (cause the changes have been saved). So I remove this feature.
  // window.addEventListener("beforeunload", (e) => {
  //   saveInputToLocalStorage();
  //   e.preventDefault(); // If you prevent default behavior in Mozilla Firefox prompt will always be shown
  //   e.returnValue = ""; // Chrome requires returnValue to be set
  //   return "Are you sure you want to leave? All your changes will be saved."; // Be ignored in modern browsers
  // });

  let filesCount = 1;
  let fileUploadDefaultMsg = "No image uploaded";
  let fileGroups = document.getElementById("fileInputGroups");

  fileGroups.addEventListener("change", (e) => {
    if (e.target.tagName.toLowerCase() == "input") {
      displayFileNameAndCancelBtn(e.target);
    }
  });
  fileGroups.addEventListener("click", (e) => {
    if (e.target.tagName.toLowerCase() == "button") {
      removeFileUploadTag(e.target.parentNode.parentNode.parentNode);
    }
  });
  const createFileIDField = () => {
    let imgURL = baseURL + generateId(FILE_ID_LENGTH);
    let d = document.createElement("span");
    let classesToAdd = ["file-name", "fake-id"];
    d.classList.add(...classesToAdd);
    d.style.paddingRight = "15px";
    d.style.cursor = "default";
    d.style.userSelect = "all";
    d.style.WebkitTransition.userSelect = "all"; // Chrome 49+
    d.textContent = imgURL;
    d.addEventListener("click", (e) => e.preventDefault()); // User has to click on this field to copy fake URL
    return d;
  };
  const createFileUploadTag = () => {
    filesCount += 1;
    let fileInputTemplate = `<label class='file-label'>
                <input class='file-input' type='file'>
                <span class='file-cta'>
                  <span class='file-icon'> 📂 </span>
                  <span class='file-label'>Upload images</span>
                </span>
                <span class='file-name'>${fileUploadDefaultMsg}</span>
              </label>`;
    let d = document.createElement("div");
    let classesToAdd = ["file", "has-name", "is-warning", "is-small", "pb-1"];
    d.classList.add(...classesToAdd);
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
  const displayFileNameAndCancelBtn = (ele) => {
    if (ele.files.length > 0) {
      let fileDeleteBtn = "<button class='delete is-small mr-2'></button>";
      let originalHTML = ele.nextElementSibling.textContent;
      ele.nextElementSibling.innerHTML = fileDeleteBtn + originalHTML;

      let val = ele.nextElementSibling.nextElementSibling.textContent;
      ele.nextElementSibling.nextElementSibling.textContent = ele.files[0].name;
      ele.parentNode.appendChild(createFileIDField());
      if (val == fileUploadDefaultMsg && filesCount < FILES_UPLOAD_LIMIT) {
        fileGroups.appendChild(createFileUploadTag());
      }
    }
  };

  tagsCount = 0;
  tagsInputBox = document.querySelector("input[name='tags']");
  tagsInputBox.addEventListener("keyup", tagsConstructor);
  tagsList = document.querySelector("#tags-list");
  tagsList.addEventListener("click", tagsDeconstructor);

  submitBtn = document.querySelector("#submitButton");
  if (submitBtn) {
    submitBtn.addEventListener("click", submitPost);
  }

  saveBtn = document.getElementById("saveButton");
  if (saveBtn) {
    saveBtn.addEventListener("click", savePost);
  }

  cancelBtn = document.getElementById("cancelButton");
  if (cancelBtn) {
    cancelBtn.addEventListener("click", () => {
      let articleId = new URLSearchParams(window.location.search).get("articleId");
      let para = "?" + new URLSearchParams({ articleId: articleId });
      window.location.href = "/articles/browse" + para;
    });
  }
})();
