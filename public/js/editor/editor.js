const TITLE_BYTES_LIMIT = 255;
const SUBTITLE_BYTES_LIMIT = 255;
const TAGS_LIMIT = 5;
const TAGS_BYTES_LIMIT = 20;
const FILE_MAX_SIZE = 8 * 1000 * 1000; // 8MB
const ACCEPT_FILE_TYPE = {
  image: ['image/png', 'image/jpeg', 'image/gif', 'image/webp', 'image/apng'],
};
const errInputMsg = {
  empty: "This field can't be empty.",
  long: 'This field can have no more than 255 characters.',
  dateFormat: 'The format of date should be yyyy-mm-dd.', // For browsers don't support input type date
  dateIllegal: 'The date is illegal.',
  dateTooOld: 'The date chosen should be greater than 1960-01-01.',
  dateFuture: 'The date chosen should be smaller than the current year.',
  tagsTooMany: `You can target up to ${TAGS_LIMIT} tags at a time.`,
  tagsTooLong: `Each tag can contain at most ${TAGS_BYTES_LIMIT} characters.`,
};
const editorPlaceholder = `Tip:\nUpload images in advance so that you can get the URL to embed them (8MB per image).\n\nShortcuts:\nCtrl-B   -   toggleBold\nCtrl-I    -   toggleItalic\nCtrl-K   -   drawLink\nCtrl-H   -   toggleHeadingSmaller\nShift-Ctrl-H  -  toggleHeadingBigger`;
const createArticleEndpoint = '/admin/create/article';
const updateArticleEndpoint = '/admin/update/article';
const browseArticleEndpoint = '/articles/browse';
let easyMDE;
let INITIAL_INPUT_SIZE = 0; // Autosave only if user have typed something

const loadMarkdownEditor = () => {
  // https://github.com/Ionaru/easy-markdown-editor
  return new EasyMDE({
    element: document.getElementById('content-text-area'),
    previewRender: function (plainText) {
      c(marked(plainText));
      return marked(plainText);
    },
    autoDownloadFontAwesome: true,
    showIcons: ['italic', '|', 'bold', 'strikethrough', 'code', 'redo', '|', 'undo'],
    // showIcons: ['strikethrough', 'code', 'table', 'redo', 'heading', 'undo', 'heading-bigger', 'heading-smaller', 'heading-1', 'heading-2', 'heading-3', 'clean-block', 'horizontal-rule'],
    lineNumbers: false,
    initialValue: '',
    minHeight: '250px',
    maxHeight: '400px',
    placeholder: editorPlaceholder,
    imageAccept: 'image/png, image/jpeg',
    spellChecker: false,
    tabSize: 4,
    toolbarTips: true,
    imageMaxSize: 1024 * 1024 * 4, // 4 Mb
  });
};

const getInputValue = () => {
  const adminOnly = document.querySelector('#adminOnly').checked;
  const title = document.getElementsByName('title')[0].value.trim();
  const subtitle = document.getElementsByName('subtitle')[0].value.trim();
  const date = document.getElementsByName('date')[0].value;
  const authors = [...document.getElementsByName('authors')].filter((author) => author.checked).map((author) => author.value);
  // If toLowerCase() is omitted, sqlite can't find out records in function GetSameCategoryArticles() (but MySQL works fine)
  const category = document.getElementsByName('category')[0].value.toLowerCase(); // From "Medication" to "medication"
  const tags = [...document.getElementsByName('tags')].filter((tag) => tag.tagName.toLowerCase() == 'span').map((tag) => tag.textContent.trim());
  const outline = document.getElementsByName('outline')[0].value;
  const content = easyMDE.value();

  return { adminOnly: adminOnly, title: title, subtitle: subtitle, date: date, authors: authors, category: category, tags: tags, outline: outline, content: content };
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
  window.location.href = resp.headers.get('Location');
  return Promise.resolve();
};

const fetchFailed = async (resp) => {
  if (resp.status >= 400 && resp.status <= 500) {
    resp.json().then(function (data) {
      showErrMsg(data.errHead, data.errBody);
      for (var key in data.errTags) {
        document.getElementById(`err_msg_${key}`).innerText = data.errTags[key];
      }
    });
  } else {
    showErrMsg('An Error Occurred !', 'Please reload the page and try again.');
  }
};

const submitArticle = async (method, url, formData) => {
  const csrfToken = document.querySelector('meta[name="csrf-token"]').content;
  const headers = new Headers({ 'X-CSRF-TOKEN': csrfToken });

  await fetch(url, {
    method: method,
    headers: headers,
    body: formData,
  })
    .then(checkStatus)
    .then(fetchSucceed)
    .catch(fetchFailed);
};

const generateForm = ({ adminOnly, title, subtitle, date, authors, category, tags, outline, content } = values) => {
  const formData = new FormData();

  formData.append('adminOnly', adminOnly);
  formData.append('title', title);
  formData.append('subtitle', subtitle);
  formData.append('date', date);
  formData.append('authors', authors);
  formData.append('category', category);
  formData.append('tags', tags);
  formData.append('outline', outline);
  formData.append('content', content);

  const coverPhoto = document.querySelector('#coverPhoto');
  if (coverPhoto.files.length > 0) {
    formData.append('coverPhoto', coverPhoto.files[0], 'coverPhoto');
  }

  const filesInContent = document.querySelectorAll('input[type="file"]:not(#coverPhoto)');
  for (const f of filesInContent) {
    if (f.files[0] === undefined) {
      continue;
    } else if (f.files[0].size > FILE_MAX_SIZE) {
      document.getElementById('err_msg_content').innerText = `File size of ${f.files[0].name} is too large (max: 8MB per image)!`;
      return;
    } else if (ACCEPT_FILE_TYPE.image.indexOf(f.files[0].type) == -1) {
      document.getElementById('err_msg_content').innerText = `File type of ${f.files[0].name} is not permitted!`;
      return;
    }

    /*
     * Note 1:
     * From "http://127.0.0.1/38c0cbb5a7" to "38c0cbb5a7" (fake image ID provide to user)
     * Frankly speaking, the fake links don't need to have the URL format
     * But images taken from outside world use URL format in markdown syntax, so i decide to use URL format
     *
     * Note 2:
     * f.files[0].name is readonly (can't change a name of a created file), so we add customized name in the 3rd argument
     * This is necessary cause we'll renamed files on server side, and change the URL in the content. So the file name must be the same as fakeID
     * (URL = protocol + domain + fakeID)
     */
    const fakeID = f.nextElementSibling.nextElementSibling.nextElementSibling.innerText.substr(-FILE_ID_LENGTH);
    formData.append('contentImages', f.files[0], fakeID);
  }
  return formData;
};

const submitHandler = async (method, endpoint, button) => {
  button.classList.add('is-loading');
  const values = getInputValue();
  const res = validateInput(values);

  if (!res) {
    button.classList.remove('is-loading');
  } else {
    const formData = generateForm(values);
    if (formData != null) {
      await submitArticle(method, endpoint, formData);
    }
    button.classList.remove('is-loading');
  }
};

const submitPost = () => {
  submitHandler('POST', createArticleEndpoint, submitBtn);
};

const savePost = () => {
  const articleId = new URLSearchParams(window.location.search).get('articleId');
  const baseURL = new URL(window.location.href);
  const url = new URL(updateArticleEndpoint, baseURL);
  url.searchParams.set('articleId', articleId);

  submitHandler('PUT', url, saveBtn);
};

const validateInput = ({ title, subtitle, date, authors, category, tags, content } = values) => {
  let canSubmit = true;

  /*
   * Notice: the TITLE_BYTES_LIMIT, SUBTITLE_BYTES_LIMIT, and TAGS_BYTES_LIMIT does not work for Mandarin
   *         since JS counts the number of words instead of total bytes in the string
   *         Currently this issue is fixed by backend (backend will check the size measured in bytes)
   */

  if (title.length == 0) {
    canSubmit = false;
    document.getElementById('err_msg_title').innerText = errInputMsg.empty;
  } else if (title.length > TITLE_BYTES_LIMIT) {
    canSubmit = false;
    document.getElementById('err_msg_title').innerText = errInputMsg.long;
  } else {
    document.getElementById('err_msg_title').innerText = '';
  }

  if (subtitle.length > SUBTITLE_BYTES_LIMIT) {
    // Subtitle can be empty
    canSubmit = false;
    document.getElementById('err_msg_subtitle').innerText = errInputMsg.long;
  } else {
    document.getElementById('err_msg_subtitle').innerText = '';
  }

  // var currentTime = new Date();
  if (!/^\d\d\d\d-\d\d-\d\d$/.test(date)) {
    canSubmit = false;
    document.getElementById('err_msg_date').innerText = errInputMsg.dateFormat;
  } else if (!/^(19[6-9][0-9]|2\d\d\d)-(3[01]|[12][0-9]|0[1-9])-(0[1-9]|[12][0-9]|3[01])$/.test(date)) {
    // The regex can't detect 02/30, 04/31 .... Nevertheless, Beckend will fix this error
    canSubmit = false;
    document.getElementById('err_msg_date').innerText = errInputMsg.dateIllegal;
  } else if (Number(date.split('-')[0]) < 1960) {
    // Can't detect if the input date is in the current year but in a future month and/or day. Nevertheless, Beckend will fix this error
    canSubmit = false;
    document.getElementById('err_msg_date').innerText = errInputMsg.dateTooOld;
    // } else if (Number(date.split('-')[0]) > currentTime.getFullYear()) {
    //     canSubmit = false;
    //     document.getElementById('err_msg_date').innerText = errInputMsg.dateFuture;
  } else {
    document.getElementById('err_msg_date').innerText = '';
  }

  if (authors.length == 0) {
    canSubmit = false;
    document.getElementById('err_msg_authors').innerText = errInputMsg.empty;
  } else {
    document.getElementById('err_msg_authors').innerText = '';
  }

  if (tags.length > TAGS_LIMIT) {
    canSubmit = false;
    document.getElementById('err_msg_tags').innerText = errInputMsg.tagsTooMany;
  } else {
    document.getElementById('err_msg_tags').innerText = '';
  }
  for (let t of tags) {
    if (t.length > TAGS_BYTES_LIMIT) {
      canSubmit = false;
      document.getElementById('err_msg_tags').innerText = errInputMsg.tagsTooLong;
      break;
    }
  }

  if (content.length == 0) {
    canSubmit = false;
    document.getElementById('err_msg_content').innerText = errInputMsg.empty;
  } else {
    document.getElementById('err_msg_content').innerText = '';
  }

  return canSubmit;
};

const editorHandler = () => {
  // Since the prompt message can't be customized, and the default message is "Changes you made may not be saved"
  // which is quite confusing (cause the changes have been saved). So I remove this feature.
  // window.addEventListener("beforeunload", (e) => {
  //   saveInputToLocalStorage();
  //   e.preventDefault(); // If you prevent default behavior in Mozilla Firefox prompt will always be shown
  //   e.returnValue = ""; // Chrome requires returnValue to be set
  //   return "Are you sure you want to leave? All your changes will be saved."; // Be ignored in modern browsers
  // });

  easyMDE = loadMarkdownEditor();

  submitBtn = document.querySelector('#submitButton');
  if (submitBtn) {
    submitBtn.addEventListener('click', submitPost);
  }

  saveBtn = document.getElementById('saveButton');
  if (saveBtn) {
    saveBtn.addEventListener('click', savePost);
  }

  cancelBtn = document.getElementById('cancelButton'); // Displays in update articles page but not create articles page
  if (cancelBtn) {
    cancelBtn.addEventListener('click', () => {
      const articleId = new URLSearchParams(window.location.search).get('articleId');
      const baseURL = new URL(window.location.href);
      const url = new URL(browseArticleEndpoint, baseURL);
      url.searchParams.set('articleId', articleId);

      window.location.href = url;
    });
  }
};

onDOMContentLoaded = (function () {
  editorHandler();
})();
