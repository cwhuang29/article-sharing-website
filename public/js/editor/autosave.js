const getLocalStorageKey = () => {
  let key = 'article-create';
  if (window.location.pathname.indexOf('create') == -1) {
    key = 'article-update-' + new URLSearchParams(window.location.search).get('articleId');
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
    // If user is updating articles, the INITIAL_INPUT_SIZE may be super large (i.e. the size of original article)
    window.localStorage.setItem(key, JSON.stringify(values));
  }
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

const writeLocalStorageValue = (values) => {
  const { adminOnly, title, subtitle, date, authors, category, tags, outline, content } = JSON.parse(values);

  document.querySelector('#adminOnly').checked = adminOnly;
  document.getElementsByName('title')[0].value = title;
  document.getElementsByName('subtitle')[0].value = subtitle;
  document.getElementsByName('date')[0].value = date;
  document.getElementsByName('outline')[0].value = outline;
  easyMDE.value(content);

  [...document.getElementsByName('category')[0]].forEach((ele, idx) => {
    if (ele.value.toLowerCase() == category) {
      document.getElementsByName('category')[0].selectedIndex = idx;
    }
  });

  [...document.getElementsByName('authors')].forEach((ele, idx) => {
    if (authors.includes(ele.value)) {
      document.getElementsByName('authors')[idx].checked = true;
    }
  });

  var tagsHTMLHead = '<span class="tag is-warning is-medium" name="tags" style="margin-right: 8px; margin-bottom: 5px">';
  var tagsHTMLTail = '<button class="delete is-small"></button></span>';
  var tagsBody = '';
  tags.forEach((ele) => {
    tagsBody += `${tagsHTMLHead}${encodeHTMLEntities(ele)}${tagsHTMLTail}`;
  });
  document.getElementById('tags-list').innerHTML = tagsBody;
};

const setupLocalStorage = () => {
  const key = getLocalStorageKey();
  const values = window.localStorage.getItem(key);
  if (!values) {
    showNoticeMsg('Article will be automatically saved', 'Enjoy your journey : )');
  } else {
    showNoticeMsg('You can now continue editing', 'Enjoy your journey : )');
    writeLocalStorageValue(values);
  }
};

const autosaveHandler = () => {
  setupLocalStorage();

  const iniValues = getInputValue();
  INITIAL_INPUT_SIZE = calculateObjectValueSize(iniValues);

  window.setInterval(() => saveInputToLocalStorage(), 5000);
  document.getElementById('clearAutosaveButton').addEventListener('click', () => clearInputLocalStorage());
  document.getElementById('saveNowButton').addEventListener('click', () => saveInputToLocalStorage());
};

onDOMContentLoaded = (function () {
  autosaveHandler();
})();
