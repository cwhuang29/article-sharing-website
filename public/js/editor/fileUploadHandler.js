const FILE_ID_LENGTH = 10;
const FILES_UPLOAD_LIMIT = 10;
const baseURL = window.location.protocol + '//' + window.location.host + '/';
const fileUploadDefaultMsg = 'No image uploaded';
let FILES_COUNT = 1;

const dec2hex = (dec) => {
  return dec.toString(16).padStart(2, '0');
};

const generateId = (len) => {
  var arr = new Uint8Array((len || 40) / 2);
  window.crypto.getRandomValues(arr);
  return Array.from(arr, dec2hex).join('');
};

const createFileIDField = () => {
  let imgURL = baseURL + generateId(FILE_ID_LENGTH);
  let d = document.createElement('span');
  let classesToAdd = ['file-name', 'fake-id'];
  d.classList.add(...classesToAdd);
  d.style.paddingRight = '15px';
  d.style.cursor = 'default';
  d.style.userSelect = 'all';
  d.style.WebkitTransition.userSelect = 'all'; // Chrome 49+
  d.textContent = imgURL;
  d.addEventListener('click', (e) => e.preventDefault()); // User has to click on this field to copy fake URL
  return d;
};

const createFileUploadTag = () => {
  FILES_COUNT += 1;
  let fileInputTemplate = `<label class='file-label'>
              <input class='file-input' type='file'>
              <span class='file-cta'>
                <span class='file-icon'> 📂 </span>
                <span class='file-label'>Upload images</span>
              </span>
              <span class='file-name'>${fileUploadDefaultMsg}</span>
            </label>`;
  let d = document.createElement('div');
  let classesToAdd = ['file', 'has-name', 'is-warning', 'is-small', 'pb-1'];
  d.classList.add(...classesToAdd);
  d.innerHTML = fileInputTemplate;
  return d;
};

const removeFileUploadTag = (target) => {
  target.remove();
  FILES_COUNT -= 1;
  if (FILES_COUNT == 0 || FILES_COUNT == FILES_UPLOAD_LIMIT - 1) {
    filesInContent.appendChild(createFileUploadTag());
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
    if (val == fileUploadDefaultMsg && FILES_COUNT < FILES_UPLOAD_LIMIT) {
      filesInContent.appendChild(createFileUploadTag());
    }
  }
};

const fileUploadHandler = () => {
  const coverPhoto = document.querySelector('#coverPhoto');
  coverPhoto.onchange = () => {
    if (coverPhoto.files.length > 0) {
      const fileName = document.querySelector('#coverPhotoName');
      fileName.textContent = coverPhoto.files[0].name;
    }
  };

  filesInContent = document.getElementById('filesGroupInContent');

  filesInContent.addEventListener('change', (e) => {
    if (e.target.tagName.toLowerCase() == 'input') {
      displayFileNameAndCancelBtn(e.target);
    }
  });

  filesInContent.addEventListener('click', (e) => {
    if (e.target.tagName.toLowerCase() == 'button') {
      removeFileUploadTag(e.target.parentNode.parentNode.parentNode);
    }
  });
};

onDOMContentLoaded = (function () {
  fileUploadHandler();
})();
